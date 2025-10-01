package payment

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/PrateekKumar15/CarZone/models"
	"github.com/PrateekKumar15/CarZone/service"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
)

// PaymentHandler handles HTTP requests for payment operations
type PaymentHandler struct {
	paymentService service.PaymentServiceInterface
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(paymentService service.PaymentServiceInterface) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// CreatePayment handles payment creation requests
func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("PaymentHandler")
	ctx, span := tracer.Start(r.Context(), "CreatePayment-Handler")
	defer span.End()

	// Handle OPTIONS request for CORS preflight
	if r.Method == "OPTIONS" {
		return // CORS middleware will handle the response
	}

	var paymentReq models.PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&paymentReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	razorpayOrder, err := h.paymentService.CreatePayment(ctx, &paymentReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(razorpayOrder)
}

// VerifyPayment handles payment verification requests
func (h *PaymentHandler) VerifyPayment(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("PaymentHandler")
	ctx, span := tracer.Start(r.Context(), "VerifyPayment-Handler")
	defer span.End()

	// Handle OPTIONS request for CORS preflight
	if r.Method == "OPTIONS" {
		return // CORS middleware will handle the response
	}

	var verificationReq models.PaymentVerificationRequest
	if err := json.NewDecoder(r.Body).Decode(&verificationReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	payment, err := h.paymentService.VerifyPayment(ctx, &verificationReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Payment verified successfully",
		"payment": payment,
	})
}

// GetPaymentByID handles requests to get a payment by ID
func (h *PaymentHandler) GetPaymentByID(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("PaymentHandler")
	ctx, span := tracer.Start(r.Context(), "GetPaymentByID-Handler")
	defer span.End()

	vars := mux.Vars(r)
	paymentID := vars["id"]

	if paymentID == "" {
		http.Error(w, "Payment ID is required", http.StatusBadRequest)
		return
	}

	payment, err := h.paymentService.GetPaymentByID(ctx, paymentID)
	if err != nil {
		if err.Error() == "payment not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payment)
}

// GetPaymentByBookingID handles requests to get payment by booking ID
func (h *PaymentHandler) GetPaymentByBookingID(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("PaymentHandler")
	ctx, span := tracer.Start(r.Context(), "GetPaymentByBookingID-Handler")
	defer span.End()

	vars := mux.Vars(r)
	bookingID := vars["booking_id"]

	if bookingID == "" {
		http.Error(w, "Booking ID is required", http.StatusBadRequest)
		return
	}

	payment, err := h.paymentService.GetPaymentByBookingID(ctx, bookingID)
	if err != nil {
		if err.Error() == "payment not found for booking" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payment)
}

// GetPaymentsByUserID handles requests to get payments by user ID
func (h *PaymentHandler) GetPaymentsByUserID(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("PaymentHandler")
	ctx, span := tracer.Start(r.Context(), "GetPaymentsByUserID-Handler")
	defer span.End()

	vars := mux.Vars(r)
	userID := vars["user_id"]

	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	payments, err := h.paymentService.GetPaymentsByUserID(ctx, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payments)
}

// ProcessRefund handles refund requests
func (h *PaymentHandler) ProcessRefund(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("PaymentHandler")
	ctx, span := tracer.Start(r.Context(), "ProcessRefund-Handler")
	defer span.End()

	vars := mux.Vars(r)
	paymentID := vars["payment_id"]

	if paymentID == "" {
		http.Error(w, "Payment ID is required", http.StatusBadRequest)
		return
	}

	// Parse refund amount from request body
	var refundReq struct {
		Amount float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&refundReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	payment, err := h.paymentService.ProcessRefund(ctx, paymentID, refundReq.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Refund processed successfully",
		"payment": payment,
	})
}

// GetAllPayments handles requests to get all payments
func (h *PaymentHandler) GetAllPayments(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("PaymentHandler")
	ctx, span := tracer.Start(r.Context(), "GetAllPayments-Handler")
	defer span.End()

	// Parse query parameters for pagination (optional)
	limit := 50 // default limit
	offset := 0 // default offset

	if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetParam := r.URL.Query().Get("offset"); offsetParam != "" {
		if o, err := strconv.Atoi(offsetParam); err == nil && o >= 0 {
			offset = o
		}
	}

	payments, err := h.paymentService.GetAllPayments(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Apply pagination (simple slice-based pagination for now)
	totalPayments := len(*payments)
	start := offset
	end := offset + limit

	if start > totalPayments {
		start = totalPayments
	}
	if end > totalPayments {
		end = totalPayments
	}

	paginatedPayments := (*payments)[start:end]

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"payments": paginatedPayments,
		"total":    totalPayments,
		"limit":    limit,
		"offset":   offset,
		"has_more": end < totalPayments,
	})
}
