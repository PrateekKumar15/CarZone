package booking

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/PrateekKumar15/CarZone/models"
	"github.com/PrateekKumar15/CarZone/service"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
)

// BookingHandler struct to handle booking-related requests
type BookingHandler struct {
	service service.BookingServiceInterface
}

// NewBookingHandler creates a new BookingHandler with the provided service
func NewBookingHandler(service service.BookingServiceInterface) *BookingHandler {
	return &BookingHandler{service: service}
}

// GetBookingByID retrieves a booking by its ID
func (h *BookingHandler) GetBookingByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer("BookingHandler")
	ctx, span := tracer.Start(ctx, "GetBookingByID-Handler")
	defer span.End()

	vars := mux.Vars(r)
	id := vars["id"]

	resp, err := h.service.GetBookingByID(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error retrieving booking by ID:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if resp == nil {
		http.Error(w, "Booking not found", http.StatusNotFound)
		return
	}

	body, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling response:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error writing response:", err)
		return
	}
}

// GetBookingsByCustomerID retrieves all bookings for a specific customer
func (h *BookingHandler) GetBookingsByCustomerID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer("BookingHandler")
	ctx, span := tracer.Start(ctx, "GetBookingsByCustomerID-Handler")
	defer span.End()

	vars := mux.Vars(r)
	customerID := vars["customerID"]

	resp, err := h.service.GetBookingsByCustomerID(ctx, customerID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error retrieving bookings by customer ID:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	body, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling response:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error writing response:", err)
		return
	}
}

// GetBookingsByCarID retrieves all bookings for a specific car
func (h *BookingHandler) GetBookingsByCarID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer("BookingHandler")
	ctx, span := tracer.Start(ctx, "GetBookingsByCarID-Handler")
	defer span.End()

	vars := mux.Vars(r)
	carID := vars["carID"]

	resp, err := h.service.GetBookingsByCarID(ctx, carID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error retrieving bookings by car ID:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	body, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling response:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error writing response:", err)
		return
	}
}

// GetBookingsByOwnerID retrieves all bookings for cars owned by a specific owner
func (h *BookingHandler) GetBookingsByOwnerID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer("BookingHandler")
	ctx, span := tracer.Start(ctx, "GetBookingsByOwnerID-Handler")
	defer span.End()

	vars := mux.Vars(r)
	ownerID := vars["ownerID"]

	resp, err := h.service.GetBookingsByOwnerID(ctx, ownerID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error retrieving bookings by owner ID:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	body, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling response:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error writing response:", err)
		return
	}
}

// CreateBooking creates a new booking
func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer("BookingHandler")
	ctx, span := tracer.Start(ctx, "CreateBooking-Handler")
	defer span.End()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error reading request body:", err)
		return
	}

	var bookingReq models.BookingRequest
	err = json.Unmarshal(body, &bookingReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error unmarshalling request body:", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	resp, err := h.service.CreateBooking(ctx, bookingReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error creating booking:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	responseBody, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling response:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write(responseBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error writing response:", err)
		return
	}
}

// UpdateBookingStatus updates the status of an existing booking
func (h *BookingHandler) UpdateBookingStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer("BookingHandler")
	ctx, span := tracer.Start(ctx, "UpdateBookingStatus-Handler")
	defer span.End()

	vars := mux.Vars(r)
	id := vars["id"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error reading request body:", err)
		return
	}

	var statusUpdate struct {
		Status models.BookingStatus `json:"status"`
	}
	err = json.Unmarshal(body, &statusUpdate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error unmarshalling request body:", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	resp, err := h.service.UpdateBookingStatus(ctx, id, statusUpdate.Status)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error updating booking status:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	responseBody, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling response:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(responseBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error writing response:", err)
		return
	}
}

// DeleteBooking deletes a booking
func (h *BookingHandler) DeleteBooking(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer("BookingHandler")
	ctx, span := tracer.Start(ctx, "DeleteBooking-Handler")
	defer span.End()

	vars := mux.Vars(r)
	id := vars["id"]

	resp, err := h.service.DeleteBooking(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error deleting booking:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	body, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling response:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error writing response:", err)
		return
	}
}

// GetAllBookings retrieves all bookings (admin function)
func (h *BookingHandler) GetAllBookings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer("BookingHandler")
	ctx, span := tracer.Start(ctx, "GetAllBookings-Handler")
	defer span.End()

	resp, err := h.service.GetAllBookings(ctx)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error retrieving all bookings:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	body, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling response:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error writing response:", err)
		return
	}
}
