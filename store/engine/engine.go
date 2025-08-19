package engine

import (
	"context"
	"database/sql"
	"github.com/PrateekKumar15/CarZone/models"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type EngineStore struct {
	db *sql.DB
}
func New(db *sql.DB) *EngineStore {
	return &EngineStore{db: db}
}

func (e EngineStore) GetEngineByID(ctx context.Context, id string) (models.Engine, error) {
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
	defer func () {
		if err != nil {
			if rbErr:=tx.Rollback(); rbErr!=nil {
			fmt.Printf("Transaction rollback error: %v\n", rbErr) // Log rollback error
			 // Rollback the transaction in case of error
			}
		}else {if cmErr := tx.Commit(); cmErr != nil {
			fmt.Printf("Transaction commit error: %v\n", cmErr) // Log commit error
		}}
	}()
	query := `SELECT  engine_id, displacement, no_of_cylinders, car_range FROM engine WHERE engine_id = $1`
	err = tx.QueryRowContext(ctx, query, engineId).Scan( &engine.ID, &engine.Displacement, &engine.NoOfCylinders, &engine.CarRange)
	if err != nil {
		if err == sql.ErrNoRows {
			return engine, nil // No engine found with the given ID
		}
		return models.Engine{}, err // Return error if query fails
	}
	return engine, nil // Return the found engine
}

func (e EngineStore) CreateEngine(ctx context.Context, engineReq models.EngineRequest) (models.Engine, error) {
	newEngine := models.Engine{
		ID: uuid.New(), // Generate a new UUID for the engine
		Displacement: engineReq.Displacement,
		NoOfCylinders: engineReq.NoOfCylinders,
		CarRange: engineReq.CarRange,
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
	err = tx.QueryRowContext(ctx, query,newEngine.ID,newEngine.Displacement,newEngine.NoOfCylinders,newEngine.CarRange).Scan(&createdEngine.ID, &createdEngine.Displacement, &createdEngine.NoOfCylinders, &createdEngine.CarRange)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Engine{}, errors.New("no rows returned from the query")
		}
		return models.Engine{}, err // Return error if query fails
	}
	return createdEngine, nil // Return the created engine
}

func (e EngineStore) EngineUpdate(ctx context.Context, id string, engineReq models.EngineRequest) (models.Engine, error) {
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
	results,err := tx.ExecContext(ctx, query, engineReq.Displacement, engineReq.NoOfCylinders, engineReq.CarRange, engineId)
	if err != nil {
		return models.Engine{}, fmt.Errorf("failed to update engine: %v", err)
	}
	if rowsAffected, _ := results.RowsAffected(); rowsAffected == 0 {
		return models.Engine{}, fmt.Errorf("no engine found with ID: %s", id)
	}
	var updatedEngine models.Engine = models.Engine{
		ID:      engineId,
		Displacement:  engineReq.Displacement,
		NoOfCylinders: engineReq.NoOfCylinders,
		CarRange:      engineReq.CarRange,
	}
	return updatedEngine, nil // Return the updated engine
}

func (e EngineStore) DeleteEngine(ctx context.Context, id string) (models.Engine,error) {
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
	defer func () {
		if err != nil {
			if rbErr:=tx.Rollback(); rbErr!=nil {
			fmt.Printf("Transaction rollback error: %v\n", rbErr) // Log rollback error
			 // Rollback the transaction in case of error
			}
		}else {if cmErr := tx.Commit(); cmErr != nil {
			fmt.Printf("Transaction commit error: %v\n", cmErr) // Log commit error
		}}
	}()
	query := `SELECT  engine_id, displacement, no_of_cylinders, car_range FROM engine WHERE engine_id = $1`
	err = tx.QueryRowContext(ctx, query, engineId).Scan( &engine.ID, &engine.Displacement, &engine.NoOfCylinders, &engine.CarRange)
	if err != nil {
		if err == sql.ErrNoRows {
			return engine, nil // No engine found with the given ID
		}
		return models.Engine{}, err // Return error if query fails
	}
	deleteQuery := `DELETE FROM engine WHERE engine_id = $1`
	_, err = tx.ExecContext(ctx, deleteQuery, engineId)
	if err != nil {		return models.Engine{}, fmt.Errorf("failed to delete engine: %v", err)
	}
	return engine, nil // Return the deleted engine
}

