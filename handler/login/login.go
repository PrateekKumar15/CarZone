package login

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/PrateekKumar15/CarZone/models"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	valid := credentials.UserName == "admin" && credentials.Password == "admin123" // Replace with real validation
	if !valid {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	tokenString, err := GenerateToken(credentials.UserName)
	if err != nil {
		log.Println("Error generating token:", err)
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"token": tokenString, "message": "Login successful"}

	// Set token as HTTP-only cookie for automatic inclusion in subsequent requests
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,  // Prevents JavaScript access (XSS protection)
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   24 * 60 * 60, // 24 hours in seconds
	})

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", "Bearer "+tokenString)           // Set token in header for easy access
	w.Header().Set("Access-Control-Expose-Headers", "Authorization") // Allow frontend to access this header
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    "CarZone",
		Subject:   username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte("your_secret_key")) // Replace with your actual secret key
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ValidateToken validates a JWT token and returns the username if valid
func ValidateToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte("your_secret_key"), nil // Same secret key used for signing
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		// Check if token is expired
		if time.Now().Unix() > claims.ExpiresAt {
			return "", errors.New("token expired")
		}
		return claims.Subject, nil // Return username
	}

	return "", errors.New("invalid token")
}
