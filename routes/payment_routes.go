package routes

import (
	"github.com/gorilla/mux"
)

// setupPaymentRoutes configures all payment-related routes
func (r *Router) setupPaymentRoutes(router *mux.Router) {
	// Payment operations - All routes require authentication

	// Create payment and get Razorpay order
	router.HandleFunc("/payments", r.PaymentHandler.CreatePayment).Methods("POST", "OPTIONS")

	// Get all payments (admin only - consider adding admin middleware later)
	router.HandleFunc("/payments", r.PaymentHandler.GetAllPayments).Methods("GET", "OPTIONS")

	// Verify payment after successful transaction
	router.HandleFunc("/payments/verify", r.PaymentHandler.VerifyPayment).Methods("POST", "OPTIONS")

	// Get payment by ID
	router.HandleFunc("/payments/{id}", r.PaymentHandler.GetPaymentByID).Methods("GET", "OPTIONS")

	// Get payment by booking ID
	router.HandleFunc("/payments/booking/{booking_id}", r.PaymentHandler.GetPaymentByBookingID).Methods("GET", "OPTIONS")

	// Get all payments for a user
	router.HandleFunc("/payments/user/{user_id}", r.PaymentHandler.GetPaymentsByUserID).Methods("GET", "OPTIONS")

	// Process refund for a payment
	router.HandleFunc("/payments/{payment_id}/refund", r.PaymentHandler.ProcessRefund).Methods("POST", "OPTIONS")
}
