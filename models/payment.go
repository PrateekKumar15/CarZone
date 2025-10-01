package models

import (
	"time"

	"github.com/google/uuid"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

// PaymentMethod represents the payment method used
type PaymentMethod string

const (
	PaymentMethodRazorpay   PaymentMethod = "razorpay"
	PaymentMethodCash       PaymentMethod = "cash"
	PaymentMethodCard       PaymentMethod = "card"
	PaymentMethodUPI        PaymentMethod = "upi"
	PaymentMethodNetbanking PaymentMethod = "netbanking"
)

// Payment represents a payment record in the database
type Payment struct {
	ID                uuid.UUID     `json:"id" db:"id"`
	BookingID         uuid.UUID     `json:"booking_id" db:"booking_id"`
	RazorpayOrderID   *string       `json:"razorpay_order_id,omitempty" db:"razorpay_order_id"`
	RazorpayPaymentID *string       `json:"razorpay_payment_id,omitempty" db:"razorpay_payment_id"`
	Amount            float64       `json:"amount" db:"amount"`     // Amount in INR
	Currency          string        `json:"currency" db:"currency"` // INR
	Status            PaymentStatus `json:"status" db:"status"`
	Method            PaymentMethod `json:"method" db:"method"`
	TransactionID     *string       `json:"transaction_id,omitempty" db:"transaction_id"`
	Description       string        `json:"description" db:"description"`
	Notes             *string       `json:"notes,omitempty" db:"notes"`
	CreatedAt         time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at" db:"updated_at"`
}

// PaymentRequest represents the request to create a payment
type PaymentRequest struct {
	BookingID   uuid.UUID     `json:"booking_id" validate:"required"`
	Amount      float64       `json:"amount" validate:"required,gt=0"`
	Method      PaymentMethod `json:"method" validate:"required"`
	Description string        `json:"description"`
	Notes       string        `json:"notes,omitempty"`
}

// RazorpayOrderRequest represents the request to create a Razorpay order
type RazorpayOrderRequest struct {
	Amount   int    `json:"amount"`   // Amount in paise (smallest currency unit)
	Currency string `json:"currency"` // INR
	Receipt  string `json:"receipt"`  // Unique receipt ID
}

// RazorpayOrderResponse represents the response from Razorpay order creation
type RazorpayOrderResponse struct {
	ID       string `json:"id"`
	Entity   string `json:"entity"`
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
	Receipt  string `json:"receipt"`
	Status   string `json:"status"`
}

// PaymentVerificationRequest represents the request to verify a payment
type PaymentVerificationRequest struct {
	RazorpayOrderID   string `json:"razorpay_order_id" validate:"required"`
	RazorpayPaymentID string `json:"razorpay_payment_id" validate:"required"`
	RazorpaySignature string `json:"razorpay_signature" validate:"required"`
}
