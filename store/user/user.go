package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/PrateekKumar15/CarZone/models"
	"go.opentelemetry.io/otel"
	"golang.org/x/crypto/bcrypt"
)

// Assuming models.UserRequest is defined in your models package
type UserStore struct {
	db *sql.DB
}

func New(db *sql.DB) UserStore {
	return UserStore{db: db}
}

func (s UserStore) CreateUser(ctx context.Context, user models.UserRequest) (err error) {
	tracer := otel.Tracer("AuthStore")
	ctx, span := tracer.Start(ctx, "CreateUser-Store")
	defer span.End()
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Begin the transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Ensure commit or rollback based on err
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// Check if a user with the same email already exists
	var exists bool
	err = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", user.Email).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("user with this email already exists")
	}

	// Insert user into the users table using the transaction
	query := `
		INSERT INTO users (username, email, password_hash, phone, role, profile_data, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	now := time.Now().UTC()

	// Convert profile_data to JSON bytes
	profileDataJSON, err := json.Marshal(map[string]interface{}{})
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, user.UserName, user.Email, string(hashedPassword), user.Phone, user.Role, profileDataJSON, now, now)
	if err != nil {
		return err
	}

	// Zero out the plain password in memory for safety
	user.Password = ""

	return nil
}

func (s UserStore) GetUser(ctx context.Context, email, password string) (models.User, error) {
	tracer := otel.Tracer("AuthStore")
	ctx, span := tracer.Start(ctx, "LoginUser-Store")
	defer span.End()
	var user models.User
	var profileDataJSON []byte
	query := "SELECT id, username, email, password_hash, phone, role, profile_data, created_at, updated_at FROM users WHERE email = $1"
	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.UserName, &user.Email, &user.PasswordHash, &user.Phone, &user.Role, &profileDataJSON, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, err // User not found
		}
		return user, err // Some other error
	}

	// Unmarshal profile_data JSON
	if len(profileDataJSON) > 0 {
		err = json.Unmarshal(profileDataJSON, &user.ProfileData)
		if err != nil {
			return user, err
		}
	} else {
		user.ProfileData = make(map[string]interface{})
	}
	// Compare the provided password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return user, err // Invalid password
	}
	// Zero out the password hash before returning the user object for security
	user.PasswordHash = ""
	return user, nil
}

func (s UserStore) UpdateUser(ctx context.Context, id string, userReq models.UserRequest) (models.User, error) {
	tracer := otel.Tracer("AuthStore")
	ctx, span := tracer.Start(ctx, "UpdateUser-Store")
	defer span.End()

	var updatedUser models.User
	var err error

	// Begin the transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return updatedUser, err
	}

	// Ensure commit or rollback based on err
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// Check if a user with the given id exists and get current data
	var exists bool
	err = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return updatedUser, err
	}
	if !exists {
		return updatedUser, sql.ErrNoRows
	}

	// Hash the new password (after confirming the user exists)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return updatedUser, err
	}

	// Update user in the users table using the transaction
	query := `
		UPDATE users
		SET username = $1, email = $2, password_hash = $3, phone = $4, role = $5, updated_at = $6
		WHERE id = $7
		RETURNING id, username, email, phone, role, profile_data, created_at, updated_at
	`
	now := time.Now().UTC()
	var profileDataJSON []byte
	err = tx.QueryRowContext(ctx, query, userReq.UserName, userReq.Email, string(hashedPassword), userReq.Phone, userReq.Role, now, id).Scan(
		&updatedUser.ID, &updatedUser.UserName, &updatedUser.Email, &updatedUser.Phone, &updatedUser.Role, &profileDataJSON, &updatedUser.CreatedAt, &updatedUser.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return updatedUser, errors.New("no user found with the given ID")
		}
		return updatedUser, err
	}

	// Unmarshal profile_data JSON
	if len(profileDataJSON) > 0 {
		err = json.Unmarshal(profileDataJSON, &updatedUser.ProfileData)
		if err != nil {
			return updatedUser, err
		}
	} else {
		updatedUser.ProfileData = make(map[string]interface{})
	}

	// Zero out the password hash for security
	updatedUser.PasswordHash = ""
	return updatedUser, nil
}

func (s UserStore) DeleteUser(ctx context.Context, id string) (models.User, error) {
	tracer := otel.Tracer("AuthStore")
	ctx, span := tracer.Start(ctx, "DeleteUser-Store")
	defer span.End()

	var deletedUser models.User
	var err error

	// Begin the transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return deletedUser, err
	}

	// Ensure commit or rollback based on err
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// Get user data before deleting (for audit purposes)
	var profileDataJSON []byte
	query := "SELECT id, username, email, phone, role, profile_data, created_at, updated_at FROM users WHERE id = $1"
	err = tx.QueryRowContext(ctx, query, id).Scan(
		&deletedUser.ID, &deletedUser.UserName, &deletedUser.Email, &deletedUser.Phone, &deletedUser.Role, &profileDataJSON, &deletedUser.CreatedAt, &deletedUser.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return deletedUser, errors.New("no user found with the given ID")
		}
		return deletedUser, err
	}

	// Unmarshal profile_data JSON
	if len(profileDataJSON) > 0 {
		err = json.Unmarshal(profileDataJSON, &deletedUser.ProfileData)
		if err != nil {
			return deletedUser, err
		}
	} else {
		deletedUser.ProfileData = make(map[string]interface{})
	}

	// Delete user from the users table using the transaction
	deleteQuery := "DELETE FROM users WHERE id = $1"
	result, err := tx.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		return deletedUser, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return deletedUser, err
	}
	if rowsAffected == 0 {
		return deletedUser, errors.New("no user found with the given ID")
	}

	return deletedUser, nil
}
func (s UserStore) GetAllUsers(ctx context.Context) (users []models.User, err error) {
	tracer := otel.Tracer("AuthStore")
	ctx, span := tracer.Start(ctx, "GetAllUsers-Store")
	defer span.End()
	query := "SELECT id, username, email, phone, role, profile_data, created_at, updated_at FROM users"
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()
	for rows.Next() {
		var user models.User
		var profileDataJSON []byte
		err := rows.Scan(&user.ID, &user.UserName, &user.Email, &user.Phone, &user.Role, &profileDataJSON, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}

		// Unmarshal profile_data JSON
		if len(profileDataJSON) > 0 {
			err = json.Unmarshal(profileDataJSON, &user.ProfileData)
			if err != nil {
				return nil, err
			}
		} else {
			user.ProfileData = make(map[string]interface{})
		}

		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

// GetUserByID retrieves a user by their ID
func (s UserStore) GetUserByID(ctx context.Context, userID string) (models.User, error) {
	tracer := otel.Tracer("AuthStore")
	ctx, span := tracer.Start(ctx, "GetUserByID-Store")
	defer span.End()

	var user models.User
	var profileDataJSON []byte
	query := "SELECT id, username, email, phone, role, profile_data, created_at, updated_at FROM users WHERE id = $1"
	err := s.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID, &user.UserName, &user.Email, &user.Phone, &user.Role, &profileDataJSON, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New("user not found")
		}
		return user, err
	}

	// Unmarshal profile_data JSON
	if len(profileDataJSON) > 0 {
		err = json.Unmarshal(profileDataJSON, &user.ProfileData)
		if err != nil {
			return user, err
		}
	} else {
		user.ProfileData = make(map[string]interface{})
	}

	return user, nil
}

// UpdateProfileData updates only the profile_data field for a user
func (s UserStore) UpdateProfileData(ctx context.Context, userID string, profileData map[string]interface{}) error {
	tracer := otel.Tracer("AuthStore")
	ctx, span := tracer.Start(ctx, "UpdateProfileData-Store")
	defer span.End()

	// Convert profile_data to JSON bytes
	profileDataJSON, err := json.Marshal(profileData)
	if err != nil {
		return err
	}

	query := `
		UPDATE users 
		SET profile_data = $1, updated_at = $2 
		WHERE id = $3
	`
	now := time.Now().UTC()
	result, err := s.db.ExecContext(ctx, query, profileDataJSON, now, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// GetUsersByRole retrieves all users with a specific role
func (s UserStore) GetUsersByRole(ctx context.Context, role string) ([]models.User, error) {
	tracer := otel.Tracer("AuthStore")
	ctx, span := tracer.Start(ctx, "GetUsersByRole-Store")
	defer span.End()

	query := "SELECT id, username, email, phone, role, profile_data, created_at, updated_at FROM users WHERE role = $1"
	rows, err := s.db.QueryContext(ctx, query, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		var profileDataJSON []byte
		err := rows.Scan(&user.ID, &user.UserName, &user.Email, &user.Phone, &user.Role, &profileDataJSON, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}

		// Unmarshal profile_data JSON
		if len(profileDataJSON) > 0 {
			err = json.Unmarshal(profileDataJSON, &user.ProfileData)
			if err != nil {
				return nil, err
			}
		} else {
			user.ProfileData = make(map[string]interface{})
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
