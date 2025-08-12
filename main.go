package main

import (
	"log"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables from env.yml
	err := godotenv.Load("env.yml")
	if err != nil {
		log.Fatalf("Error loading env.yml file: %v", err)
	}
	// register db
	registerDB()
	
	// Setup routes
	SetupRoutes()
}