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

	// GetCarWithOwnerByID retrieves a single car record with owner information by its unique identifier.
	// Parameters:
	//   - ctx: Request context for cancellation and timeout
	//   - id: Unique identifier of the car (UUID string format)
	// Returns:
	//   - models.Car: The car record with populated owner field if found
	//   - error: Error if car not found or database operation fails
	GetCarWithOwnerByID(ctx context.Context, id string) (models.Car, error)

	// GetCarByBrand retrieves multiple car records filtered by brand name.
	// Parameters:
	//   - ctx: Request context for cancellation and timeout
	//   - brand: Brand name to filter by (e.g., "Toyota", "BMW")
	// Returns:
	//   - []models.Car: Slice of car records matching the brand
	//   - error: Error if database operation fails
	GetCarByBrand(ctx context.Context, brand string) ([]models.Car, error)

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

	GetAllCars(ctx context.Context) ([]models.Car, error)
}

// UserStoreInterface defines the contract for user authentication and management operations.
// This interface abstracts all database operations related to user entities,
// following the Repository pattern to decouple business logic from data persistence.
// All methods accept a context for request scoping, cancellation, and timeout handling.
type UserStoreInterface interface {
	// CreateUser inserts a new user record into the database.
	// The method generates a new UUID for the user and handles all creation logic.
	// Parameters:
	//   - ctx: Request context for transaction management
	//   - userReq: User data to be inserted (without ID, timestamps)
	// Returns:
	//   - error: Error if creation fails or validation errors occur
	CreateUser(ctx context.Context, userReq models.UserRequest) error

	// GetUser retrieves a user by email and validates password for authentication.
	// Parameters:
	//   - ctx: Request context for cancellation and timeout
	//   - email: User's email address
	//   - password: Plain text password for validation
	// Returns:
	//   - models.User: User record if authentication successful
	//   - error: Error if user not found or password invalid
	GetUser(ctx context.Context, email, password string) (models.User, error)

	// GetUserByID retrieves a user by their unique ID.
	// Parameters:
	//   - ctx: Request context for cancellation and timeout
	//   - userID: User's unique identifier (UUID)
	// Returns:
	//   - models.User: User record if found
	//   - error: Error if user not found or database operation fails
	GetUserByID(ctx context.Context, userID string) (models.User, error)

	// UpdateUser modifies an existing user record.
	// Parameters:
	//   - ctx: Request context for transaction management
	//   - id: User's unique identifier
	//   - userReq: Updated user data
	// Returns:
	//   - models.User: Updated user record
	//   - error: Error if user not found or update fails
	UpdateUser(ctx context.Context, id string, userReq models.UserRequest) (models.User, error)

	// UpdateProfileData updates only the profile_data field for a user.
	// Parameters:
	//   - ctx: Request context for transaction management
	//   - userID: User's unique identifier
	//   - profileData: Profile data as map[string]interface{}
	// Returns:
	//   - error: Error if user not found or update fails
	UpdateProfileData(ctx context.Context, userID string, profileData map[string]interface{}) error

	// DeleteUser removes a user record from the database.
	// Parameters:
	//   - ctx: Request context for transaction management
	//   - id: User's unique identifier
	// Returns:
	//   - models.User: Deleted user record for audit purposes
	//   - error: Error if user not found or deletion fails
	DeleteUser(ctx context.Context, id string) (models.User, error)

	// GetAllUsers retrieves all user records from the database.
	// Parameters:
	//   - ctx: Request context for cancellation and timeout
	// Returns:
	//   - []models.User: Slice of all user records
	//   - error: Error if database operation fails
	GetAllUsers(ctx context.Context) ([]models.User, error)

	// GetUsersByRole retrieves all users with a specific role.
	// Parameters:
	//   - ctx: Request context for cancellation and timeout
	//   - role: User role to filter by (owner, renter, admin)
	// Returns:
	//   - []models.User: Slice of users with specified role
	//   - error: Error if database operation fails
	GetUsersByRole(ctx context.Context, role string) ([]models.User, error)
}

// BookingStoreInterface defines the contract for booking data access operations.
// This interface abstracts all database operations related to booking entities,
// following the Repository pattern to decouple business logic from data persistence.
type BookingStoreInterface interface {
	// GetBookingByID retrieves a single booking record by its unique identifier.
	// Parameters:
	//   - ctx: Request context for cancellation and timeout
	//   - id: Unique identifier of the booking (UUID string format)
	// Returns:
	//   - models.Booking: The booking record if found
	//   - error: Error if booking not found or database operation fails
	GetBookingByID(ctx context.Context, id string) (models.Booking, error)

	// GetBookingsByCustomerID retrieves all bookings for a specific customer.
	// Parameters:
	//   - ctx: Request context for cancellation and timeout
	//   - customerID: Customer's unique identifier
	// Returns:
	//   - []models.Booking: Slice of booking records for the customer
	//   - error: Error if database operation fails
	GetBookingsByCustomerID(ctx context.Context, customerID string) ([]models.Booking, error)

	// GetBookingsByCarID retrieves all bookings for a specific car.
	// Parameters:
	//   - ctx: Request context for cancellation and timeout
	//   - carID: Car's unique identifier
	// Returns:
	//   - []models.Booking: Slice of booking records for the car
	//   - error: Error if database operation fails
	GetBookingsByCarID(ctx context.Context, carID string) ([]models.Booking, error)

	// GetBookingsByOwnerID retrieves all bookings for cars owned by a specific owner.
	// Parameters:
	//   - ctx: Request context for cancellation and timeout
	//   - ownerID: Owner's unique identifier
	// Returns:
	//   - []models.Booking: Slice of booking records for the owner's cars
	//   - error: Error if database operation fails
	GetBookingsByOwnerID(ctx context.Context, ownerID string) ([]models.Booking, error)

	// CreateBooking inserts a new booking record into the database.
	// Parameters:
	//   - ctx: Request context for transaction management
	//   - bookingReq: Booking data to be inserted
	// Returns:
	//   - models.Booking: The created booking record with generated ID and timestamps
	//   - error: Error if creation fails or validation errors occur
	CreateBooking(ctx context.Context, bookingReq models.BookingRequest, totalAmount float64) (models.Booking, error)

	// UpdateBookingStatus updates the status of an existing booking.
	// Parameters:
	//   - ctx: Request context for transaction management
	//   - id: Unique identifier of the booking to update
	//   - status: New booking status
	// Returns:
	//   - models.Booking: The updated booking record
	//   - error: Error if booking not found or update operation fails
	UpdateBookingStatus(ctx context.Context, id string, status models.BookingStatus) (models.Booking, error)

	// DeleteBooking removes a booking record from the database.
	// Parameters:
	//   - ctx: Request context for transaction management
	//   - id: Unique identifier of the booking to delete
	// Returns:
	//   - models.Booking: The deleted booking record
	//   - error: Error if booking not found or deletion fails
	DeleteBooking(ctx context.Context, id string) (models.Booking, error)

	// GetAllBookings retrieves all booking records.
	// Parameters:
	//   - ctx: Request context for cancellation and timeout
	// Returns:
	//   - []models.Booking: Slice of all booking records
	//   - error: Error if database operation fails
	GetAllBookings(ctx context.Context) ([]models.Booking, error)
}
