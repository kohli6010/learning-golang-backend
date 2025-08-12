package handler

import (
	"encoding/json"
	"golang-crud-ddd/apprequest"
	"golang-crud-ddd/helper"
	"golang-crud-ddd/middleware"
	"golang-crud-ddd/service"
	"log"
	"net/http"
	"strings"
	"time"
)

// CreateEndpointForGetUserByID ...
func CreateEndpointForGetUserByID(svc service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Logic to handle getting user by ID
		userID, ok := middleware.GetUserIDFromContext(r.Context())
		if !ok {
			http.Error(w, "Request Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := svc.GetUserByID(userID)
		if err != nil {
			http.Error(w, "Failed to get user", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

// CreateEndpointForCreateUser ...
func CreateEndpointForCreateUser(svc service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Logic to handle user creation

		var userRequest apprequest.UserRequest
		// Decode the request body into userRequest
		err := json.NewDecoder(r.Body).Decode(&userRequest)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Call the service to create the user
		resp, err := svc.CreateUser(&userRequest)
		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}
}

// CreateEndpointForAuthenticateUser ...
func CreateEndpointForAuthenticateUser(svc service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Logic to handle user authentication

		var authRequest apprequest.AuthenticateUserRequest
		err := json.NewDecoder(r.Body).Decode(&authRequest)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		response, err := svc.AuthenticateUser(&authRequest)
		if err != nil {
			http.Error(w, "Authentication failed", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// CreateEndpointForChangePassword ...
func CreateEndpointForChangePassword(svc service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Logic to handle changing user password
		_, ok := middleware.GetUserIDFromContext(r.Context())
		if !ok {
			http.Error(w, "Request Unauthorized", http.StatusUnauthorized)
			return
		}

		var changePasswordRequest apprequest.ChangePasswordRequest
		err := json.NewDecoder(r.Body).Decode(&changePasswordRequest)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err = svc.ChangePassword(&changePasswordRequest)
		if err != nil {
			http.Error(w, "Failed to change password", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Password changed successfully"))
	}
}

// CreateEndpointForRefreshTokenss ...
func CreateEndpointForRefreshTokenss(svc service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Logic to handle refreshing tokens
		var refreshRequest apprequest.RefreshTokensRequest
		err := json.NewDecoder(r.Body).Decode(&refreshRequest)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		response, err := svc.RefreshTokenss(refreshRequest.OldRefreshTokens)
		if err != nil {
			http.Error(w, "Failed to refresh tokens", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// CreateEndpointForLogout ...
func CreateEndpointForLogout(svc service.UserService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()

        // 1. Get refresh token from request body (JSON)
        var req struct {
            RefreshToken string `json:"refresh_token"`
        }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RefreshToken == "" {
            http.Error(w, "Refresh token is required", http.StatusBadRequest)
            return
        }

        // 2. Extract access token from Authorization header
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Access token is required", http.StatusBadRequest)
            return
        }

        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
            http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
            return
        }
        accessToken := parts[1]

        // 3. Verify and decode access token (to get jti & expiry)
        claims, err := helper.VerifyAndDecodeToken(accessToken)
        if err != nil {
            http.Error(w, "Invalid or expired access token", http.StatusUnauthorized)
            return
        }

        jti, ok := claims["jti"].(string)
        if !ok || jti == "" {
            http.Error(w, "No token ID (jti) present", http.StatusUnauthorized)
            return
        }

        expFloat, ok := claims["exp"].(float64)
        if !ok {
            http.Error(w, "Invalid token expiry", http.StatusUnauthorized)
            return
        }
        expTime := time.Unix(int64(expFloat), 0)
        ttl := time.Until(expTime)
        if ttl <= 0 {
            ttl = time.Second // ensure positive TTL to store temporarily
        }

        // 4. Blacklist the current access token's JTI in Redis
        if err := helper.BlacklistToken(ctx, jti, ttl); err != nil {
            log.Printf("Failed to blacklist token in Redis: %v", err)
        }

        // 5. Revoke refresh token in DB
        if err := svc.Logout(req.RefreshToken); err != nil {
            http.Error(w, "Failed to logout", http.StatusInternalServerError)
            return
        }

        // 6. Send success response
        w.Header().Set("Content-Type", "application/json")
        _ = json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
    }
}

// CreateEndpointForUpdateUserByID ...
func CreateEndpointForUpdateUserByID(svc service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Logic to handle updating user by ID
		userID, ok := middleware.GetUserIDFromContext(r.Context())
		if !ok {
			http.Error(w, "Request Unauthorized", http.StatusUnauthorized)
			return
		}

		var updateRequest apprequest.UpsertUserPersonalDetailRequest
		err := json.NewDecoder(r.Body).Decode(&updateRequest)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		response, err := svc.UpdateUserByID(userID, &updateRequest)
		if err != nil {
			http.Error(w, "Failed to update user", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

