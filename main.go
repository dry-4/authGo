package main

import (
	"fmt"
	"log"
	"net/http"

	configs "hells/config"
	"hells/routes"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database connection
	db, err := configs.InitDatabase()

	if err != nil {
		fmt.Errorf("Database initialiaztion failed: %v", err)
		log.Fatalf("Database initialiaztion failed: %v", err)
	}
	fmt.Println(db)

	// Create router
	router := mux.NewRouter()

	// Setup routes
	routes.SetupRoutes(router)

	// Start server
	port := ":8080"
	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
