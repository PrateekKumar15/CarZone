package models

import (
	"errors"
	"regexp"
	"time"

	"github.com/google/uuid"
)

// User represents a user account in the system.
// Fields follow the style used in existing models (UUID, JSON tags, timestamps).
type User struct {
	ID           uuid.UUID              `json:"id"`
	Email        string                 `json:"email"`
	PasswordHash string                 `json:"password_hash"`
	UserName     string                 `json:"username"`
	Phone        string                 `json:"phone"`
	Role         string                 `json:"role"`
	ProfileData  map[string]interface{} `json:"profile_data"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

// UserRequest represents the payload used to create or update a user.
// It intentionally excludes fields like ID and timestamps which are managed by the system.
type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	UserName string `json:"username"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ValidateUserRequest validates a UserRequest. Returns nil when valid, otherwise an error.
func ValidateUserRequest(req UserRequest) error {
	if err := validateEmail(req.Email); err != nil {
		return err
	}
	if err := validatePassword(req.Password); err != nil {
		return err
	}
	if len(req.UserName) == 0 {
		return errors.New("username cannot be empty")
	}
	if err := validatePhone(req.Phone); err != nil {
		return err
	}
	if err := validateRole(req.Role); err != nil {
		return err
	}

	return nil
}

// validateEmail uses a simple regex to check email format.
func validateEmail(email string) error {
	if email == "" {
		return errors.New("email cannot be empty")
	}
	// very small, permissive email regex
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

// validatePassword enforces a minimal password length requirement.
func validatePassword(pw string) error {
	if len(pw) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	return nil
}

// validatePhone checks if phone number has a valid format.
func validatePhone(phone string) error {
	if phone == "" {
		return errors.New("phone number cannot be empty")
	}
	// Phone number should be between 10-15 digits and can include + and - characters
	re := regexp.MustCompile(`^\+?[1-9]\d{1,14}$|^[0-9]{10,15}$`)
	if !re.MatchString(phone) {
		return errors.New("invalid phone number format (should be 10-15 digits, optionally starting with +)")
	}
	return nil
}

// validateRole ensures role is one of the allowed values.
func validateRole(role string) error {
	if role == "" {
		return errors.New("role cannot be empty")
	}
	allowedRoles := []string{"owner", "renter", "admin"}
	for _, allowedRole := range allowedRoles {
		if role == allowedRole {
			return nil
		}
	}
	return errors.New("role must be one of: owner, renter, admin")
}

// NewUserFromRequest creates a new User from a validated UserRequest.
// Note: this does NOT hash the password; hashing should be performed by the caller
// before assigning to PasswordHash (to avoid importing crypto libraries in models).
func NewUserFromRequest(req UserRequest) *User {
	return &User{
		ID:          uuid.New(),
		Email:       req.Email,
		UserName:    req.UserName,
		Phone:       req.Phone,
		Role:        req.Role,
		ProfileData: make(map[string]interface{}), // Initialize empty profile data
		// PasswordHash should be set by caller after hashing the provided password.
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func ValidateLoginRequest(req LoginRequest) error {
	if err := validateEmail(req.Email); err != nil {
		return err
	}
	if err := validatePassword(req.Password); err != nil {
		return err
	}
	return nil
}
