// Package main is the entry point for the CarZone application.
// It sets up the HTTP server, initializes dependencies, and configures routing.
// This file follows the dependency injection pattern to wire up all components.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
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
	loginHandler "github.com/PrateekKumar15/CarZone/handler/login"
	"github.com/PrateekKumar15/CarZone/middleware"
	"github.com/gorilla/mux"   // HTTP router and URL matcher
	"github.com/joho/godotenv" // Environment variable loader
	// "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.4.0"
	// Prometheus metrics
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	traceProvider, err := startTracing()
	if err != nil {
		log.Fatalf("Failed to start tracing: %v", err)
	}
	defer func() {
		if err := traceProvider.Shutdown(context.Background()); err != nil {
			log.Fatalf("Failed to shutdown tracer provider: %v", err)
		}
	}()
	// Set global tracer provider
	// This enables tracing throughout the application
	otel.SetTracerProvider(traceProvider)

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

	schemaFile := "store/schema.sql"
	if err := executeSchemaFile(db, schemaFile); err != nil {
		log.Fatalf("Failed to execute schema file %s: %v", schemaFile, err)
	}
	router.Use((otelmux.Middleware("CarZone")))
	// Authentication endpoint
	router.HandleFunc("/login", loginHandler.LoginHandler).Methods("POST")
	// Prometheus metrics endpoint
	router.Handle("/metrics", promhttp.Handler())
	// Public endpoints (no auth required)
	// You can add more public endpoints here if needed

	// Protected endpoints (require authentication)
	protected := router.PathPrefix("/").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.Use(middleware.MetricMiddleware)

	// Car-related endpoints
	// GET /cars/all - Retrieve all cars (must come before /cars/{id})
	protected.HandleFunc("/cars", carHandler.GetAllCars).Methods("GET")

	// GET /cars/{id} - Retrieve a specific car by its UUID
	protected.HandleFunc("/cars/{id}", carHandler.GetCarByID).Methods("GET")

	// GET /cars/brand?brand={brand}&engine={true/false} - Retrieve cars by brand with optional engine details
	protected.HandleFunc("/carbybrand", carHandler.GetCarByBrand).Methods("GET")

	// POST /cars - Create a new car record
	protected.HandleFunc("/cars", carHandler.CreateCar).Methods("POST")

	// PUT /cars/{id} - Update an existing car by its UUID
	protected.HandleFunc("/cars/{id}", carHandler.UpdateCar).Methods("PUT")

	// DELETE /cars/{id} - Delete a car by its UUID
	protected.HandleFunc("/cars/{id}", carHandler.DeleteCar).Methods("DELETE")

	// Engine-related endpoints
	// GET /engines/{id} - Retrieve a specific engine by its UUID
	protected.HandleFunc("/engines/{id}", engineHandler.GetEngineByID).Methods("GET")


	// POST /engines - Create a new engine record
	protected.HandleFunc("/engines", engineHandler.CreateEngine).Methods("POST")

	// PUT /engines/{id} - Update an existing engine by its UUID
	protected.HandleFunc("/engines/{id}", engineHandler.UpdateEngine).Methods("PUT")

	// DELETE /engines/{id} - Delete an engine by its UUID
	protected.HandleFunc("/engines/{id}", engineHandler.DeleteEngine).Methods("DELETE")

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

func startTracing() (*trace.TracerProvider, error) {
	header := map[string]string{
		"Content-Type": "application/json",
	}
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint("jaeger:4318"),
			otlptracehttp.WithHeaders(header),
			otlptracehttp.WithInsecure(),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", err)
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter,
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
			trace.WithBatchTimeout(trace.DefaultScheduleDelay*time.Millisecond),
		),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("CarZone"),
		)),
	)

	return traceProvider, nil
}
