// Package main is the entry point for the CarZone application.
// It sets up the HTTP server, initializes dependencies, and configures routing.
// This file follows the dependency injection pattern to wire up all components.
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	// Database connection management
	"github.com/PrateekKumar15/CarZone/driver"

	// HTTP handlers for car and engine endpoints
	carHandler "github.com/PrateekKumar15/CarZone/handler/car"
	engineHandler "github.com/PrateekKumar15/CarZone/handler/engine"

	// Business logic services
	carService "github.com/PrateekKumar15/CarZone/service/car"
	engineService "github.com/PrateekKumar15/CarZone/service/engine"

	// Data access layer stores
	carStore "github.com/PrateekKumar15/CarZone/store/car"
	engineStore "github.com/PrateekKumar15/CarZone/store/engine"

	// Third-party dependencies
	"github.com/gorilla/mux"   // HTTP router and URL matcher
	"github.com/joho/godotenv" // Environment variable loader
)

// main is the entry point of the CarZone application.
// It initializes all dependencies, sets up the HTTP server, and starts the application.
// The function follows these main steps:
// 1. Load environment variables from .env file
// 2. Initialize database connection and ensure proper cleanup
// 3. Set up dependency injection chain: stores -> services -> handlers
// 4. Configure HTTP routes using Gorilla Mux router
// 5. Start the HTTP server on the specified port
func main() {
	// Step 1: Load environment variables from .env file
	// This allows configuration without hardcoding values
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Step 2: Initialize database connection
	// The driver package handles PostgreSQL connection setup
	driver.InitDB()
	// Ensure database connection is properly closed when application exits
	// This is critical for preventing connection leaks
	defer driver.CloseDB()

	// Get the database connection instance
	db := driver.GetDB()
	if db == nil {
		log.Fatal("Database connection is nil - cannot proceed")
	}

	// Step 3: Set up dependency injection chain following clean architecture
	// Data Access Layer (Stores) - Handle database operations
	carStore := carStore.New(db)
	engineStore := engineStore.New(db)

	// Business Logic Layer (Services) - Handle domain logic and validation
	carService := carService.NewCarService(carStore)
	engineService := engineService.NewEngineService(engineStore)

	// Presentation Layer (Handlers) - Handle HTTP requests/responses
	carHandler := carHandler.NewCarHandler(carService)
	engineHandler := engineHandler.NewEngineHandler(engineService)

	// Step 4: Configure HTTP routing using Gorilla Mux
	router := mux.NewRouter()

	// Execute schema file to set up database structure
	// This is typically done once during application startup
	// It ensures the database is ready for operations
	// This function executes the SQL commands in the schema file
	executeSchemaFile := func(db *sql.DB, schemaFile string) error {
		schema, err := os.ReadFile(schemaFile)
		if err != nil {
			fmt.Printf("Error reading schema file %s: %v\n", schemaFile, err)
			return err
		}
		// Execute the schema SQL commands
		_, err = db.Exec(string(schema))
		if err != nil {
			fmt.Printf("Error executing schema file %s: %v\n", schemaFile, err)
			return err
		}
		return nil
	}


	schemaFile :="store/schema.sql"
	if err := executeSchemaFile(db, schemaFile); err != nil {
		log.Fatalf("Failed to execute schema file %s: %v", schemaFile, err)
	}


	// Car-related endpoints
	// GET /cars/{id} - Retrieve a specific car by its UUID
	router.HandleFunc("/cars/{id}", carHandler.GetCarByID).Methods("GET")

	// GET /cars?brand={brand}&engine={true/false} - Retrieve cars by brand with optional engine details
	router.HandleFunc("/cars", carHandler.GetCarByBrand).Methods("GET")

	// POST /cars - Create a new car record
	router.HandleFunc("/cars", carHandler.CreateCar).Methods("POST")

	// PUT /cars/{id} - Update an existing car by its UUID
	router.HandleFunc("/cars/{id}", carHandler.UpdateCar).Methods("PUT")

	// DELETE /cars/{id} - Delete a car by its UUID
	router.HandleFunc("/cars/{id}", carHandler.DeleteCar).Methods("DELETE")

	// Engine-related endpoints
	// GET /engines/{id} - Retrieve a specific engine by its UUID
	router.HandleFunc("/engines/{id}", engineHandler.GetEngineByID).Methods("GET")

	// GET /engines?brand={brand} - Retrieve engines by brand
	router.HandleFunc("/engines", engineHandler.GetEngineByBrand).Methods("GET")

	// POST /engines - Create a new engine record
	router.HandleFunc("/engines", engineHandler.CreateEngine).Methods("POST")

	// PUT /engines/{id} - Update an existing engine by its UUID
	router.HandleFunc("/engines/{id}", engineHandler.UpdateEngine).Methods("PUT")

	// DELETE /engines/{id} - Delete an engine by its UUID
	router.HandleFunc("/engines/{id}", engineHandler.DeleteEngine).Methods("DELETE")

	// Step 5: Start the HTTP server
	// Get port from environment variables with fallback to default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not set in environment variables
	}

	// Log server startup information
	log.Printf("Starting CarZone server on port %s", port)
	log.Println("Available endpoints:")
	log.Println("Cars: GET,POST /cars | GET,PUT,DELETE /cars/{id}")
	log.Println("Engines: GET,POST /engines | GET,PUT,DELETE /engines/{id}")

	// Start the HTTP server - this blocks until server shuts down
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
