package models

import (
	"time"

	"github.com/google/uuid"
)

// BookingType represents the type of booking
type BookingType string

const (
	BookingTypeRental   BookingType = "rental"
	BookingTypePurchase BookingType = "purchase"
)

// BookingStatus represents the current status of a booking
type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCompleted BookingStatus = "completed"
	BookingStatusCancelled BookingStatus = "cancelled"
)

// Booking represents a car booking in the system
type Booking struct {
	ID          uuid.UUID     `json:"id"`
	CustomerID  uuid.UUID     `json:"customer_id"`
	CarID       uuid.UUID     `json:"car_id"`
	OwnerID     uuid.UUID     `json:"owner_id"`
	BookingType BookingType   `json:"booking_type"`
	Status      BookingStatus `json:"status"`
	TotalAmount float64       `json:"total_amount"`
	StartDate   *time.Time    `json:"start_date"` // For rentals
	EndDate     *time.Time    `json:"end_date"`   // For rentals
	Notes       string        `json:"notes"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// BookingRequest represents the payload to create a booking
type BookingRequest struct {
	CustomerID  uuid.UUID   `json:"customer_id"`
	CarID       uuid.UUID   `json:"car_id"`
	OwnerID     uuid.UUID   `json:"owner_id"`
	BookingType BookingType `json:"booking_type"`
	StartDate   *time.Time  `json:"start_date"` // For rentals
	EndDate     *time.Time  `json:"end_date"`   // For rentals
	Notes       string      `json:"notes"`
}
