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

	// Routes layer
	"github.com/PrateekKumar15/CarZone/routes"

	// HTTP handlers for car endpoints
	carHandler "github.com/PrateekKumar15/CarZone/handler/car"

	// HTTP handlers for booking endpoints
	bookingHandler "github.com/PrateekKumar15/CarZone/handler/booking"

	// Business logic services
	carService "github.com/PrateekKumar15/CarZone/service/car"

	// Business logic services for booking
	bookingService "github.com/PrateekKumar15/CarZone/service/booking"

	// Data access layer stores
	carStore "github.com/PrateekKumar15/CarZone/store/car"

	// Data access layer for booking
	bookingStore "github.com/PrateekKumar15/CarZone/store/booking"

	// Data access layer for payment
	paymentStore "github.com/PrateekKumar15/CarZone/store/payment"

	// Third-party dependencies
	authHandler "github.com/PrateekKumar15/CarZone/handler/auth"
	authService "github.com/PrateekKumar15/CarZone/service/auth"
	userStore "github.com/PrateekKumar15/CarZone/store/user"

	// Payment components
	paymentHandler "github.com/PrateekKumar15/CarZone/handler/payment"
	paymentService "github.com/PrateekKumar15/CarZone/service/payment"
	"github.com/joho/godotenv" // Environment variable loader
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
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

	bookingStore := bookingStore.New(db)

	userStore := userStore.New(db)

	paymentStore := paymentStore.New(db)

	// Business Logic Layer (Services) - Handle domain logic and validation
	carService := carService.NewCarService(carStore)
	bookingService := bookingService.NewBookingService(bookingStore, carStore)
	authService := authService.NewAuthService(userStore)
	paymentService := paymentService.NewPaymentService(paymentStore, bookingStore)

	// Presentation Layer (Handlers) - Handle HTTP requests/responses
	carHandler := carHandler.NewCarHandler(carService)
	bookingHandler := bookingHandler.NewBookingHandler(bookingService)
	authHandler := authHandler.NewAuthHandler(authService)
	paymentHandler := paymentHandler.NewPaymentHandler(paymentService)

	// Step 4: Initialize routes using the routes layer
	// Create router with all handler dependencies injected
	routeManager := routes.NewRouter(authHandler, carHandler, bookingHandler, paymentHandler)
	router := routeManager.SetupRoutes()

	// Execute schema file to set up database structure
	// This is typically done once during application startup
	// It ensures the database is ready for operations
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

	// Step 5: Start the HTTP server
	// Get port from environment variables with fallback to default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not set in environment variables
	}

	// Log server startup information with organized route categories
	log.Printf("Starting CarZone server on port %s", port)
	log.Println("🚀 CarZone API Server Started Successfully!")
	log.Println("")
	log.Println("📋 Available API Routes:")
	log.Println("  🔐 Authentication (Public):")
	log.Println("    POST /auth/register  - Register new user account")
	log.Println("    POST /auth/login     - User authentication")
	log.Println("    GET  /auth/logout    - User logout")
	log.Println("")
	log.Println("  🚗 Car Management (Protected):")
	log.Println("    GET    /cars           - Get all cars")
	log.Println("    GET    /cars/{id}      - Get car by ID")
	log.Println("    GET    /cars/brand     - Get cars by brand")
	log.Println("    POST   /cars           - Create new car")
	log.Println("    PUT    /cars/{id}      - Update car")
	log.Println("    DELETE /cars/{id}      - Delete car")
	log.Println("")
	log.Println("  📅 Booking Management (Protected):")
	log.Println("    GET    /bookings                    - Get all bookings")
	log.Println("    GET    /bookings/{id}               - Get booking by ID")
	log.Println("    POST   /bookings                    - Create new booking")
	log.Println("    DELETE /bookings/{id}               - Delete booking")
	log.Println("    PUT    /bookings/{id}/status        - Update booking status")
	log.Println("    GET    /bookings/customer/{id}      - Get bookings by customer")
	log.Println("    GET    /bookings/car/{id}           - Get bookings by car")
	log.Println("    GET    /bookings/owner/{id}         - Get bookings by owner")
	log.Println("")
	log.Println("  💳 Payment Management (Protected):")
	log.Println("    POST   /payments                     - Create payment and Razorpay order")
	log.Println("    POST   /payments/verify              - Verify payment signature")
	log.Println("    GET    /payments/{id}                - Get payment by ID")
	log.Println("    GET    /payments/booking/{booking_id} - Get payment by booking ID")
	log.Println("    GET    /payments/user/{user_id}      - Get payments by user ID")
	log.Println("    POST   /payments/{payment_id}/refund - Process payment refund")
	log.Println("    GET    /payments                     - Get all payments")
	log.Println("")
	log.Println("  📊 Monitoring:")
	log.Println("    GET /metrics - Prometheus metrics")
	log.Println("")
	log.Println("✨ Routes are organized using the new routes layer for better maintainability!")

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
