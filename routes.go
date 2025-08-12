package main

import (
	"golang-crud-ddd/handler"
	"golang-crud-ddd/middleware"
	"golang-crud-ddd/service"
	"log"
	"net/http"
	"os"

	"golang-crud-ddd/repository"

	"github.com/gorilla/mux"
)

// SetupRoutes ...
func SetupRoutes() {
	router := mux.NewRouter()
	// initialize the user repository and service
	userRepo := repository.NewUserRepo()
	RefreshTokensRepo := repository.NewRefreshTokensRepo()
	userService := service.NewUserService(userRepo, RefreshTokensRepo)

	// Define your routes here
	// Example route for getting user by ID
	router.Handle("/api/v1/users", middleware.AuthMiddleware(handler.CreateEndpointForGetUserByID(userService))).Methods("GET")
	// Example route for user creation
	router.HandleFunc("/api/v1/users/create", handler.CreateEndpointForCreateUser(userService)).Methods("POST")
	// Example route for user authentication
	router.HandleFunc("/api/v1/users/authenticate", handler.CreateEndpointForAuthenticateUser(userService)).Methods("POST")
	// Example route for changing password
	router.Handle("/api/v1/users/change-password", middleware.AuthMiddleware(handler.CreateEndpointForChangePassword(userService))).Methods("POST")
	// Example route for refreshing tokens
	router.Handle("/api/v1/users/refresh-tokens", middleware.AuthMiddleware(handler.CreateEndpointForRefreshTokenss(userService))).Methods("POST")
	// Example route for user logout
	router.Handle("/api/v1/users/logout", middleware.AuthMiddleware(handler.CreateEndpointForLogout(userService))).Methods("POST")
	// Example route for updating user by ID
	router.Handle("/api/v1/users", middleware.AuthMiddleware(handler.CreateEndpointForUpdateUserByID(userService))).Methods("PUT")

	// Start the HTTP server
	serverAddr := os.Getenv("SERVER_ADDR")
	log.Printf("Starting server on %s", serverAddr)
	if err := http.ListenAndServe(serverAddr, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	log.Println("Server stopped")
	// Note: In a real application, you might want to handle graceful shutdowns and other
	// advanced features like logging, middleware, etc.
	log.Println("Routes have been set up successfully")
}
