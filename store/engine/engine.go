// Package engine implements the data access layer for engine entities in the CarZone application.
// This package provides concrete implementations of the EngineStoreInterface,
// handling all database operations related to engine management including CRUD operations,
// transaction management, and data validation.
package engine

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/PrateekKumar15/CarZone/models"
	"github.com/google/uuid"
)

// EngineStore is a concrete implementation of the EngineStoreInterface.
// It encapsulates database operations for engine entities and manages
// PostgreSQL database interactions with proper transaction handling.
type EngineStore struct {
	db *sql.DB // Database connection pool for PostgreSQL operations
}

// New creates and returns a new instance of EngineStore.
// It initializes the store with a database connection pool for performing
// engine-related database operations.
//
// Parameters:
//   - db: Database connection pool (*sql.DB) for PostgreSQL operations
//
// Returns:
//   - *EngineStore: A new instance of EngineStore ready for database operations
func New(db *sql.DB) *EngineStore {
	return &EngineStore{db: db}
}

// GetEngineByID retrieves a single engine record from the database by its unique identifier.
// This method implements the EngineStoreInterface contract and handles all aspects of
// engine retrieval including UUID validation, transaction management, and error handling.
//
// The method performs the following operations:
// 1. Validates and parses the engine ID from string to UUID format
// 2. Begins a database transaction for consistency
// 3. Executes a parameterized query to prevent SQL injection
// 4. Handles transaction rollback/commit based on operation success
// 5. Returns the engine record or appropriate error
//
// Parameters:
//   - ctx: Request context for cancellation, timeout, and request scoping
//   - id: Engine unique identifier in UUID string format
//
// Returns:
//   - models.Engine: The engine record if found, empty engine if not found
//   - error: Error for invalid ID format, database failures, or other issues
func (e EngineStore) GetEngineByID(ctx context.Context, id string) (models.Engine, error) {
	var engine models.Engine

	// Parse and validate the engine ID from string to UUID format
	// This ensures the ID is in correct UUID format before database operations
	engineId, err := uuid.Parse(id)
	if err != nil {
		return models.Engine{}, fmt.Errorf("invalid engine ID format: %v", err)
	}

	// Begin database transaction for data consistency and isolation
	// Transactions ensure ACID properties even for read operations
	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Engine{}, err
	}

	// Defer transaction cleanup with proper error handling
	// This ensures transaction is always properly closed regardless of execution path
	defer func() {
		if err != nil {
			// Rollback transaction on any error to maintain database consistency
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("Transaction rollback error: %v\n", rbErr)
			}
		} else {
			// Commit transaction on successful completion
			if cmErr := tx.Commit(); cmErr != nil {
				fmt.Printf("Transaction commit error: %v\n", cmErr)
			}
		}
	}()

	// Execute parameterized query to prevent SQL injection attacks
	// Query retrieves all engine fields based on the provided engine ID
	query := `SELECT engine_id, displacement, no_of_cylinders, car_range FROM engine WHERE engine_id = $1`
	err = tx.QueryRowContext(ctx, query, engineId).Scan(&engine.ID, &engine.Displacement, &engine.NoOfCylinders, &engine.CarRange)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return empty engine when no record found (not considered an error)
			return engine, nil
		}
		return models.Engine{}, err
	}

	return engine, nil
}

func (e EngineStore) CreateEngine(ctx context.Context, engineReq models.EngineRequest) (models.Engine, error) {
	newEngine := models.Engine{
		ID:            uuid.New(), // Generate a new UUID for the engine
		Displacement:  engineReq.Displacement,
		NoOfCylinders: engineReq.NoOfCylinders,
		CarRange:      engineReq.CarRange,
	}
	var createdEngine models.Engine
	// Begin the transaction
	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Engine{}, err // Return error if transaction cannot be started
	}
	// Insert the engine into the database
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("Transaction rollback error: %v\n", rbErr) // Log rollback error
			}
		} else {
			if cmErr := tx.Commit(); cmErr != nil {
				fmt.Printf("Transaction commit error: %v\n", cmErr) // Log commit error
			}
		}
	}()
	query := `INSERT INTO engine (engine_id, displacement, no_of_cylinders, car_range) VALUES ($1, $2, $3, $4)
	RETURNING engine_id, displacement, no_of_cylinders, car_range`
	err = tx.QueryRowContext(ctx, query, newEngine.ID, newEngine.Displacement, newEngine.NoOfCylinders, newEngine.CarRange).Scan(&createdEngine.ID, &createdEngine.Displacement, &createdEngine.NoOfCylinders, &createdEngine.CarRange)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Engine{}, errors.New("no rows returned from the query")
		}
		return models.Engine{}, err // Return error if query fails
	}
	return createdEngine, nil // Return the created engine
}

func (e EngineStore) UpdateEngine(ctx context.Context, id string, engineReq models.EngineRequest) (models.Engine, error) {
	engineId, err := uuid.Parse(id)
	if err != nil {
		return models.Engine{}, fmt.Errorf("invalid engine ID format: %v", err)
	}
	// Begin the transaction
	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Engine{}, err // Return error if transaction cannot be started
	}
	// Insert the engine into the database
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("Transaction rollback error: %v\n", rbErr) // Log rollback error
			}
		} else {
			if cmErr := tx.Commit(); cmErr != nil {
				fmt.Printf("Transaction commit error: %v\n", cmErr) // Log commit error
			}
		}
	}()

	query := `UPDATE engine SET displacement = $1, no_of_cylinders = $2, car_range = $3 WHERE engine_id = $4
	RETURNING engine_id, displacement, no_of_cylinders, car_range`
	results, err := tx.ExecContext(ctx, query, engineReq.Displacement, engineReq.NoOfCylinders, engineReq.CarRange, engineId)
	if err != nil {
		return models.Engine{}, fmt.Errorf("failed to update engine: %v", err)
	}
	if rowsAffected, _ := results.RowsAffected(); rowsAffected == 0 {
		return models.Engine{}, fmt.Errorf("no engine found with ID: %s", id)
	}
	var updatedEngine models.Engine = models.Engine{
		ID:            engineId,
		Displacement:  engineReq.Displacement,
		NoOfCylinders: engineReq.NoOfCylinders,
		CarRange:      engineReq.CarRange,
	}
	return updatedEngine, nil // Return the updated engine
}

func (e EngineStore) DeleteEngine(ctx context.Context, id string) (models.Engine, error) {
	var engine models.Engine

	// Parse the string ID to UUID
	engineId, err := uuid.Parse(id)
	if err != nil {
		return models.Engine{}, fmt.Errorf("invalid engine ID format: %v", err)
	}

	//  Begin the transaction
	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Engine{}, err // Return error if transaction cannot be started
	}
	// Insert the car into the database
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("Transaction rollback error: %v\n", rbErr) // Log rollback error
				// Rollback the transaction in case of error
			}
		} else {
			if cmErr := tx.Commit(); cmErr != nil {
				fmt.Printf("Transaction commit error: %v\n", cmErr) // Log commit error
			}
		}
	}()
	query := `SELECT  engine_id, displacement, no_of_cylinders, car_range FROM engine WHERE engine_id = $1`
	err = tx.QueryRowContext(ctx, query, engineId).Scan(&engine.ID, &engine.Displacement, &engine.NoOfCylinders, &engine.CarRange)
	if err != nil {
		if err == sql.ErrNoRows {
			return engine, nil // No engine found with the given ID
		}
		return models.Engine{}, err // Return error if query fails
	}
	deleteQuery := `DELETE FROM engine WHERE engine_id = $1`
	_, err = tx.ExecContext(ctx, deleteQuery, engineId)
	if err != nil {
		return models.Engine{}, fmt.Errorf("failed to delete engine: %v", err)
	}
	return engine, nil // Return the deleted engine
}

func (e EngineStore) GetEngineByBrand(ctx context.Context, brand string) ([]models.Engine, error) {
	var engines []models.Engine
	// Begin the transaction
	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err // Return error if transaction cannot be started
	}
	// Insert the car into the database
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("Transaction rollback error: %v\n", rbErr) // Log rollback error
			}
		} else {
			if cmErr := tx.Commit(); cmErr != nil {
				fmt.Printf("Transaction commit error: %v\n", cmErr) // Log commit error
			}
		}
	}()
	query := `SELECT engine_id, displacement, no_of_cylinders, car_range FROM engine WHERE brand = $1`
	rows, err := tx.QueryContext(ctx, query, brand)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch engines by brand: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var engine models.Engine
		if err := rows.Scan(&engine.ID, &engine.Displacement, &engine.NoOfCylinders, &engine.CarRange); err != nil {
			return nil, fmt.Errorf("failed to scan engine row: %v", err)
		}
		engines = append(engines, engine)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating over rows: %v", err)
	}
	return engines, nil // Return the found engines
}
