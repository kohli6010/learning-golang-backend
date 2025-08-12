package helper

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"sync"

	"github.com/go-redis/redis/v8"
)

// CheckPasswordHash ...
func CheckPasswordHash(password string, encryptedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(password))
	return err == nil
}

// GenerateToken ...
func GenerateToken(id int, email string) (string, string, error) {

	expiresAt := time.Now().Add(time.Hour * 24)
	claims := jwt.MapClaims{
		"id":    id,
		"email": email,
		"exp":   expiresAt.Unix(),
		"iat":   time.Now().Unix(),
		"jti":   uuid.New().String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("your_secret_key"))
	if err != nil {
		return "", "", err
	}

	return tokenString, expiresAt.Format(time.RFC3339), nil

}

// VerifyAndDecodeToken ...
func VerifyAndDecodeToken(tokenStr string) (jwt.MapClaims, error) {
	// Parse and verify token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	// Validate token and extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	fmt.Println("claims:", token.Claims)
	return nil, fmt.Errorf("invalid token")
}

type ctxKey string

const tokenKey ctxKey = "authToken"

// StoreTokenInContext stores token in context
func StoreTokenInContext(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tokenKey, token)
}

// GetTokenFromContext retrieves the JWT token from context
func GetTokenFromContext(c context.Context) (string, error) {
	token, ok := c.Value(tokenKey).(string)
	if !ok || token == "" {
		return "", errors.New("no token found in context")
	}
	return token, nil
}

// ExtractTokenFromHeader takes 'Authorization: Bearer <token>'
func ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header empty")
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("authorization header format must be Bearer {token}")
	}
	return parts[1], nil
}

// GenerateTokens ...
func GenerateTokens(id int, email string) (accessToken, RefreshTokens string, accessExpiry, refreshExpiry time.Time, err error) {
	accessExpiry = time.Now().Add(15 * time.Minute)
	refreshExpiry = time.Now().Add(7 * 24 * time.Hour)

	claims := jwt.MapClaims{
		"id":    id,
		"email": email,
		"exp":   accessExpiry.Unix(),
		"iat":   time.Now().Unix(),
		"jti":   uuid.New().String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	accessToken, err = token.SignedString([]byte(secret))
	if err != nil {
		return
	}

	refreshBytes := make([]byte, 32)
	if _, err = rand.Read(refreshBytes); err != nil {
		return
	}
	RefreshTokens = base64.URLEncoding.EncodeToString(refreshBytes)
	return
}

var (
	redisClient *redis.Client
	once sync.Once
)

func GetRedisClient() *redis.Client {
	once.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: "",
			DB:       0,
		})
	})

	return redisClient
}

// BlacklistToken ...
func BlacklistToken(ctx context.Context, jti string, expiration time.Duration) error {
    client := GetRedisClient()
    key := "blacklist:" + jti
    return client.Set(ctx, key, "revoked", expiration).Err()
}

// IsTokenBlacklisted ...
func IsTokenBlacklisted(ctx context.Context, jti string) (bool, error) {
    client := GetRedisClient()
    key := "blacklist:" + jti
    exists, err := client.Exists(ctx, key).Result()
    if err != nil {
        return false, err
    }
    return exists == 1, nil
}

