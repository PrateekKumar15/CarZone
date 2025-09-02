package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/PrateekKumar15/CarZone/models"
	"go.opentelemetry.io/otel"
	"golang.org/x/crypto/bcrypt"
)

// Assuming models.UserRequest is defined in your models package
type AuthStore struct {
	db *sql.DB
}

func New(db *sql.DB) AuthStore {
	return AuthStore{db: db}
}

func (s AuthStore) CreateUser(ctx context.Context, user models.UserRequest) (err error) {
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
		INSERT INTO users (name, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	now := time.Now().UTC()
	_, err = tx.ExecContext(ctx, query, user.UserName, user.Email, string(hashedPassword), now, now)
	if err != nil {
		return err
	}

	// Zero out the plain password in memory for safety
	user.Password = ""

	return nil
}

func (s AuthStore) GetUser(ctx context.Context, email, password string) (models.User, error) {
	tracer := otel.Tracer("AuthStore")
	ctx, span := tracer.Start(ctx, "LoginUser-Store")
	defer span.End()
	var user models.User
	query := "SELECT id, name, email, password_hash, created_at, updated_at FROM users WHERE email = $1"
	err := s.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.UserName, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, err // User not found
		}
		return user, err // Some other error
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

func (s AuthStore) UpdateUser(ctx context.Context, id string, userReq models.UserRequest) (models.User, error) {
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
		SET name = $1, email = $2, password_hash = $3, updated_at = $4
		WHERE id = $5
		RETURNING id, name, email, created_at, updated_at
	`
	now := time.Now().UTC()
	err = tx.QueryRowContext(ctx, query, userReq.UserName, userReq.Email, string(hashedPassword), now, id).Scan(
		&updatedUser.ID, &updatedUser.UserName, &updatedUser.Email, &updatedUser.CreatedAt, &updatedUser.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return updatedUser, errors.New("no user found with the given ID")
		}
		return updatedUser, err
	}

	// Zero out the password hash for security
	updatedUser.PasswordHash = ""
	return updatedUser, nil
}

func (s AuthStore) DeleteUser(ctx context.Context, id string) (models.User, error) {
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
	query := "SELECT id, name, email, created_at, updated_at FROM users WHERE id = $1"
	err = tx.QueryRowContext(ctx, query, id).Scan(&deletedUser.ID, &deletedUser.UserName, &deletedUser.Email, &deletedUser.CreatedAt, &deletedUser.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return deletedUser, errors.New("no user found with the given ID")
		}
		return deletedUser, err
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
func (s AuthStore) GetAllUsers(ctx context.Context) (users []models.User, err error) {
	tracer := otel.Tracer("AuthStore")
	ctx, span := tracer.Start(ctx, "GetAllUsers-Store")
	defer span.End()
	query := "SELECT id, name, email, created_at, updated_at FROM users"
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
		err := rows.Scan(&user.ID, &user.UserName, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
