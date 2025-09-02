package middleware

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// Define a custom type for context keys to avoid collisions
type contextKey string

const (
	emailContextKey contextKey = "email"
)

func getSecretKey() string {
	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		return "your_secret_key" // fallback for development
	}
	return secret
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

	secretKey := getSecretKey()
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

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenString string

		// Try to get token from Authorization header first
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			// If no Authorization header, try to get from cookie
			if cookie, err := r.Cookie("auth_token"); err == nil {
				tokenString = cookie.Value
			}
		}

		// If no token found, return unauthorized
		if tokenString == "" {
			http.Error(w, "Missing authentication token", http.StatusUnauthorized)
			return
		}

		// Validate the token using the same logic as in auth handler
		email, err := ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Add the email to the request context
		ctx := context.WithValue(r.Context(), emailContextKey, email)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
