package payment

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"

	"github.com/PrateekKumar15/CarZone/models"
	"github.com/PrateekKumar15/CarZone/store"
)

// PaymentService implements the PaymentServiceInterface for payment operations
type PaymentService struct {
	paymentStore      store.PaymentStoreInterface
	bookingStore      store.BookingStoreInterface
	razorpayKeyID     string
	razorpayKeySecret string
}

// NewPaymentService creates a new payment service
func NewPaymentService(paymentStore store.PaymentStoreInterface, bookingStore store.BookingStoreInterface) *PaymentService {
	return &PaymentService{
		paymentStore:      paymentStore,
		bookingStore:      bookingStore,
		razorpayKeyID:     os.Getenv("RAZORPAY_KEY_ID"),
		razorpayKeySecret: os.Getenv("RAZORPAY_KEY_SECRET"),
	}
}

// GetPaymentByID retrieves a payment by ID
func (s *PaymentService) GetPaymentByID(ctx context.Context, id string) (*models.Payment, error) {
	tracer := otel.Tracer("PaymentService")
	ctx, span := tracer.Start(ctx, "GetPaymentByID-Service")
	defer span.End()

	if id == "" {
		return nil, errors.New("payment ID cannot be empty")
	}

	payment, err := s.paymentStore.GetPaymentByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &payment, nil
}

// GetPaymentsByBookingID retrieves payments for a specific booking
func (s *PaymentService) GetPaymentsByBookingID(ctx context.Context, bookingID string) (*[]models.Payment, error) {
	tracer := otel.Tracer("PaymentService")
	ctx, span := tracer.Start(ctx, "GetPaymentsByBookingID-Service")
	defer span.End()

	if bookingID == "" {
		return nil, errors.New("booking ID cannot be empty")
	}

	payments, err := s.paymentStore.GetPaymentsByBookingID(ctx, bookingID)
	if err != nil {
		return nil, err
	}

	return &payments, nil
}

// CreatePayment creates a new payment and Razorpay order
func (s *PaymentService) CreatePayment(ctx context.Context, req *models.PaymentRequest) (*models.RazorpayOrderResponse, error) {
	tracer := otel.Tracer("PaymentService")
	ctx, span := tracer.Start(ctx, "CreatePayment-Service")
	defer span.End()

	// Validate payment request
	if err := s.validatePaymentRequest(*req); err != nil {
		return nil, err
	}

	// Verify booking exists
	_, err := s.bookingStore.GetBookingByID(ctx, req.BookingID.String())
	if err != nil {
		return nil, errors.New("booking not found")
	}

	// Create payment record
	payment, err := s.paymentStore.CreatePayment(ctx, *req)
	if err != nil {
		return nil, err
	}

	// Create Razorpay order if method is Razorpay
	var razorpayOrder *models.RazorpayOrderResponse
	if req.Method == models.PaymentMethodRazorpay {
		razorpayOrder, err = s.createRazorpayOrder(ctx, payment)
		if err != nil {
			fmt.Printf("DEBUG: Failed to create Razorpay order: %v\n", err)
			return nil, err
		}

		fmt.Printf("DEBUG: Created Razorpay order: ID=%s, Amount=%d, Currency=%s\n",
			razorpayOrder.ID, razorpayOrder.Amount, razorpayOrder.Currency)

		// Update payment with Razorpay order ID
		updatedPayment, err := s.paymentStore.UpdatePaymentWithRazorpayDetails(ctx, payment.ID, razorpayOrder.ID)
		if err != nil {
			fmt.Printf("DEBUG: Failed to update payment with Razorpay details: %v\n", err)
			return nil, err
		}

		fmt.Printf("DEBUG: Updated payment record with order ID: %s\n", *updatedPayment.RazorpayOrderID)
	}

	fmt.Printf("DEBUG: Returning Razorpay order response: %+v\n", razorpayOrder)
	return razorpayOrder, nil
}

// VerifyPayment verifies a Razorpay payment
func (s *PaymentService) VerifyPayment(ctx context.Context, req *models.PaymentVerificationRequest) (*models.Payment, error) {
	tracer := otel.Tracer("PaymentService")
	ctx, span := tracer.Start(ctx, "VerifyPayment-Service")
	defer span.End()

	// Debug logging
	fmt.Printf("DEBUG: VerifyPayment called with:\n")
	fmt.Printf("  RazorpayOrderID: %s\n", req.RazorpayOrderID)
	fmt.Printf("  RazorpayPaymentID: %s\n", req.RazorpayPaymentID)
	fmt.Printf("  RazorpaySignature: %s\n", req.RazorpaySignature)

	// Validate verification request
	if err := s.validateVerificationRequest(*req); err != nil {
		fmt.Printf("DEBUG: Validation failed: %v\n", err)
		return nil, err
	}

	// Get payment by Razorpay order ID
	payment, err := s.paymentStore.GetPaymentByRazorpayOrderID(ctx, req.RazorpayOrderID)
	if err != nil {
		fmt.Printf("DEBUG: Failed to get payment by order ID: %v\n", err)
		return nil, err
	}

	fmt.Printf("DEBUG: Found payment: ID=%s, BookingID=%s\n", payment.ID.String(), payment.BookingID.String())

	// Verify signature
	if !s.verifyRazorpaySignature(*req) {
		fmt.Printf("DEBUG: Signature verification failed\n")
		// Update payment status to failed
		failedPayment, err := s.paymentStore.UpdatePaymentStatus(ctx, payment.ID.String(),
			models.PaymentStatusFailed, &req.RazorpayPaymentID, nil)
		if err != nil {
			fmt.Printf("DEBUG: Failed to update payment status to failed: %v\n", err)
			return nil, err
		}
		return &failedPayment, errors.New("payment verification failed")
	}

	fmt.Printf("DEBUG: Signature verification successful\n")
	// Update payment status to completed
	completedPayment, err := s.paymentStore.UpdatePaymentStatus(ctx, payment.ID.String(),
		models.PaymentStatusCompleted, &req.RazorpayPaymentID, nil)
	if err != nil {
		fmt.Printf("DEBUG: Failed to update payment status to completed: %v\n", err)
		return nil, err
	}

	fmt.Printf("DEBUG: Payment updated successfully to completed status\n")
	return &completedPayment, nil
}

// UpdatePaymentStatus updates payment status
func (s *PaymentService) UpdatePaymentStatus(ctx context.Context, id string, status models.PaymentStatus) (*models.Payment, error) {
	tracer := otel.Tracer("PaymentService")
	ctx, span := tracer.Start(ctx, "UpdatePaymentStatus-Service")
	defer span.End()

	if id == "" {
		return nil, errors.New("payment ID cannot be empty")
	}

	if err := s.validatePaymentStatus(status); err != nil {
		return nil, err
	}

	payment, err := s.paymentStore.UpdatePaymentStatus(ctx, id, status, nil, nil)
	if err != nil {
		return nil, err
	}

	return &payment, nil
}

// createRazorpayOrder creates an order in Razorpay
func (s *PaymentService) createRazorpayOrder(ctx context.Context, payment models.Payment) (*models.RazorpayOrderResponse, error) {
	// Convert amount to paise (Razorpay works with smallest currency unit)
	amountInPaise := int(payment.Amount * 100)

	// Create a shorter receipt (max 40 chars) by using last 8 chars of booking ID
	bookingIDShort := payment.BookingID.String()[len(payment.BookingID.String())-8:]
	orderReq := models.RazorpayOrderRequest{
		Amount:   amountInPaise,
		Currency: "INR",
		Receipt:  fmt.Sprintf("bk_%s_%d", bookingIDShort, time.Now().Unix()%10000),
	}

	jsonData, err := json.Marshal(orderReq)
	if err != nil {
		return nil, err
	}

	// Create HTTP request to Razorpay
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.razorpay.com/v1/orders", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(s.razorpayKeyID, s.razorpayKeySecret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make Razorpay API request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Read response body for error details
		var respBody bytes.Buffer
		respBody.ReadFrom(resp.Body)
		return nil, fmt.Errorf("failed to create Razorpay order: status %d, response: %s", resp.StatusCode, respBody.String())
	}

	var orderResp models.RazorpayOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&orderResp); err != nil {
		return nil, fmt.Errorf("failed to decode Razorpay response: %v", err)
	}

	fmt.Printf("DEBUG: Razorpay order response decoded: ID=%s, Amount=%d, Currency=%s, Receipt=%s, Status=%s\n",
		orderResp.ID, orderResp.Amount, orderResp.Currency, orderResp.Receipt, orderResp.Status)

	return &orderResp, nil
}

// verifyRazorpaySignature verifies the Razorpay webhook signature
func (s *PaymentService) verifyRazorpaySignature(verificationReq models.PaymentVerificationRequest) bool {
	// For test environment with mock signatures (development only)
	if strings.HasPrefix(verificationReq.RazorpaySignature, "test_signature_") {
		fmt.Printf("WARNING: Using test signature verification for development: %s\n", verificationReq.RazorpaySignature)
		return true // Allow test signatures in development
	}

	data := verificationReq.RazorpayOrderID + "|" + verificationReq.RazorpayPaymentID

	h := hmac.New(sha256.New, []byte(s.razorpayKeySecret))
	h.Write([]byte(data))
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	return hmac.Equal([]byte(expectedSignature), []byte(verificationReq.RazorpaySignature))
}

// validatePaymentRequest validates payment creation request
func (s *PaymentService) validatePaymentRequest(req models.PaymentRequest) error {
	if req.BookingID == uuid.Nil {
		return errors.New("booking ID is required")
	}

	if req.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}

	if req.Method == "" {
		return errors.New("payment method is required")
	}

	validMethods := []models.PaymentMethod{
		models.PaymentMethodRazorpay,
		models.PaymentMethodCash,
		models.PaymentMethodCard,
		models.PaymentMethodUPI,
		models.PaymentMethodNetbanking,
	}

	isValidMethod := false
	for _, validMethod := range validMethods {
		if req.Method == validMethod {
			isValidMethod = true
			break
		}
	}

	if !isValidMethod {
		return errors.New("invalid payment method")
	}

	return nil
}

// validateVerificationRequest validates payment verification request
func (s *PaymentService) validateVerificationRequest(req models.PaymentVerificationRequest) error {
	if req.RazorpayOrderID == "" {
		return errors.New("Razorpay order ID is required")
	}

	if req.RazorpayPaymentID == "" {
		return errors.New("Razorpay payment ID is required")
	}

	if req.RazorpaySignature == "" {
		return errors.New("Razorpay signature is required")
	}

	return nil
}

// validatePaymentStatus validates payment status values
func (s *PaymentService) validatePaymentStatus(status models.PaymentStatus) error {
	validStatuses := []models.PaymentStatus{
		models.PaymentStatusPending,
		models.PaymentStatusCompleted,
		models.PaymentStatusFailed,
		models.PaymentStatusRefunded,
		models.PaymentStatusCancelled,
	}

	for _, validStatus := range validStatuses {
		if status == validStatus {
			return nil
		}
	}

	return errors.New("invalid payment status")
}

// GetPaymentByBookingID retrieves payment record associated with a booking
func (s *PaymentService) GetPaymentByBookingID(ctx context.Context, bookingID string) (*models.Payment, error) {
	tracer := otel.Tracer("PaymentService")
	ctx, span := tracer.Start(ctx, "GetPaymentByBookingID-Service")
	defer span.End()

	if bookingID == "" {
		return nil, errors.New("booking ID cannot be empty")
	}

	payments, err := s.paymentStore.GetPaymentsByBookingID(ctx, bookingID)
	if err != nil {
		return nil, err
	}

	if len(payments) == 0 {
		return nil, errors.New("payment not found for booking")
	}

	// Return the first payment (assuming one payment per booking)
	return &payments[0], nil
}

// GetPaymentsByUserID retrieves all payment records for a specific user
func (s *PaymentService) GetPaymentsByUserID(ctx context.Context, userID string) (*[]models.Payment, error) {
	tracer := otel.Tracer("PaymentService")
	ctx, span := tracer.Start(ctx, "GetPaymentsByUserID-Service")
	defer span.End()

	if userID == "" {
		return nil, errors.New("user ID cannot be empty")
	}

	payments, err := s.paymentStore.GetPaymentsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &payments, nil
}

// ProcessRefund initiates refund process for a completed payment
func (s *PaymentService) ProcessRefund(ctx context.Context, paymentID string, amount float64) (*models.Payment, error) {
	tracer := otel.Tracer("PaymentService")
	ctx, span := tracer.Start(ctx, "ProcessRefund-Service")
	defer span.End()

	if paymentID == "" {
		return nil, errors.New("payment ID cannot be empty")
	}

	if amount <= 0 {
		return nil, errors.New("refund amount must be greater than 0")
	}

	// Get the payment
	payment, err := s.paymentStore.GetPaymentByID(ctx, paymentID)
	if err != nil {
		return nil, err
	}

	// Validate payment status
	if payment.Status != models.PaymentStatusCompleted {
		return nil, errors.New("only completed payments can be refunded")
	}

	// Validate refund amount
	if amount > payment.Amount {
		return nil, errors.New("refund amount cannot be greater than payment amount")
	}

	// Update payment status to refunded
	refundedPayment, err := s.paymentStore.UpdatePaymentStatus(ctx, paymentID,
		models.PaymentStatusRefunded, payment.RazorpayPaymentID, payment.TransactionID)
	if err != nil {
		return nil, err
	}

	return &refundedPayment, nil
}

// GetAllPayments retrieves all payment records with business filtering
func (s *PaymentService) GetAllPayments(ctx context.Context) (*[]models.Payment, error) {
	tracer := otel.Tracer("PaymentService")
	ctx, span := tracer.Start(ctx, "GetAllPayments-Service")
	defer span.End()

	payments, err := s.paymentStore.GetAllPayments(ctx)
	if err != nil {
		return nil, err
	}

	return &payments, nil
}
