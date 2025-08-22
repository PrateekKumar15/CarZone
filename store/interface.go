// Package store defines the data access layer interfaces for the CarZone application.
// These interfaces establish contracts for data persistence operations and follow
// the Repository pattern to abstract database operations from business logic.
// All store implementations must adhere to these interfaces to ensure consistency.
package store

import (
	"context"

	"github.com/PrateekKumar15/CarZone/models"
)

// CarStoreInterface defines the contract for car data access operations.
// This interface abstracts all database operations related to car entities,
// following the Repository pattern to decouple business logic from data persistence.
// All methods accept a context for request scoping, cancellation, and timeout handling.
type CarStoreInterface interface {
	// GetCarByID retrieves a single car record by its unique identifier.
	// Parameters:
	//   - ctx: Request context for cancellation and timeout
	//   - id: Unique identifier of the car (UUID string format)
	// Returns:
	//   - models.Car: The car record if found
	//   - error: Error if car not found or database operation fails
	GetCarByID(ctx context.Context, id string) (models.Car, error)

	// GetCarByBrand retrieves multiple car records filtered by brand name.
	// The isEngine parameter controls whether engine information is included.
	// Parameters:
	//   - ctx: Request context for cancellation and timeout
	//   - brand: Brand name to filter by (e.g., "Toyota", "BMW")
	//   - isEngine: Boolean flag to include/exclude engine details
	// Returns:
	//   - []models.Car: Slice of car records matching the brand
	//   - error: Error if database operation fails
	GetCarByBrand(ctx context.Context, brand string, isEngine bool) ([]models.Car, error)

	// CreateCar inserts a new car record into the database.
	// The method generates a new UUID for the car and handles all creation logic.
	// Parameters:
	//   - ctx: Request context for transaction management
	//   - carReq: Car data to be inserted (without ID, timestamps)
	// Returns:
	//   - models.Car: The created car record with generated ID and timestamps
	//   - error: Error if creation fails or validation errors occur
	CreateCar(ctx context.Context, carReq models.CarRequest) (models.Car, error)

	// UpdateCar modifies an existing car record with new data.
	// Only the fields provided in carReq will be updated.
	// Parameters:
	//   - ctx: Request context for transaction management
	//   - id: Unique identifier of the car to update
	//   - carReq: New car data to replace existing values
	// Returns:
	//   - models.Car: The updated car record
	//   - error: Error if car not found or update operation fails
	UpdateCar(ctx context.Context, id string, carReq models.CarRequest) (models.Car, error)

	// DeleteCar removes a car record from the database.
	// This operation is typically irreversible and should be used with caution.
	// Parameters:
	//   - ctx: Request context for transaction management
	//   - id: Unique identifier of the car to delete
	// Returns:
	//   - models.Car: The deleted car record (for logging/audit purposes)
	//   - error: Error if car not found or deletion fails
	DeleteCar(ctx context.Context, id string) (models.Car, error)
}

// EngineStoreInterface defines the contract for engine data access operations.
// This interface abstracts all database operations related to engine entities,
// providing a consistent API for engine data management across the application.
// All methods use context for proper request handling and cancellation support.
type EngineStoreInterface interface {
	// GetEngineByID retrieves a single engine record by its unique identifier.
	// Parameters:
	//   - ctx: Request context for cancellation and timeout
	//   - id: Unique identifier of the engine (UUID string format)
	// Returns:
	//   - models.Engine: The engine record if found
	//   - error: Error if engine not found or database operation fails
	GetEngineByID(ctx context.Context, id string) (models.Engine, error)

	// CreateEngine inserts a new engine record into the database.
	// The method generates a new UUID for the engine and validates all input data.
	// Parameters:
	//   - ctx: Request context for transaction management
	//   - engineReq: Engine specification data to be inserted
	// Returns:
	//   - models.Engine: The created engine record with generated ID
	//   - error: Error if creation fails or validation errors occur
	CreateEngine(ctx context.Context, engineReq models.EngineRequest) (models.Engine, error)

	// UpdateEngine modifies an existing engine record with new specifications.
	// All engine parameters can be updated while maintaining referential integrity.
	// Parameters:
	//   - ctx: Request context for transaction management
	//   - id: Unique identifier of the engine to update
	//   - engineReq: New engine specifications to replace existing values
	// Returns:
	//   - models.Engine: The updated engine record
	//   - error: Error if engine not found or update operation fails
	UpdateEngine(ctx context.Context, id string, engineReq models.EngineRequest) (models.Engine, error)

	// DeleteEngine removes an engine record from the database.
	// This operation may cascade to related car records depending on foreign key constraints.
	// Parameters:
	//   - ctx: Request context for transaction management
	//   - id: Unique identifier of the engine to delete
	// Returns:
	//   - models.Engine: The deleted engine record (for logging/audit purposes)
	//   - error: Error if engine not found or deletion fails due to constraints
	DeleteEngine(ctx context.Context, id string) (models.Engine, error)

	// GetEngineByBrand retrieves multiple engine records filtered by brand.
	// This method allows querying engines based on the brand of cars they are associated with.
	// Parameters:
	//   - ctx: Request context for cancellation and timeout
	//   - brand: Brand name to filter engines by
	// Returns:
	//   - []models.Engine: Slice of engine records associated with the specified brand
	//   - error: Error if database operation fails
	GetEngineByBrand(ctx context.Context, brand string) ([]models.Engine, error)
}
