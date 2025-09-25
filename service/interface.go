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

	// GetCarByBrand retrieves multiple cars filtered by brand name.
	// Applies business rules for data filtering and presentation logic.
	// Parameters:
	//   - ctx: Request context for cancellation and timeout
	//   - brand: Brand name to filter by (case-sensitive)
	// Returns:
	//   - *[]models.Car: Pointer to slice of car records matching the criteria
	//   - error: Business logic error or data access error
	GetCarByBrand(ctx context.Context, brand string) (*[]models.Car, error)

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

// AuthServiceInterface defines the contract for user authentication and management.
// This interface encapsulates all business operations related to user accounts,
// including registration, authentication, and user data management.
// It ensures security best practices and compliance with authentication standards.
type AuthServiceInterface interface {
	// RegisterUser registers a new user with full validation and security checks.
	// Validates user input, enforces password policies, and coordinates with data persistence.
	// This includes validation of phone number format and role restrictions.
	// Parameters:
	//   - ctx: Request context for transaction management
	//   - userReq: User registration request with necessary fields (email, password, username, phone, role)
	// Returns:
	//   - error: Validation error, business rule violation, or data access error
	RegisterUser(ctx context.Context, userReq models.UserRequest) error

	// LoginUser authenticates a user with email and password credentials.
	// Validates credentials against stored user data and returns complete user profile.
	// Parameters:
	//   - ctx: Request context for cancellation and timeout
	//   - loginReq: Login request with email and password
	// Returns:
	//   - models.User: Complete user record including phone, role, and profile_data
	//   - error: Authentication error or data access error
	LoginUser(ctx context.Context, loginReq models.LoginRequest) (models.User, error)
}

// BookingServiceInterface defines the contract for booking business logic operations.
// This interface abstracts all business operations related to booking entities,
// including validation, business rule enforcement, and coordination with the data layer.
type BookingServiceInterface interface {
	// GetBookingByID retrieves a booking by its unique identifier with business logic applied.
	// Parameters:
	//   - ctx: Request context for cancellation, timeout, and request scoping
	//   - id: Unique identifier of the booking (UUID string format)
	// Returns:
	//   - *models.Booking: Pointer to the booking record if found, nil if not found
	//   - error: Business logic error or underlying data access error
	GetBookingByID(ctx context.Context, id string) (*models.Booking, error)

	// GetBookingsByCustomerID retrieves all bookings for a specific customer.
	// Parameters:
	//   - ctx: Request context for cancellation and timeout
	//   - customerID: Customer's unique identifier
	// Returns:
	//   - *[]models.Booking: Pointer to slice of booking records for the customer
	//   - error: Business logic error or data access error
	GetBookingsByCustomerID(ctx context.Context, customerID string) (*[]models.Booking, error)

	// GetBookingsByCarID retrieves all bookings for a specific car.
	// Parameters:
	//   - ctx: Request context for cancellation and timeout
	//   - carID: Car's unique identifier
	// Returns:
	//   - *[]models.Booking: Pointer to slice of booking records for the car
	//   - error: Business logic error or data access error
	GetBookingsByCarID(ctx context.Context, carID string) (*[]models.Booking, error)

	// GetBookingsByOwnerID retrieves all bookings for cars owned by a specific owner.
	// Parameters:
	//   - ctx: Request context for cancellation and timeout
	//   - ownerID: Owner's unique identifier
	// Returns:
	//   - *[]models.Booking: Pointer to slice of booking records for the owner's cars
	//   - error: Business logic error or data access error
	GetBookingsByOwnerID(ctx context.Context, ownerID string) (*[]models.Booking, error)

	// CreateBooking creates a new booking with full business validation.
	// Validates input data, enforces business rules, and coordinates with data persistence.
	// Parameters:
	//   - ctx: Request context for transaction management
	//   - bookingReq: Booking creation request with all required fields
	// Returns:
	//   - *models.Booking: Pointer to the created booking record with generated fields
	//   - error: Validation error, business rule violation, or data access error
	CreateBooking(ctx context.Context, bookingReq models.BookingRequest) (*models.Booking, error)

	// UpdateBookingStatus modifies booking status with business validation.
	// Validates status transitions and enforces business rules.
	// Parameters:
	//   - ctx: Request context for transaction management
	//   - id: Unique identifier of the booking to update
	//   - status: New booking status
	// Returns:
	//   - *models.Booking: Pointer to the updated booking record
	//   - error: Validation error, business rule violation, or update failure
	UpdateBookingStatus(ctx context.Context, id string, status models.BookingStatus) (*models.Booking, error)

	// DeleteBooking removes a booking record with business rule validation.
	// Parameters:
	//   - ctx: Request context for transaction management
	//   - id: Unique identifier of the booking to delete
	// Returns:
	//   - *models.Booking: Pointer to the deleted booking record (for audit purposes)
	//   - error: Business rule violation or deletion failure
	DeleteBooking(ctx context.Context, id string) (*models.Booking, error)

	// GetAllBookings retrieves all booking records with business filtering.
	// Parameters:
	//   - ctx: Request context for cancellation and timeout
	// Returns:
	//   - *[]models.Booking: Pointer to slice of all booking records
	//   - error: Business logic error or data access error
	GetAllBookings(ctx context.Context) (*[]models.Booking, error)
}
