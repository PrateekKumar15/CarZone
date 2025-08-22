// Package driver provides database connection management for the CarZone application.
// It handles PostgreSQL database connections using the lib/pq driver and manages
// connection lifecycle including initialization, retrieval, and cleanup.
package driver

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	// PostgreSQL driver - imported for side effects (driver registration)
	_ "github.com/lib/pq"
)

// db is a package-level variable that holds the database connection pool.
// Using a singleton pattern ensures all parts of the application share the same connection pool.
var db *sql.DB

// InitDB initializes the PostgreSQL database connection pool.
// It reads database configuration from environment variables and establishes
// a connection with proper error handling and connection validation.
// This function should be called once during application startup.
func InitDB() {
	// Build connection string from environment variables
	// Format: "host=localhost port=5432 user=username password=password dbname=database sslmode=disable"
	host := os.Getenv("DB_HOST")
	portStr := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	// Set default values for missing environment variables
	if host == "" {
		host = "localhost"
	}
	if portStr == "" {
		portStr = "5432"
	}
	if sslmode == "" {
		sslmode = "disable"
	}

	// Convert port string to integer for validation
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid DB_PORT value: %v", err)
	}

	// Validate required environment variables
	if user == "" {
		log.Fatal("DB_USER environment variable is required")
	}
	if password == "" {
		log.Fatal("DB_PASSWORD environment variable is required")
	}
	if dbname == "" {
		log.Fatal("DB_NAME environment variable is required")
	}

	// Construct PostgreSQL connection string
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	log.Printf("Connecting to database at %s:%d/%s...", host, port, dbname)

	// Add a small delay for containerized environments where database might be starting
	log.Println("Waiting for database to be ready...")
	time.Sleep(5 * time.Second)

	// Open database connection pool
	// sql.Open() doesn't actually connect - it just prepares the database connection pool
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}

	// Configure connection pool settings for optimal performance
	db.SetMaxOpenConns(25)                 // Maximum number of open connections
	db.SetMaxIdleConns(10)                 // Maximum number of idle connections
	db.SetConnMaxLifetime(5 * time.Minute) // Maximum lifetime of a connection

	// Test the database connection by pinging it
	// This actually establishes a connection to verify everything is working
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Successfully connected to PostgreSQL database")
	log.Printf("Connection pool configured with MaxOpen=%d, MaxIdle=%d",
		25, 10)
}

// GetDB returns the singleton database connection pool instance.
// This function provides access to the database connection throughout the application.
// It returns the same *sql.DB instance that was initialized by InitDB().
//
// Returns:
//   - *sql.DB: The database connection pool instance, or nil if not initialized
func GetDB() *sql.DB {
	if db == nil {
		log.Println("Warning: Database connection is nil. Did you call InitDB()?")
		return nil // Return nil if database connection is not initialized
	}
	return db
}

// CloseDB gracefully closes the database connection pool.
// This function should be called during application shutdown to ensure
// all database connections are properly closed and resources are freed.
// It's typically called using defer in the main function.
func CloseDB() {
	if db == nil {
		log.Println("Database connection is already nil, nothing to close")
		return
	}

	// Close the database connection pool
	if err := db.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
	} else {
		log.Println("Database connection closed successfully")
	}

	// Set the package variable to nil to prevent further use
	db = nil
}
