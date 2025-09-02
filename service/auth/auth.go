package auth

import (
	"errors"
	"fmt"
	"net/mail"

	"context"

	"github.com/PrateekKumar15/CarZone/models"
	"github.com/PrateekKumar15/CarZone/store"
)

// Assuming models.UserRequest is defined in your models package
type AuthService struct {
	store store.UserStoreInterface
}

func NewAuthService(store store.UserStoreInterface) *AuthService {
	return &AuthService{store: store}
}

func (s *AuthService) RegisterUser(ctx context.Context, userReq models.UserRequest) error {
	// Validate the user request
	if err := models.ValidateUserRequest(userReq); err != nil {
		return err
	}
	// Validate email format
	if _, err := mail.ParseAddress(userReq.Email); err != nil {
		return errors.New("invalid email format")
	}
	// Create the user in the store
	if err := s.store.CreateUser(ctx,userReq); err != nil {
		return err
	}
	fmt.Printf("User %s registered successfully\n", userReq.Email)

	return nil
}

func (s *AuthService) LoginUser(ctx context.Context, loginReq models.LoginRequest) ( models.User,error) {
	var user models.User
	// Validate the login request
	if err := models.ValidateLoginRequest(loginReq); err != nil {
		return user,  err
	}
	// Authenticate the user in the store
	user, err := s.store.GetUser(ctx, loginReq.Email, loginReq.Password)
	if err != nil {
		return user, err
	}
	return user, nil
}
// UserStoreInterface defines the contract for user data persistence operations.
// This interface abstracts the underlying data store (e.g., SQL, NoSQL) and provides

