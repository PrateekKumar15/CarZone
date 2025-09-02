// Package service defines the business logic layer interfaces for the CarZone application.
// These interfaces establish contracts for business operations and coordinate between
// the presentation layer (handlers) and data access layer (stores).
// They encapsulate domain logic, validation, and business rules.
package service

import (
	"context"

	"github.com/PrateekKumar15/CarZone/models"
)

// CarServiceInterface defines the contract for car business logic operations.
// This interface abstracts all business operations related to car entities,
// including validation, business rule enforcement, and coordination with the data layer.
// All methods return pointers to models for consistency and memory efficiency.
type CarServiceInterface interface {
	// GetCarByID retrieves a car by its unique identifier with business logic applied.
	// This method may include additional validation, logging, and business rule checks.
	// Parameters:
	//   - ctx: Request context for cancellation, timeout, and request scoping
	//   - id: Unique identifier of the car (UUID string format)
	// Returns:
	//   - *models.Car: Pointer to the car record if found, nil if not found
	//   - error: Business logic error or underlying data access error
	GetCarByID(ctx context.Context, id string) (*models.Car, error)

	// GetCarByBrand retrieves multiple cars filtered by brand with optional engine details.
	// Applies business rules for data filtering and presentation logic.
	// Parameters:
	//   - ctx: Request context for cancellation and timeout
	//   - brand: Brand name to filter by (case-sensitive)
	//   - isEngine: Whether to include detailed engine information in results
	// Returns:
	//   - *[]models.Car: Pointer to slice of car records matching the criteria
	//   - error: Business logic error or data access error
	GetCarByBrand(ctx context.Context, brand string, isEngine bool) (*[]models.Car, error)

	// CreateCar creates a new car record with full business validation.
	// Validates input data, enforces business rules, and coordinates with data persistence.
	// Parameters:
	//   - ctx: Request context for transaction management
	//   - carReq: Car creation request with all required fields
	// Returns:
	//   - *models.Car: Pointer to the created car record with generated fields
	//   - error: Validation error, business rule violation, or data access error
	CreateCar(ctx context.Context, carReq models.CarRequest) (*models.Car, error)

	// UpdateCar modifies an existing car record with business validation.
	// Validates changes against business rules and ensures data consistency.
	// Parameters:
	//   - ctx: Request context for transaction management
	//   - id: Unique identifier of the car to update
	//   - carReq: Updated car data with new field values
	// Returns:
	//   - *models.Car: Pointer to the updated car record
	//   - error: Validation error, business rule violation, or update failure
	UpdateCar(ctx context.Context, id string, carReq models.CarRequest) (*models.Car, error)

	// DeleteCar removes a car record with business rule validation.
	// May enforce cascade rules, audit logging, and referential integrity checks.
	// Parameters:
	//   - ctx: Request context for transaction management
	//   - id: Unique identifier of the car to delete
	// Returns:
	//   - *models.Car: Pointer to the deleted car record (for audit purposes)
	//   - error: Business rule violation or deletion failure
	DeleteCar(ctx context.Context, id string) (*models.Car, error)
	GetAllCars(ctx context.Context) (*[]models.Car, error)
}

// EngineServiceInterface defines the contract for engine business logic operations.
// This interface handles all business operations related to engine entities,
// including technical validation, performance calculations, and business rule enforcement.
// All operations include comprehensive validation and business logic application.
type EngineServiceInterface interface {
	// GetEngineByID retrieves an engine by its unique identifier.
	// May include performance calculations and technical specification validation.
	// Parameters:
	//   - ctx: Request context for cancellation and timeout
	//   - id: Unique identifier of the engine (UUID string format)
	// Returns:
	//   - *models.Engine: Pointer to the engine record if found
	//   - error: Business logic error or data access error
	GetEngineByID(ctx context.Context, id string) (*models.Engine, error)

	// CreateEngine creates a new engine record with technical validation.
	// Validates engine specifications against technical standards and business rules.
	// Parameters:
	//   - ctx: Request context for transaction management
	//   - engineReq: Engine creation request with technical specifications
	// Returns:
	//   - *models.Engine: Pointer to the created engine record
	//   - error: Technical validation error or creation failure
	CreateEngine(ctx context.Context, engineReq models.EngineRequest) (*models.Engine, error)

	// UpdateEngine modifies existing engine specifications with validation.
	// Ensures updated specifications meet technical and business requirements.
	// Parameters:
	//   - ctx: Request context for transaction management
	//   - id: Unique identifier of the engine to update
	//   - engineReq: Updated engine specifications
	// Returns:
	//   - *models.Engine: Pointer to the updated engine record
	//   - error: Validation error or update failure
	UpdateEngine(ctx context.Context, id string, engineReq models.EngineRequest) (*models.Engine, error)

	// DeleteEngine removes an engine record with dependency checks.
	// Validates that engine can be safely deleted without violating constraints.
	// Parameters:
	//   - ctx: Request context for transaction management
	//   - id: Unique identifier of the engine to delete
	// Returns:
	//   - *models.Engine: Pointer to the deleted engine record
	//   - error: Dependency violation or deletion failure
	DeleteEngine(ctx context.Context, id string) (*models.Engine, error)
}

// AuthServiceInterface defines the contract for user authentication and management.
// This interface encapsulates all business operations related to user accounts,
// including registration, authentication, and user data management.
// It ensures security best practices and compliance with authentication standards.
type AuthServiceInterface interface {
	// RegisterUser registers a new user with full validation and security checks.
	// Validates user input, enforces password policies, and coordinates with data persistence.
	// Parameters:
	//   - ctx: Request context for transaction management
	//   - userReq: User registration request with necessary fields
	// Returns:
	//   - error: Validation error, business rule violation, or data access error
	RegisterUser(ctx context.Context, userReq models.UserRequest) error

	// Additional authentication-related methods can be defined here,
	// such as Login, Logout, PasswordReset, etc., following similar patterns.
	LoginUser(ctx context.Context, loginReq models.LoginRequest) (models.User, error)
}
	

