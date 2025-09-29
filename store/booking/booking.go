package booking

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/PrateekKumar15/CarZone/models"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

type BookingStore struct {
	db *sql.DB
}

func New(db *sql.DB) BookingStore {
	return BookingStore{db: db}
}

func (s BookingStore) GetBookingByID(ctx context.Context, id string) (models.Booking, error) {
	tracer := otel.Tracer("BookingStore")
	ctx, span := tracer.Start(ctx, "GetBookingByID-Store")
	defer span.End()

	var booking models.Booking

	query := `SELECT id, customer_id, car_id, owner_id, booking_type, status, total_amount, 
	         start_date, end_date, notes, created_at, updated_at 
	         FROM booking WHERE id = $1`

	row := s.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&booking.ID, &booking.CustomerID, &booking.CarID, &booking.OwnerID,
		&booking.BookingType, &booking.Status, &booking.TotalAmount, &booking.StartDate,
		&booking.EndDate, &booking.Notes, &booking.CreatedAt, &booking.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Booking{}, errors.New("no booking found with the given ID")
		}
		return models.Booking{}, err
	}

	return booking, nil
}

func (s BookingStore) GetBookingsByCustomerID(ctx context.Context, customerID string) ([]models.Booking, error) {
	tracer := otel.Tracer("BookingStore")
	ctx, span := tracer.Start(ctx, "GetBookingsByCustomerID-Store")
	defer span.End()

	var bookings []models.Booking

	query := `SELECT id, customer_id, car_id, owner_id, booking_type, status, total_amount, 
	         start_date, end_date, notes, created_at, updated_at 
	         FROM booking WHERE customer_id = $1 ORDER BY created_at DESC`

	rows, err := s.db.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var booking models.Booking
		err = rows.Scan(&booking.ID, &booking.CustomerID, &booking.CarID, &booking.OwnerID,
			&booking.BookingType, &booking.Status, &booking.TotalAmount, &booking.StartDate,
			&booking.EndDate, &booking.Notes, &booking.CreatedAt, &booking.UpdatedAt)

		if err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	return bookings, nil
}

func (s BookingStore) GetBookingsByCarID(ctx context.Context, carID string) ([]models.Booking, error) {
	tracer := otel.Tracer("BookingStore")
	ctx, span := tracer.Start(ctx, "GetBookingsByCarID-Store")
	defer span.End()

	var bookings []models.Booking

	query := `SELECT id, customer_id, car_id, owner_id, booking_type, status, total_amount, 
	         start_date, end_date, notes, created_at, updated_at 
	         FROM booking WHERE car_id = $1 ORDER BY created_at DESC`

	rows, err := s.db.QueryContext(ctx, query, carID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var booking models.Booking
		err = rows.Scan(&booking.ID, &booking.CustomerID, &booking.CarID, &booking.OwnerID,
			&booking.BookingType, &booking.Status, &booking.TotalAmount, &booking.StartDate,
			&booking.EndDate, &booking.Notes, &booking.CreatedAt, &booking.UpdatedAt)

		if err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	return bookings, nil
}

func (s BookingStore) GetBookingsByOwnerID(ctx context.Context, ownerID string) ([]models.Booking, error) {
	tracer := otel.Tracer("BookingStore")
	ctx, span := tracer.Start(ctx, "GetBookingsByOwnerID-Store")
	defer span.End()

	var bookings []models.Booking

	query := `SELECT id, customer_id, car_id, owner_id, booking_type, status, total_amount, 
	         start_date, end_date, notes, created_at, updated_at 
	         FROM booking WHERE owner_id = $1 ORDER BY created_at DESC`

	rows, err := s.db.QueryContext(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var booking models.Booking
		err = rows.Scan(&booking.ID, &booking.CustomerID, &booking.CarID, &booking.OwnerID,
			&booking.BookingType, &booking.Status, &booking.TotalAmount, &booking.StartDate,
			&booking.EndDate, &booking.Notes, &booking.CreatedAt, &booking.UpdatedAt)

		if err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	return bookings, nil
}

func (s BookingStore) CreateBooking(ctx context.Context, bookingReq models.BookingRequest, totalAmount float64) (models.Booking, error) {
	tracer := otel.Tracer("BookingStore")
	ctx, span := tracer.Start(ctx, "CreateBooking-Store")
	defer span.End()

	var createdBooking models.Booking

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Booking{}, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// Generate new UUID for booking
	bookingId := uuid.New()
	createdAt := time.Now()
	updatedAt := createdAt

	query := `INSERT INTO booking (id, customer_id, car_id, owner_id, booking_type, status, total_amount, 
	         start_date, end_date, notes, created_at, updated_at)
	         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	         RETURNING id, customer_id, car_id, owner_id, booking_type, status, total_amount, 
	         start_date, end_date, notes, created_at, updated_at`

	err = tx.QueryRowContext(ctx, query, bookingId, bookingReq.CustomerID, bookingReq.CarID,
		bookingReq.OwnerID, bookingReq.BookingType, models.BookingStatusPending, totalAmount,
		bookingReq.StartDate, bookingReq.EndDate, bookingReq.Notes, createdAt, updatedAt).Scan(
		&createdBooking.ID, &createdBooking.CustomerID, &createdBooking.CarID, &createdBooking.OwnerID,
		&createdBooking.BookingType, &createdBooking.Status, &createdBooking.TotalAmount,
		&createdBooking.StartDate, &createdBooking.EndDate, &createdBooking.Notes,
		&createdBooking.CreatedAt, &createdBooking.UpdatedAt)

	if err != nil {
		return models.Booking{}, err
	}

	return createdBooking, nil
}

func (s BookingStore) UpdateBookingStatus(ctx context.Context, id string, status models.BookingStatus) (models.Booking, error) {
	tracer := otel.Tracer("BookingStore")
	ctx, span := tracer.Start(ctx, "UpdateBookingStatus-Store")
	defer span.End()

	var updatedBooking models.Booking

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Booking{}, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	query := `UPDATE booking SET status = $1, updated_at = $2 WHERE id = $3 
	         RETURNING id, customer_id, car_id, owner_id, booking_type, status, total_amount, 
	         start_date, end_date, notes, created_at, updated_at`

	err = tx.QueryRowContext(ctx, query, status, time.Now(), id).Scan(
		&updatedBooking.ID, &updatedBooking.CustomerID, &updatedBooking.CarID, &updatedBooking.OwnerID,
		&updatedBooking.BookingType, &updatedBooking.Status, &updatedBooking.TotalAmount,
		&updatedBooking.StartDate, &updatedBooking.EndDate, &updatedBooking.Notes,
		&updatedBooking.CreatedAt, &updatedBooking.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Booking{}, errors.New("no booking found with the given ID")
		}
		return models.Booking{}, err
	}

	return updatedBooking, nil
}

func (s BookingStore) DeleteBooking(ctx context.Context, id string) (models.Booking, error) {
	tracer := otel.Tracer("BookingStore")
	ctx, span := tracer.Start(ctx, "DeleteBooking-Store")
	defer span.End()

	var deletedBooking models.Booking

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Booking{}, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// First get the booking data before deleting
	query := `SELECT id, customer_id, car_id, owner_id, booking_type, status, total_amount, 
	         start_date, end_date, notes, created_at, updated_at 
	         FROM booking WHERE id = $1`

	err = tx.QueryRowContext(ctx, query, id).Scan(&deletedBooking.ID, &deletedBooking.CustomerID,
		&deletedBooking.CarID, &deletedBooking.OwnerID, &deletedBooking.BookingType, &deletedBooking.Status,
		&deletedBooking.TotalAmount, &deletedBooking.StartDate, &deletedBooking.EndDate,
		&deletedBooking.Notes, &deletedBooking.CreatedAt, &deletedBooking.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Booking{}, errors.New("no booking found with the given ID")
		}
		return models.Booking{}, err
	}

	// Now delete the booking
	result, err := tx.ExecContext(ctx, "DELETE FROM booking WHERE id = $1", id)
	if err != nil {
		return models.Booking{}, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.Booking{}, err
	}
	if rowsAffected == 0 {
		return models.Booking{}, errors.New("no booking found with the given ID")
	}

	return deletedBooking, nil
}

func (s BookingStore) GetAllBookings(ctx context.Context) ([]models.Booking, error) {
	tracer := otel.Tracer("BookingStore")
	ctx, span := tracer.Start(ctx, "GetAllBookings-Store")
	defer span.End()

	var bookings []models.Booking

	query := `SELECT id, customer_id, car_id, owner_id, booking_type, status, total_amount, 
	         start_date, end_date, notes, created_at, updated_at 
	         FROM booking ORDER BY created_at DESC`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var booking models.Booking
		err = rows.Scan(&booking.ID, &booking.CustomerID, &booking.CarID, &booking.OwnerID,
			&booking.BookingType, &booking.Status, &booking.TotalAmount, &booking.StartDate,
			&booking.EndDate, &booking.Notes, &booking.CreatedAt, &booking.UpdatedAt)

		if err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	return bookings, nil
}
