package main

import (
	"log"
	"net/http"
	"smart-scheduler/config"
	"smart-scheduler/db"
	"smart-scheduler/repository"
	"smart-scheduler/routes"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db.InitDB(cfg.DatabaseURL, cfg.DBName)

	// Populate dummy data for testing
	if err := repository.CreateDummyData(db.DB); err != nil {
		log.Printf("Warning: Failed to create dummy data: %v", err)
	} else {
		log.Println("Dummy data created successfully")
	}

	// Setup routes
	router := routes.SetupRoutes()

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
