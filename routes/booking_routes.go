package routes

import (
	"github.com/gorilla/mux"
)

// setupBookingRoutes configures all booking-related routes
func (r *Router) setupBookingRoutes(router *mux.Router) {
	// Booking CRUD operations

	// GET /bookings - Retrieve all bookings for authenticated user
	// Returns bookings based on user's role and permissions
	router.HandleFunc("/bookings", r.BookingHandler.GetAllBookings).Methods("GET")

	// GET /bookings/{id} - Retrieve a specific booking by its UUID
	// Path parameter: UUID of the booking
	router.HandleFunc("/bookings/{id}", r.BookingHandler.GetBookingByID).Methods("GET")

	// POST /bookings - Create a new booking
	// Body: Booking JSON data with customer_id, car_id, booking details
	router.HandleFunc("/bookings", r.BookingHandler.CreateBooking).Methods("POST")

	// DELETE /bookings/{id} - Delete a booking by its UUID
	// Path parameter: UUID of the booking to delete
	router.HandleFunc("/bookings/{id}", r.BookingHandler.DeleteBooking).Methods("DELETE")

	// Booking status management

	// PUT /bookings/{id}/status - Update booking status
	// Path parameter: UUID of the booking
	// Body: { "status": "confirmed|cancelled|completed" }
	router.HandleFunc("/bookings/{id}/status", r.BookingHandler.UpdateBookingStatus).Methods("PUT")

	// Booking query endpoints

	// GET /bookings/customer/{customerID} - Get all bookings for a specific customer
	// Path parameter: UUID of the customer
	router.HandleFunc("/bookings/customer/{customerID}", r.BookingHandler.GetBookingsByCustomerID).Methods("GET")

	// GET /bookings/car/{carID} - Get all bookings for a specific car
	// Path parameter: UUID of the car
	router.HandleFunc("/bookings/car/{carID}", r.BookingHandler.GetBookingsByCarID).Methods("GET")

	// GET /bookings/owner/{ownerID} - Get all bookings for cars owned by a specific owner
	// Path parameter: UUID of the car owner
	router.HandleFunc("/bookings/owner/{ownerID}", r.BookingHandler.GetBookingsByOwnerID).Methods("GET")
}
