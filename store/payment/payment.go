package payment

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"

	"github.com/PrateekKumar15/CarZone/models"
)

// PaymentStore implements payment data access operations
type PaymentStore struct {
	db *sql.DB
}

// New creates a new PaymentStore instance
func New(db *sql.DB) *PaymentStore {
	return &PaymentStore{db: db}
}

// GetPaymentByID retrieves a payment by its ID
func (s *PaymentStore) GetPaymentByID(ctx context.Context, id string) (models.Payment, error) {
	tracer := otel.Tracer("PaymentStore")
	ctx, span := tracer.Start(ctx, "GetPaymentByID-Store")
	defer span.End()

	var payment models.Payment

	query := `SELECT id, booking_id, razorpay_order_id, razorpay_payment_id, amount, currency, 
	         status, method, transaction_id, description, notes, created_at, updated_at 
	         FROM payment WHERE id = $1`

	row := s.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&payment.ID, &payment.BookingID, &payment.RazorpayOrderID, &payment.RazorpayPaymentID,
		&payment.Amount, &payment.Currency, &payment.Status, &payment.Method, &payment.TransactionID,
		&payment.Description, &payment.Notes, &payment.CreatedAt, &payment.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Payment{}, errors.New("no payment found with the given ID")
		}
		return models.Payment{}, err
	}

	return payment, nil
}

// GetPaymentsByBookingID retrieves all payments for a specific booking
func (s *PaymentStore) GetPaymentsByBookingID(ctx context.Context, bookingID string) ([]models.Payment, error) {
	tracer := otel.Tracer("PaymentStore")
	ctx, span := tracer.Start(ctx, "GetPaymentsByBookingID-Store")
	defer span.End()

	var payments []models.Payment

	query := `SELECT id, booking_id, razorpay_order_id, razorpay_payment_id, amount, currency, 
	         status, method, transaction_id, description, notes, created_at, updated_at 
	         FROM payment WHERE booking_id = $1 ORDER BY created_at DESC`

	rows, err := s.db.QueryContext(ctx, query, bookingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var payment models.Payment
		err = rows.Scan(&payment.ID, &payment.BookingID, &payment.RazorpayOrderID, &payment.RazorpayPaymentID,
			&payment.Amount, &payment.Currency, &payment.Status, &payment.Method, &payment.TransactionID,
			&payment.Description, &payment.Notes, &payment.CreatedAt, &payment.UpdatedAt)

		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	return payments, nil
}

// GetPaymentByRazorpayOrderID retrieves a payment by Razorpay order ID
func (s *PaymentStore) GetPaymentByRazorpayOrderID(ctx context.Context, orderID string) (models.Payment, error) {
	tracer := otel.Tracer("PaymentStore")
	ctx, span := tracer.Start(ctx, "GetPaymentByRazorpayOrderID-Store")
	defer span.End()

	var payment models.Payment

	query := `SELECT id, booking_id, razorpay_order_id, razorpay_payment_id, amount, currency, 
	         status, method, transaction_id, description, notes, created_at, updated_at 
	         FROM payment WHERE razorpay_order_id = $1`

	row := s.db.QueryRowContext(ctx, query, orderID)
	err := row.Scan(&payment.ID, &payment.BookingID, &payment.RazorpayOrderID, &payment.RazorpayPaymentID,
		&payment.Amount, &payment.Currency, &payment.Status, &payment.Method, &payment.TransactionID,
		&payment.Description, &payment.Notes, &payment.CreatedAt, &payment.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Payment{}, errors.New("no payment found with the given Razorpay order ID")
		}
		return models.Payment{}, err
	}

	return payment, nil
}

// CreatePayment creates a new payment record
func (s *PaymentStore) CreatePayment(ctx context.Context, paymentReq models.PaymentRequest) (models.Payment, error) {
	tracer := otel.Tracer("PaymentStore")
	ctx, span := tracer.Start(ctx, "CreatePayment-Store")
	defer span.End()

	var createdPayment models.Payment

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Payment{}, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// Generate new UUID for payment
	paymentId := uuid.New()
	createdAt := time.Now()
	updatedAt := createdAt

	query := `INSERT INTO payment (id, booking_id, amount, currency, status, method, 
	         description, notes, created_at, updated_at)
	         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	         RETURNING id, booking_id, razorpay_order_id, razorpay_payment_id, amount, currency, 
	         status, method, transaction_id, description, notes, created_at, updated_at`

	err = tx.QueryRowContext(ctx, query, paymentId, paymentReq.BookingID, paymentReq.Amount, "INR",
		models.PaymentStatusPending, paymentReq.Method, paymentReq.Description,
		&paymentReq.Notes, createdAt, updatedAt).Scan(
		&createdPayment.ID, &createdPayment.BookingID, &createdPayment.RazorpayOrderID,
		&createdPayment.RazorpayPaymentID, &createdPayment.Amount, &createdPayment.Currency,
		&createdPayment.Status, &createdPayment.Method, &createdPayment.TransactionID,
		&createdPayment.Description, &createdPayment.Notes, &createdPayment.CreatedAt,
		&createdPayment.UpdatedAt)

	if err != nil {
		return models.Payment{}, err
	}

	return createdPayment, nil
}

// UpdatePaymentWithRazorpayDetails updates payment with Razorpay order details
func (s *PaymentStore) UpdatePaymentWithRazorpayDetails(ctx context.Context, paymentID uuid.UUID, orderID string) (models.Payment, error) {
	tracer := otel.Tracer("PaymentStore")
	ctx, span := tracer.Start(ctx, "UpdatePaymentWithRazorpayDetails-Store")
	defer span.End()

	var updatedPayment models.Payment

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Payment{}, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	query := `UPDATE payment SET razorpay_order_id = $1, updated_at = $2 WHERE id = $3 
	         RETURNING id, booking_id, razorpay_order_id, razorpay_payment_id, amount, currency, 
	         status, method, transaction_id, description, notes, created_at, updated_at`

	err = tx.QueryRowContext(ctx, query, orderID, time.Now(), paymentID).Scan(
		&updatedPayment.ID, &updatedPayment.BookingID, &updatedPayment.RazorpayOrderID,
		&updatedPayment.RazorpayPaymentID, &updatedPayment.Amount, &updatedPayment.Currency,
		&updatedPayment.Status, &updatedPayment.Method, &updatedPayment.TransactionID,
		&updatedPayment.Description, &updatedPayment.Notes, &updatedPayment.CreatedAt,
		&updatedPayment.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Payment{}, errors.New("no payment found with the given ID")
		}
		return models.Payment{}, err
	}

	return updatedPayment, nil
}

// UpdatePaymentStatus updates the payment status
func (s *PaymentStore) UpdatePaymentStatus(ctx context.Context, id string, status models.PaymentStatus, paymentID *string, transactionID *string) (models.Payment, error) {
	tracer := otel.Tracer("PaymentStore")
	ctx, span := tracer.Start(ctx, "UpdatePaymentStatus-Store")
	defer span.End()

	// Debug logging
	fmt.Printf("DEBUG: UpdatePaymentStatus called with:\n")
	fmt.Printf("  ID: %s\n", id)
	fmt.Printf("  Status: %s\n", status)
	if paymentID != nil {
		fmt.Printf("  PaymentID: %s\n", *paymentID)
	} else {
		fmt.Printf("  PaymentID: nil\n")
	}
	if transactionID != nil {
		fmt.Printf("  TransactionID: %s\n", *transactionID)
	} else {
		fmt.Printf("  TransactionID: nil\n")
	}

	var updatedPayment models.Payment

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Printf("DEBUG: Failed to begin transaction: %v\n", err)
		return models.Payment{}, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	query := `UPDATE payment SET status = $1, razorpay_payment_id = $2, transaction_id = $3, updated_at = $4 
	         WHERE id = $5 
	         RETURNING id, booking_id, razorpay_order_id, razorpay_payment_id, amount, currency, 
	         status, method, transaction_id, description, notes, created_at, updated_at`

	err = tx.QueryRowContext(ctx, query, status, paymentID, transactionID, time.Now(), id).Scan(
		&updatedPayment.ID, &updatedPayment.BookingID, &updatedPayment.RazorpayOrderID,
		&updatedPayment.RazorpayPaymentID, &updatedPayment.Amount, &updatedPayment.Currency,
		&updatedPayment.Status, &updatedPayment.Method, &updatedPayment.TransactionID,
		&updatedPayment.Description, &updatedPayment.Notes, &updatedPayment.CreatedAt,
		&updatedPayment.UpdatedAt)

	if err != nil {
		fmt.Printf("DEBUG: Failed to execute update query: %v\n", err)
		if err == sql.ErrNoRows {
			return models.Payment{}, errors.New("no payment found with the given ID")
		}
		return models.Payment{}, err
	}

	fmt.Printf("DEBUG: Payment updated successfully:\n")
	fmt.Printf("  ID: %s\n", updatedPayment.ID.String())
	if updatedPayment.RazorpayPaymentID != nil {
		fmt.Printf("  RazorpayPaymentID: %s\n", *updatedPayment.RazorpayPaymentID)
	} else {
		fmt.Printf("  RazorpayPaymentID: nil\n")
	}
	fmt.Printf("  Status: %s\n", updatedPayment.Status)

	return updatedPayment, nil
}

// DeletePayment deletes a payment by ID
func (s *PaymentStore) DeletePayment(ctx context.Context, id string) (models.Payment, error) {
	tracer := otel.Tracer("PaymentStore")
	ctx, span := tracer.Start(ctx, "DeletePayment-Store")
	defer span.End()

	var deletedPayment models.Payment

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Payment{}, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// First get the payment data before deleting
	query := `SELECT id, booking_id, razorpay_order_id, razorpay_payment_id, amount, currency, 
	         status, method, transaction_id, description, notes, created_at, updated_at 
	         FROM payment WHERE id = $1`

	err = tx.QueryRowContext(ctx, query, id).Scan(&deletedPayment.ID, &deletedPayment.BookingID,
		&deletedPayment.RazorpayOrderID, &deletedPayment.RazorpayPaymentID, &deletedPayment.Amount,
		&deletedPayment.Currency, &deletedPayment.Status, &deletedPayment.Method,
		&deletedPayment.TransactionID, &deletedPayment.Description, &deletedPayment.Notes,
		&deletedPayment.CreatedAt, &deletedPayment.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Payment{}, errors.New("no payment found with the given ID")
		}
		return models.Payment{}, err
	}

	// Now delete the payment
	result, err := tx.ExecContext(ctx, "DELETE FROM payment WHERE id = $1", id)
	if err != nil {
		return models.Payment{}, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.Payment{}, err
	}
	if rowsAffected == 0 {
		return models.Payment{}, errors.New("no payment found with the given ID")
	}

	return deletedPayment, nil
}

// GetPaymentsByUserID retrieves all payments for a specific user
func (ps *PaymentStore) GetPaymentsByUserID(ctx context.Context, userID string) ([]models.Payment, error) {
	tracer := otel.Tracer("PaymentStore")
	ctx, span := tracer.Start(ctx, "GetPaymentsByUserID-Store")
	defer span.End()

	if userID == "" {
		return nil, errors.New("user ID cannot be empty")
	}

	// Join payment with booking to get user information
	query := `
		SELECT p.id, p.booking_id, p.razorpay_order_id, p.razorpay_payment_id, p.amount, 
			   p.currency, p.status, p.method, p.transaction_id, p.description,
			   p.notes, p.created_at, p.updated_at
		FROM payment p
		INNER JOIN booking b ON p.booking_id = b.id
		WHERE b.customer_id = $1
		ORDER BY p.created_at DESC`

	rows, err := ps.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var payment models.Payment
		err := rows.Scan(&payment.ID, &payment.BookingID, &payment.RazorpayOrderID,
			&payment.RazorpayPaymentID, &payment.Amount, &payment.Currency, &payment.Status,
			&payment.Method, &payment.TransactionID, &payment.Description,
			&payment.Notes, &payment.CreatedAt, &payment.UpdatedAt)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}

// GetAllPayments retrieves all payment records
func (ps *PaymentStore) GetAllPayments(ctx context.Context) ([]models.Payment, error) {
	tracer := otel.Tracer("PaymentStore")
	ctx, span := tracer.Start(ctx, "GetAllPayments-Store")
	defer span.End()

	query := `
		SELECT p.id, p.booking_id, p.razorpay_order_id, p.razorpay_payment_id, p.amount, 
			   p.currency, p.status, p.method, p.transaction_id, p.description,
			   p.notes, p.created_at, p.updated_at
		FROM payment p
		ORDER BY p.created_at DESC`

	rows, err := ps.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var payment models.Payment
		err := rows.Scan(&payment.ID, &payment.BookingID, &payment.RazorpayOrderID,
			&payment.RazorpayPaymentID, &payment.Amount, &payment.Currency, &payment.Status,
			&payment.Method, &payment.TransactionID, &payment.Description,
			&payment.Notes, &payment.CreatedAt, &payment.UpdatedAt)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}
