package car

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/PrateekKumar15/CarZone/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"go.opentelemetry.io/otel"
)

type CarStore struct {
	db *sql.DB
}

func New(db *sql.DB) CarStore {
	return CarStore{db: db}
}

func (s CarStore) GetCarByID(ctx context.Context, id string) (models.Car, error) {
	tracer := otel.Tracer("CarStore")
	ctx, span := tracer.Start(ctx, "GetCarByID-Store")
	defer span.End()

	var car models.Car
	var engineJSON, featuresJSON []byte
	var images pq.StringArray

	query := `SELECT id, owner_id, name, model, year, brand, fuel_type, engine, location_city, 
	         location_state, location_country, price, status, is_available, 
	         features, description, images, mileage, created_at, updated_at 
	         FROM car WHERE id = $1`

	row := s.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&car.ID, &car.OwnerID, &car.Name, &car.Model, &car.Year, &car.Brand,
		&car.FuelType, &engineJSON, &car.LocationCity, &car.LocationState, &car.LocationCountry,
		&car.Price, &car.Status, &car.IsAvailable, &featuresJSON,
		&car.Description, &images, &car.Mileage, &car.CreatedAt, &car.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Car{}, nil // No car found with the given ID
		}
		return models.Car{}, err
	}

	// Parse JSON fields
	if err = json.Unmarshal(engineJSON, &car.Engine); err != nil {
		return models.Car{}, err
	}

	if err = json.Unmarshal(featuresJSON, &car.Features); err != nil {
		return models.Car{}, err
	}
	car.Images = []string(images)

	return car, nil
}

// GetCarWithOwnerByID retrieves a car by ID and includes owner information
func (s CarStore) GetCarWithOwnerByID(ctx context.Context, id string) (models.Car, error) {
	tracer := otel.Tracer("CarStore")
	ctx, span := tracer.Start(ctx, "GetCarWithOwnerByID-Store")
	defer span.End()

	var car models.Car
	var owner models.User
	var engineJSON, featuresJSON, ownerProfileDataJSON []byte
	var images pq.StringArray

	// Join query to get car data with owner information (INNER JOIN since owner is mandatory)
	query := `SELECT 
		c.id, c.owner_id, c.name, c.model, c.year, c.brand, c.fuel_type, c.engine, 
		c.location_city, c.location_state, c.location_country, c.price, c.status, c.is_available, c.features, c.description, c.images, 
		c.mileage, c.created_at, c.updated_at,
		u.id, u.username, u.email, u.phone, u.role, u.profile_data, u.created_at, u.updated_at
		FROM car c 
		INNER JOIN users u ON c.owner_id = u.id 
		WHERE c.id = $1`

	row := s.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&car.ID, &car.OwnerID, &car.Name, &car.Model, &car.Year, &car.Brand,
		&car.FuelType, &engineJSON, &car.LocationCity, &car.LocationState, &car.LocationCountry,
		&car.Price, &car.Status, &car.IsAvailable, &featuresJSON,
		&car.Description, &images, &car.Mileage, &car.CreatedAt, &car.UpdatedAt,
		&owner.ID, &owner.UserName, &owner.Email, &owner.Phone, &owner.Role,
		&ownerProfileDataJSON, &owner.CreatedAt, &owner.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Car{}, nil // No car found with the given ID
		}
		return models.Car{}, err
	}

	// Parse JSON fields for car
	if err = json.Unmarshal(engineJSON, &car.Engine); err != nil {
		return models.Car{}, err
	}
	if err = json.Unmarshal(featuresJSON, &car.Features); err != nil {
		return models.Car{}, err
	}
	car.Images = []string(images)

	// Parse owner profile data (owner is mandatory)
	if len(ownerProfileDataJSON) > 0 {
		err = json.Unmarshal(ownerProfileDataJSON, &owner.ProfileData)
		if err != nil {
			return models.Car{}, err
		}
	} else {
		owner.ProfileData = make(map[string]interface{})
	}
	car.Owner = &owner

	return car, nil
}

func (s CarStore) GetCarByBrand(ctx context.Context, brand string) ([]models.Car, error) {
	tracer := otel.Tracer("CarStore")
	ctx, span := tracer.Start(ctx, "GetCarByBrand-Store")
	defer span.End()

	var cars []models.Car
	query := `SELECT id, owner_id, name, model, year, brand, fuel_type, engine, location_city, 
	         location_state, location_country, price, status, is_available, 
	         features, description, images, mileage, created_at, updated_at 
	         FROM car WHERE brand = $1`

	rows, err := s.db.QueryContext(ctx, query, brand)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var car models.Car
		var engineJSON, featuresJSON []byte
		var images pq.StringArray

		err = rows.Scan(&car.ID, &car.OwnerID, &car.Name, &car.Model, &car.Year, &car.Brand,
			&car.FuelType, &engineJSON, &car.LocationCity, &car.LocationState, &car.LocationCountry,
			&car.Price, &car.Status, &car.IsAvailable, &featuresJSON,
			&car.Description, &images, &car.Mileage, &car.CreatedAt, &car.UpdatedAt)

		if err != nil {
			return nil, err
		}

		// Parse JSON fields
		if err = json.Unmarshal(engineJSON, &car.Engine); err != nil {
			return nil, err
		}

		if err = json.Unmarshal(featuresJSON, &car.Features); err != nil {
			return nil, err
		}
		car.Images = []string(images)

		cars = append(cars, car)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cars, nil
}

func (s CarStore) CreateCar(ctx context.Context, carReq models.CarRequest) (models.Car, error) {
	tracer := otel.Tracer("CarStore")
	ctx, span := tracer.Start(ctx, "CreateCar-Store")
	defer span.End()

	var createdCar models.Car
	carId := uuid.New()
	createdAt := time.Now()
	updatedAt := createdAt

	// Marshal JSON fields
	engineJSON, err := json.Marshal(carReq.Engine)
	if err != nil {
		return models.Car{}, err
	}
	featuresJSON, err := json.Marshal(carReq.Features)
	if err != nil {
		return models.Car{}, err
	}
	images := pq.StringArray(carReq.Images)

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Car{}, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	query := `INSERT INTO car (id, owner_id, name, model, year, brand, fuel_type, engine, 
	         location_city, location_state, location_country, price, status, 
	         is_available, features, description, images, mileage, created_at, updated_at) 
	         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
	         RETURNING id, owner_id, name, model, year, brand, fuel_type, engine, location_city, 
	         location_state, location_country, price, status, is_available, 
	         features, description, images, mileage, created_at, updated_at`

	var returnedEngineJSON, returnedPriceJSON, returnedFeaturesJSON []byte
	var returnedImages pq.StringArray

	err = tx.QueryRowContext(ctx, query, carId, carReq.OwnerID, carReq.Name, carReq.Model, carReq.Year,
		carReq.Brand, carReq.FuelType, engineJSON, carReq.LocationCity, carReq.LocationState,
		carReq.LocationCountry, carReq.Price, carReq.Status, carReq.IsAvailable,
		featuresJSON, carReq.Description, images, carReq.Mileage, createdAt, updatedAt).Scan(
		&createdCar.ID, &createdCar.OwnerID, &createdCar.Name, &createdCar.Model, &createdCar.Year,
		&createdCar.Brand, &createdCar.FuelType, &returnedEngineJSON, &createdCar.LocationCity,
		&createdCar.LocationState, &createdCar.LocationCountry, &returnedPriceJSON, &createdCar.Status, &createdCar.IsAvailable, &returnedFeaturesJSON,
		&createdCar.Description, &returnedImages, &createdCar.Mileage, &createdCar.CreatedAt, &createdCar.UpdatedAt)

	if err != nil {
		return models.Car{}, err
	}

	// Parse returned JSON fields
	if err = json.Unmarshal(returnedEngineJSON, &createdCar.Engine); err != nil {
		return models.Car{}, err
	}
	if err = json.Unmarshal(returnedFeaturesJSON, &createdCar.Features); err != nil {
		return models.Car{}, err
	}
	createdCar.Images = []string(returnedImages)

	return createdCar, nil
}

func (s CarStore) UpdateCar(ctx context.Context, id string, carReq models.CarRequest) (models.Car, error) {
	tracer := otel.Tracer("CarStore")
	ctx, span := tracer.Start(ctx, "UpdateCar-Store")
	defer span.End()

	var updatedCar models.Car

	// Marshal JSON fields
	engineJSON, err := json.Marshal(carReq.Engine)
	if err != nil {
		return models.Car{}, err
	}

	featuresJSON, err := json.Marshal(carReq.Features)
	if err != nil {
		return models.Car{}, err
	}
	images := pq.StringArray(carReq.Images)

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Car{}, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	query := `UPDATE car SET owner_id = $1, name = $2, model = $3, year = $4, brand = $5, fuel_type = $6, 
	         engine = $7, location_city = $8, location_state = $9, location_country = $10, price = $11, 
	         status = $12, is_available = $14, features = $15, description = $16, 
	         images = $17, mileage = $18, updated_at = $19 WHERE id = $20 
	         RETURNING id, owner_id, name, model, year, brand, fuel_type, engine, location_city, 
	         location_state, location_country, price, status, availability_type, is_available, 
	         features, description, images, mileage, created_at, updated_at`

	var returnedEngineJSON, returnedPriceJSON, returnedFeaturesJSON []byte
	var returnedImages pq.StringArray

	err = tx.QueryRowContext(ctx, query, carReq.OwnerID, carReq.Name, carReq.Model, carReq.Year,
		carReq.Brand, carReq.FuelType, engineJSON, carReq.LocationCity, carReq.LocationState,
		carReq.LocationCountry,carReq.Price, carReq.Status, carReq.IsAvailable,
		featuresJSON, carReq.Description, images, carReq.Mileage, time.Now(), id).Scan(
		&updatedCar.ID, &updatedCar.OwnerID, &updatedCar.Name, &updatedCar.Model, &updatedCar.Year,
		&updatedCar.Brand, &updatedCar.FuelType, &returnedEngineJSON, &updatedCar.LocationCity,
		&updatedCar.LocationState, &updatedCar.LocationCountry, &returnedPriceJSON, &updatedCar.Status, &updatedCar.IsAvailable, &returnedFeaturesJSON,
		&updatedCar.Description, &returnedImages, &updatedCar.Mileage, &updatedCar.CreatedAt, &updatedCar.UpdatedAt)

	if err != nil {
		return models.Car{}, err
	}

	// Parse returned JSON fields
	if err = json.Unmarshal(returnedEngineJSON, &updatedCar.Engine); err != nil {
		return models.Car{}, err
	}
	if err = json.Unmarshal(returnedFeaturesJSON, &updatedCar.Features); err != nil {
		return models.Car{}, err
	}
	updatedCar.Images = []string(returnedImages)

	return updatedCar, nil
}

func (s CarStore) DeleteCar(ctx context.Context, id string) (models.Car, error) {
	tracer := otel.Tracer("CarStore")
	ctx, span := tracer.Start(ctx, "DeleteCar-Store")
	defer span.End()

	var deletedCar models.Car

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Car{}, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// First get the car data before deleting
	query := `SELECT id, owner_id, name, model, year, brand, fuel_type, engine, location_city, 
	         location_state, location_country, price, status, is_available, 
	         features, description, images, mileage, created_at, updated_at 
	         FROM car WHERE id = $1`

	var engineJSON, featuresJSON []byte
	var images pq.StringArray

	err = tx.QueryRowContext(ctx, query, id).Scan(&deletedCar.ID, &deletedCar.OwnerID, &deletedCar.Name,
		&deletedCar.Model, &deletedCar.Year, &deletedCar.Brand, &deletedCar.FuelType, &engineJSON,
		&deletedCar.LocationCity, &deletedCar.LocationState, &deletedCar.LocationCountry, &deletedCar.Price,
		&deletedCar.Status, &deletedCar.IsAvailable, &featuresJSON,
		&deletedCar.Description, &images, &deletedCar.Mileage, &deletedCar.CreatedAt, &deletedCar.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Car{}, errors.New("no car found with the given ID")
		}
		return models.Car{}, err
	}

	// Parse JSON fields
	if err = json.Unmarshal(engineJSON, &deletedCar.Engine); err != nil {
		return models.Car{}, err
	}
	deletedCar.Images = []string(images)

	// Now delete the car
	result, err := tx.ExecContext(ctx, "DELETE FROM car WHERE id = $1", id)
	if err != nil {
		return models.Car{}, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.Car{}, err
	}
	if rowsAffected == 0 {
		return models.Car{}, errors.New("no car found with the given ID")
	}

	return deletedCar, nil
}

func (s CarStore) GetAllCars(ctx context.Context) ([]models.Car, error) {
	tracer := otel.Tracer("CarStore")
	ctx, span := tracer.Start(ctx, "GetAllCars-Store")
	defer span.End()

	var cars []models.Car

	query := `SELECT id, owner_id, name, model, year, brand, fuel_type, engine, location_city, 
	         location_state, location_country, price, status, is_available, 
	         features, description, images, mileage, created_at, updated_at 
	         FROM car`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var car models.Car
		var engineJSON, featuresJSON []byte
		var images pq.StringArray

		err = rows.Scan(&car.ID, &car.OwnerID, &car.Name, &car.Model, &car.Year, &car.Brand,
			&car.FuelType, &engineJSON, &car.LocationCity, &car.LocationState, &car.LocationCountry,
			&car.Price, &car.Status, &car.IsAvailable, &featuresJSON,
			&car.Description, &images, &car.Mileage, &car.CreatedAt, &car.UpdatedAt)

		if err != nil {
			return nil, err
		}

		// Parse JSON fields
		if err = json.Unmarshal(engineJSON, &car.Engine); err != nil {
			return nil, err
		}
		if err = json.Unmarshal(featuresJSON, &car.Features); err != nil {
			return nil, err
		}
		car.Images = []string(images)

		cars = append(cars, car)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cars, nil
}
