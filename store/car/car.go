package car

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/PrateekKumar15/CarZone/models"
	"github.com/google/uuid"
)

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) Store {
	return Store{db: db}
}

func (s Store) GetCarByID(ctx context.Context, id string) (models.Car, error) {
	var car models.Car
	query := "SELECT c.id, c.name, c.year, c.brand, c.fuel_type, c.engine_id, c.price, c.created_at, c.updated_at, e.id, e.displacement, e.no_of_cylinders, e.car_range FROM car c LEFT JOIN engine e ON c.engine_id = e.id WHERE c.id = $1"
	row := s.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&car.ID, &car.Name, &car.Year, &car.Brand, &car.FuelType, &car.Price, &car.CreatedAt, &car.UpdatedAt, &car.Engine.EngineID, &car.Engine.Displacement, &car.Engine.NoOfCylinders, &car.Engine.CarRange)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return models.Car{}, nil // No car found with the given ID
		case sql.ErrConnDone:
			return models.Car{}, errors.New("database connection is closed")
		case sql.ErrTxDone:
			return models.Car{}, errors.New("transaction is already committed or rolled back")
		}
	}
	return car, nil // Return the found car
}

func (s Store) GetCarByBrand(ctx context.Context, brand string, isEngine bool) ([]models.Car, error) {
	var cars []models.Car
	var query string
	if isEngine {
		query = "SELECT c.id, c.name, c.year, c.brand, c.fuel_type, c.engine_id, c.price, c.created_at, c.updated_at, e.id, e.displacement, e.no_of_cylinders, e.car_range FROM car c JOIN engine e ON c.engine_id = e.id WHERE c.brand = $1"
	} else {
		query = "SELECT c.id, c.name, c.year, c.brand, c.fuel_type, c.engine_id, c.price, c.created_at, c.updated_at FROM car c WHERE c.brand = $1"
	}
	rows, err := s.db.QueryContext(ctx, query, brand)
	if err != nil {
		return nil, err // Return error if query fails
	}
	if err = rows.Err(); err != nil {
		return nil, err // Return error if there was an issue with the rows
	}
	defer rows.Close()
	for rows.Next() {
		var car models.Car
		if isEngine {
			err = rows.Scan(&car.ID, &car.Name, &car.Year, &car.Brand, &car.FuelType, &car.Price, &car.CreatedAt, &car.UpdatedAt, &car.Engine.EngineID, &car.Engine.Displacement, &car.Engine.NoOfCylinders, &car.Engine.CarRange)
			if err != nil {
				return nil, err // Return error if row scan fails
			}
		} else {
			err = rows.Scan(&car.ID, &car.Name, &car.Year, &car.Brand, &car.FuelType, &car.Price, &car.CreatedAt, &car.UpdatedAt)
		}
		if err != nil {
			return nil, err // Return error if row scan fails
		}
		cars = append(cars, car)

	}

	return cars, nil // Return the list of cars
}

func (s Store) CreateCar(ctx context.Context, carReq models.CarRequest) (models.Car, error) {
	var createdCar models.Car
	var engineID uuid.UUID

	err := s.db.QueryRowContext(ctx, "Select id from engine where id = $1", carReq.Engine.EngineID).Scan(&engineID)
	if err!=nil {
		if err == sql.ErrNoRows {
			return models.Car{}, errors.New("engine id not found")
		}
		return  createdCar, err
	}
	carId := uuid.New()
	createdAt := time.Now()
	updatedAt := createdAt
	newCar := models.Car{
		ID:        carId,
		Name:      carReq.Name,
		Year:      carReq.Year,
		Brand:     carReq.Brand,
		FuelType:  carReq.FuelType,
		Engine:    carReq.Engine,
		Price:     carReq.Price,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
	//  Begin the transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Car{}, err // Return error if transaction cannot be started
	}	
	// Insert the car into the database
	defer func () {
		if err != nil {
			tx.Rollback() // Rollback the transaction in case of error
			return
		}
		err = tx.Commit() // Commit the transaction if no error
		if err != nil {
			return
		}
	}()
	query := `INSERT INTO car (id, name, year, brand, fuel_type, engine_id, price, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id, name, year, brand, fuel_type, engine_id, price, created_at, updated_at`
		err = tx.QueryRowContext(ctx, query, newCar.ID, newCar.Name, newCar.Year, newCar.Brand, newCar.FuelType, engineID, newCar.Price, newCar.CreatedAt, newCar.UpdatedAt).Scan(&createdCar.ID, &createdCar.Name, &createdCar.Year, &createdCar.Brand, &createdCar.FuelType, &createdCar.Engine, &createdCar.Price, &createdCar.CreatedAt, &createdCar.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Car{}, errors.New("no rows returned from the query")
		}
		return models.Car{}, err // Return error if query fails
	}
	return createdCar, nil // Return the created car
}

func (s Store) UpdateCar(ctx context.Context, id string, carReq models.CarRequest) (models.Car, error) {
 var updatedCar models.Car
 //  Begin the transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Car{}, err // Return error if transaction cannot be started
	}	
	// Insert the car into the database
	defer func () {
		if err != nil {
			tx.Rollback() // Rollback the transaction in case of error
			return
		}
		err = tx.Commit() // Commit the transaction if no error
		if err != nil {
			return
		}
	}()
	query := `UPDATE car SET name = $1, year = $2, brand = $3, fuel_type = $4, engine_id = $5, price = $6, updated_at = $7 WHERE id = $8 RETURNING id, name, year, brand, fuel_type, engine_id, price, created_at, updated_at`
	err = tx.QueryRowContext(ctx, query, carReq.Name, carReq.Year, carReq.Brand, carReq.FuelType, carReq.Engine.EngineID, carReq.Price, time.Now(), id).Scan(&updatedCar.ID, &updatedCar.Name, &updatedCar.Year, &updatedCar.Brand, &updatedCar.FuelType, &updatedCar.Engine.EngineID, &updatedCar.Price, &updatedCar.CreatedAt, &updatedCar.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Car{}, errors.New("no rows returned from the query")
		}
		return models.Car{}, err // Return error if query fails
	}
	return updatedCar, nil // Return the updated car

}

func (s Store) DeleteCar(ctx context.Context, id string) (models.Car,error) {
	var deletedCar models.Car
	 //  Begin the transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Car{}, err // Return error if transaction cannot be started
	}	
	// Insert the car into the database
	defer func () {
		if err != nil {
			tx.Rollback() // Rollback the transaction in case of error
			return
		}
		err = tx.Commit() // Commit the transaction if no error
		if err != nil {
			return
		}
	}()
	query := `Select id, name, year, brand, fuel_type, engine_id, price, created_at, updated_at FROM car WHERE id = $1`
	err = tx.QueryRowContext(ctx, query, id).Scan(&deletedCar.ID, &deletedCar.Name, &deletedCar.Year, &deletedCar.Brand, &deletedCar.FuelType, &deletedCar.Engine.EngineID, &deletedCar.Price, &deletedCar.CreatedAt, &deletedCar.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Car{},errors.New("no car found with the given ID")
		}
		return models.Car{}, err // Return error if query fails

}

	result, err := tx.ExecContext(ctx, "DELETE FROM car WHERE id = $1", id)
	if err != nil {
		return models.Car{}, err // Return error if delete query fails
}
 rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.Car{}, err // Return error if rows affected query fails
	}
	if rowsAffected == 0 {
		return models.Car{}, errors.New("no car found with the given ID")
	}
	return deletedCar, nil // Return nil if the car was successfully deleted
	}