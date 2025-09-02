package auth

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PrateekKumar15/CarZone/models"
	"github.com/PrateekKumar15/CarZone/service"
	jwt "github.com/dgrijalva/jwt-go"
	"go.opentelemetry.io/otel"
)

type AuthHandler struct {
	service service.AuthServiceInterface
}

// NewCarHandler creates a new CarHandler with the provided service
func NewAuthHandler(service service.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer("AuthHandler")
	ctx, span := tracer.Start(ctx, "LoginUser-Handler")
	defer span.End()

	var credentials models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Use the login service to authenticate user
	user, err := h.service.LoginUser(ctx, credentials)
	if err != nil {
		log.Println("Error logging in user:", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	tokenString, err := GenerateTokenAndSetCookie(w, credentials.Email)
	if err != nil {
		log.Println("Error generating token:", err)
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"user":    user,
		"token":   tokenString,
		"message": "Login successful",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GenerateTokenAndSetCookie(w http.ResponseWriter, email string) (string, error) {
	// Create the JWT claims, which includes the username and expiry time
	secretKey := os.Getenv("SECRET_KEY")
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    "CarZone",
		Subject:   email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	// Set token as HTTP-only cookie for automatic inclusion in subsequent requests
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    signedToken,
		Path:     "/",
		HttpOnly: true,  // Prevents JavaScript access (XSS protection)
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   24 * 60 * 60, // 24 hours in seconds
	})

	// Set token in header for easy access
	w.Header().Set("Authorization", "Bearer "+signedToken)
	w.Header().Set("Access-Control-Expose-Headers", "Authorization") // Allow frontend to access this header

	return signedToken, nil
}



// ValidateToken validates a JWT token and returns the email (stored in Subject) if valid
func ValidateToken(tokenString string) (string, error) {
	if tokenString == "" {
		return "", errors.New("empty token")
	}

	// Accept tokens prefixed with "Bearer "
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	secretKey := os.Getenv("SECRET_KEY")
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}

	// Check expiry if present
	if claims.ExpiresAt != 0 && time.Now().Unix() > claims.ExpiresAt {
		return "", errors.New("token expired")
	}

	if claims.Subject == "" {
		return "", errors.New("email not found in token")
	}

	// Subject contains the email
	return claims.Subject, nil
}

func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer("AuthHandler")
	ctx, span := tracer.Start(ctx, "RegisterUser-Handler")
	defer span.End()

	var userReq models.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Use the registration service to create a new user
	if err := h.service.RegisterUser(ctx, userReq); err != nil {
		log.Println("Error registering user:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// After successful registration, authenticate the user and generate token
	loginReq := models.LoginRequest{
		Email:    userReq.Email,
		Password: userReq.Password,
	}

	// Get the newly created user from the database
	user, err := h.service.LoginUser(ctx, loginReq)
	if err != nil {
		log.Println("Error retrieving newly registered user:", err)
		http.Error(w, "Registration successful but failed to authenticate", http.StatusInternalServerError)
		return
	}

	// Generate token and set cookie/headers
	tokenString, err := GenerateTokenAndSetCookie(w, userReq.Email)
	if err != nil {
		log.Println("Error generating token for new user:", err)
		http.Error(w, "Registration successful but failed to generate token", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"user":    user,
		"token":   tokenString,
		"message": "User registered and logged in successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear the auth_token cookie by setting its MaxAge to -1
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",		
		MaxAge:  -1 ,
	})
	
	response := map[string]interface{}{
		"message": "Logout successful",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
