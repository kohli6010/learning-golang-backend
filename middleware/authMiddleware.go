package middleware

import (
	"context"
	"golang-crud-ddd/helper"
	"net/http"
	"strings"
)

type ctxKey string

const userIDKey ctxKey = "userID"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenStr := parts[1]

		claims, err := helper.VerifyAndDecodeToken(tokenStr)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Claims["id"] will be float64
		idFloat, ok := claims["id"].(float64)
		if !ok {
			http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
			return
		}

		jti, ok := claims["jti"].(string)
		if !ok || jti == "" {
			http.Error(w, "Invalid token identifier", http.StatusUnauthorized)
			return
		}

		ctx := context.Background()
		isTokenBlacklisted, err := helper.IsTokenBlacklisted(ctx, jti)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if isTokenBlacklisted {
			http.Error(w, "Token is blacklisted", http.StatusUnauthorized)
			return
		}

		// Store userID in context
		ctx = context.WithValue(r.Context(), userIDKey, int(idFloat))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserIDFromContext returns the user ID stored by the middleware
func GetUserIDFromContext(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(userIDKey).(int)
	return id, ok
}
