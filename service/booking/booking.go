package booking

import (
	"context"
	"errors"
	"time"

	"github.com/PrateekKumar15/CarZone/models"
	"github.com/PrateekKumar15/CarZone/store"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

type BookingService struct {
	bookingStore store.BookingStoreInterface
	carStore     store.CarStoreInterface
}

func NewBookingService(bookingStore store.BookingStoreInterface, carStore store.CarStoreInterface) *BookingService {
	return &BookingService{
		bookingStore: bookingStore,
		carStore:     carStore,
	}
}

func (s *BookingService) GetBookingByID(ctx context.Context, id string) (*models.Booking, error) {
	tracer := otel.Tracer("BookingService")
	ctx, span := tracer.Start(ctx, "GetBookingByID-Service")
	defer span.End()

	if id == "" {
		return nil, errors.New("booking ID cannot be empty")
	}

	booking, err := s.bookingStore.GetBookingByID(ctx, id)
	if err != nil {
		if err.Error() == "no booking found with the given ID" {
			return nil, nil
		}
		return nil, err
	}

	return &booking, nil
}

func (s *BookingService) GetBookingsByCustomerID(ctx context.Context, customerID string) (*[]models.Booking, error) {
	tracer := otel.Tracer("BookingService")
	ctx, span := tracer.Start(ctx, "GetBookingsByCustomerID-Service")
	defer span.End()

	if customerID == "" {
		return nil, errors.New("customer ID cannot be empty")
	}

	bookings, err := s.bookingStore.GetBookingsByCustomerID(ctx, customerID)
	if err != nil {
		return nil, err
	}

	return &bookings, nil
}

func (s *BookingService) GetBookingsByCarID(ctx context.Context, carID string) (*[]models.Booking, error) {
	tracer := otel.Tracer("BookingService")
	ctx, span := tracer.Start(ctx, "GetBookingsByCarID-Service")
	defer span.End()

	if carID == "" {
		return nil, errors.New("car ID cannot be empty")
	}

	bookings, err := s.bookingStore.GetBookingsByCarID(ctx, carID)
	if err != nil {
		return nil, err
	}

	return &bookings, nil
}

func (s *BookingService) GetBookingsByOwnerID(ctx context.Context, ownerID string) (*[]models.Booking, error) {
	tracer := otel.Tracer("BookingService")
	ctx, span := tracer.Start(ctx, "GetBookingsByOwnerID-Service")
	defer span.End()

	if ownerID == "" {
		return nil, errors.New("owner ID cannot be empty")
	}

	bookings, err := s.bookingStore.GetBookingsByOwnerID(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	return &bookings, nil
}

func (s *BookingService) CreateBooking(ctx context.Context, bookingReq models.BookingRequest) (*models.Booking, error) {
	tracer := otel.Tracer("BookingService")
	ctx, span := tracer.Start(ctx, "CreateBooking-Service")
	defer span.End()

	// Validate booking request
	if err := s.validateBookingRequest(bookingReq); err != nil {
		return nil, err
	}

	// Verify car exists and is available
	car, err := s.carStore.GetCarByID(ctx, bookingReq.CarID.String())
	if err != nil {
		return nil, errors.New("failed to verify car availability")
	}

	if car.ID.String() == "00000000-0000-0000-0000-000000000000" {
		return nil, errors.New("car not found")
	}

	if !car.IsAvailable {
		return nil, errors.New("car is not available for booking")
	}

	// Verify owner ID matches the car's owner
	if car.OwnerID == nil || *car.OwnerID != bookingReq.OwnerID {
		return nil, errors.New("owner ID does not match car owner")
	}

	// Check for booking conflicts if it's a rental
	if bookingReq.BookingType == models.BookingTypeRental {
		if err := s.checkBookingConflicts(ctx, bookingReq); err != nil {
			return nil, err
		}
	}

	// Calculate total amount based on booking type and duration
	totalAmount, err := s.calculateTotalAmount(car, bookingReq)
	if err != nil {
		return nil, err
	}

	booking, err := s.bookingStore.CreateBooking(ctx, bookingReq, totalAmount)
	if err != nil {
		return nil, err
	}

	return &booking, nil
}

func (s *BookingService) calculateTotalAmount(car models.Car, bookingReq models.BookingRequest) (float64, error) {
	if bookingReq.BookingType == models.BookingTypePurchase {
		// For purchases, use the sale price
		if car.Price.SalePrice == nil {
			return 0, errors.New("sale price not available for this car")
		}
		return *car.Price.SalePrice, nil
	}

	// For rentals, calculate based on daily rate and duration
	if bookingReq.StartDate == nil || bookingReq.EndDate == nil {
		return 0, errors.New("start and end dates are required for rental bookings")
	}

	dailyRate := car.Price.RentalPriceDaily
	if dailyRate <= 0 {
		return 0, errors.New("invalid daily rental price for this car")
	}

	// Calculate duration in days
	duration := bookingReq.EndDate.Sub(*bookingReq.StartDate)
	days := int(duration.Hours() / 24)
	if days < 1 {
		days = 1 // Minimum 1 day
	}

	totalAmount := dailyRate * float64(days)
	return totalAmount, nil
}

func (s *BookingService) UpdateBookingStatus(ctx context.Context, id string, status models.BookingStatus) (*models.Booking, error) {
	tracer := otel.Tracer("BookingService")
	ctx, span := tracer.Start(ctx, "UpdateBookingStatus-Service")
	defer span.End()

	if id == "" {
		return nil, errors.New("booking ID cannot be empty")
	}

	// Validate status
	if err := s.validateBookingStatus(status); err != nil {
		return nil, err
	}

	// Get current booking to validate status transition
	currentBooking, err := s.bookingStore.GetBookingByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate status transition
	if err := s.validateStatusTransition(currentBooking.Status, status); err != nil {
		return nil, err
	}

	booking, err := s.bookingStore.UpdateBookingStatus(ctx, id, status)
	if err != nil {
		return nil, err
	}

	return &booking, nil
}

func (s *BookingService) DeleteBooking(ctx context.Context, id string) (*models.Booking, error) {
	tracer := otel.Tracer("BookingService")
	ctx, span := tracer.Start(ctx, "DeleteBooking-Service")
	defer span.End()

	if id == "" {
		return nil, errors.New("booking ID cannot be empty")
	}

	// Get booking to check if it can be deleted
	booking, err := s.bookingStore.GetBookingByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Business rule: Only pending or cancelled bookings can be deleted
	if booking.Status != models.BookingStatusPending && booking.Status != models.BookingStatusCancelled {
		return nil, errors.New("only pending or cancelled bookings can be deleted")
	}

	deletedBooking, err := s.bookingStore.DeleteBooking(ctx, id)
	if err != nil {
		return nil, err
	}

	return &deletedBooking, nil
}

func (s *BookingService) GetAllBookings(ctx context.Context) (*[]models.Booking, error) {
	tracer := otel.Tracer("BookingService")
	ctx, span := tracer.Start(ctx, "GetAllBookings-Service")
	defer span.End()

	bookings, err := s.bookingStore.GetAllBookings(ctx)
	if err != nil {
		return nil, err
	}

	return &bookings, nil
}

// validateBookingRequest validates the booking request
func (s *BookingService) validateBookingRequest(req models.BookingRequest) error {
	if req.CustomerID == uuid.Nil {
		return errors.New("customer ID is required")
	}

	if req.CarID == uuid.Nil {
		return errors.New("car ID is required")
	}

	if req.OwnerID == uuid.Nil {
		return errors.New("owner ID is required")
	}

	if req.BookingType == "" {
		return errors.New("booking type is required")
	}

	if req.BookingType != models.BookingTypeRental && req.BookingType != models.BookingTypePurchase {
		return errors.New("booking type must be 'rental' or 'purchase'")
	}

	// Validate rental-specific fields
	if req.BookingType == models.BookingTypeRental {
		return s.validateRentalRequest(req)
	}

	return nil
}

// validateRentalRequest validates rental-specific fields
func (s *BookingService) validateRentalRequest(req models.BookingRequest) error {
	if req.StartDate == nil {
		return errors.New("start date is required for rental bookings")
	}

	if req.EndDate == nil {
		return errors.New("end date is required for rental bookings")
	}

	// Validate date logic
	if req.StartDate.After(*req.EndDate) {
		return errors.New("start date cannot be after end date")
	}

	if req.StartDate.Before(time.Now().Add(-24 * time.Hour)) {
		return errors.New("start date cannot be in the past")
	}

	// Validate minimum rental duration (at least 1 day)
	duration := req.EndDate.Sub(*req.StartDate)
	if duration < 24*time.Hour {
		return errors.New("minimum rental duration is 1 day")
	}

	return nil
}

// validateBookingStatus validates booking status values
func (s *BookingService) validateBookingStatus(status models.BookingStatus) error {
	validStatuses := []models.BookingStatus{
		models.BookingStatusPending,
		models.BookingStatusConfirmed,
		models.BookingStatusCompleted,
		models.BookingStatusCancelled,
	}

	for _, validStatus := range validStatuses {
		if status == validStatus {
			return nil
		}
	}

	return errors.New("invalid booking status")
}

// validateStatusTransition validates if a status transition is allowed
func (s *BookingService) validateStatusTransition(current, new models.BookingStatus) error {
	// Define allowed status transitions
	allowedTransitions := map[models.BookingStatus][]models.BookingStatus{
		models.BookingStatusPending: {
			models.BookingStatusConfirmed,
			models.BookingStatusCancelled,
		},
		models.BookingStatusConfirmed: {
			models.BookingStatusCompleted,
			models.BookingStatusCancelled,
		},
		models.BookingStatusCompleted: {}, // Terminal state
		models.BookingStatusCancelled: {}, // Terminal state
	}

	allowed, exists := allowedTransitions[current]
	if !exists {
		return errors.New("invalid current booking status")
	}

	for _, status := range allowed {
		if status == new {
			return nil
		}
	}

	return errors.New("invalid status transition from " + string(current) + " to " + string(new))
}

// checkBookingConflicts checks for conflicting bookings for rental requests
func (s *BookingService) checkBookingConflicts(ctx context.Context, req models.BookingRequest) error {
	// Get existing bookings for the car
	existingBookings, err := s.bookingStore.GetBookingsByCarID(ctx, req.CarID.String())
	if err != nil {
		return errors.New("failed to check booking conflicts")
	}

	// Check for date conflicts with confirmed/active rentals
	for _, booking := range existingBookings {
		if booking.BookingType == models.BookingTypeRental &&
			(booking.Status == models.BookingStatusConfirmed || booking.Status == models.BookingStatusPending) &&
			booking.StartDate != nil && booking.EndDate != nil {

			// Check if dates overlap
			if s.datesOverlap(*req.StartDate, *req.EndDate, *booking.StartDate, *booking.EndDate) {
				return errors.New("booking conflicts with existing rental for the same period")
			}
		}
	}

	return nil
}

// datesOverlap checks if two date ranges overlap
func (s *BookingService) datesOverlap(start1, end1, start2, end2 time.Time) bool {
	return start1.Before(end2) && end1.After(start2)
}
