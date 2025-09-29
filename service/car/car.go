package car

import (
	"context"
	"errors"

	"github.com/PrateekKumar15/CarZone/models"
	"github.com/PrateekKumar15/CarZone/store"
	"go.opentelemetry.io/otel"
)

type CarService struct {
	store store.CarStoreInterface
}

func NewCarService(store store.CarStoreInterface) *CarService {
	return &CarService{store: store}
}

func (s *CarService) GetCarByID(ctx context.Context, id string) (*models.Car, error) {
	tracer := otel.Tracer("CarService")
	ctx, span := tracer.Start(ctx, "GetCarByID-Service")
	defer span.End()

	if id == "" {
		return nil, errors.New("car ID cannot be empty")
	}

	// Use the method that includes owner information
	car, err := s.store.GetCarWithOwnerByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if car was found (store returns empty car if not found)
	if car.ID.String() == "00000000-0000-0000-0000-000000000000" {
		return nil, nil
	}

	return &car, nil
}

func (s *CarService) GetCarByBrand(ctx context.Context, brand string) (*[]models.Car, error) {
	tracer := otel.Tracer("CarService")
	ctx, span := tracer.Start(ctx, "GetCarByBrand-Service")
	defer span.End()

	if brand == "" {
		return nil, errors.New("brand cannot be empty")
	}

	cars, err := s.store.GetCarByBrand(ctx, brand)
	if err != nil {
		return nil, err
	}

	return &cars, nil
}

func (s *CarService) CreateCar(ctx context.Context, carReq models.CarRequest) (*models.Car, error) {
	tracer := otel.Tracer("CarService")
	ctx, span := tracer.Start(ctx, "CreateCar-Service")
	defer span.End()

	// Validate the car request
	if err := s.validateCarRequest(carReq); err != nil {
		return nil, err
	}

	createdCar, err := s.store.CreateCar(ctx, carReq)
	if err != nil {
		return nil, err
	}

	return &createdCar, nil
}

func (s *CarService) UpdateCar(ctx context.Context, id string, carReq models.CarRequest) (*models.Car, error) {
	tracer := otel.Tracer("CarService")
	ctx, span := tracer.Start(ctx, "UpdateCar-Service")
	defer span.End()

	if id == "" {
		return nil, errors.New("car ID cannot be empty")
	}

	// Validate the car request
	if err := s.validateCarRequest(carReq); err != nil {
		return nil, err
	}

	updatedCar, err := s.store.UpdateCar(ctx, id, carReq)
	if err != nil {
		return nil, err
	}

	return &updatedCar, nil
}
func (s *CarService) DeleteCar(ctx context.Context, id string) (*models.Car, error) {
	tracer := otel.Tracer("CarService")
	ctx, span := tracer.Start(ctx, "DeleteCar-Service")
	defer span.End()

	if id == "" {
		return nil, errors.New("car ID cannot be empty")
	}

	deletedCar, err := s.store.DeleteCar(ctx, id)
	if err != nil {
		return nil, err
	}

	return &deletedCar, nil
}

func (s *CarService) GetAllCars(ctx context.Context) (*[]models.Car, error) {
	tracer := otel.Tracer("CarService")
	ctx, span := tracer.Start(ctx, "GetAllCars-Service")
	defer span.End()
	cars, err := s.store.GetAllCars(ctx)
	if err != nil {
		return nil, err // Return error if fetching all cars fails
	}
	return &cars, nil // Return the list of all cars
}

// validateCarRequest validates the car request data
func (s *CarService) validateCarRequest(carReq models.CarRequest) error {
	if carReq.Name == "" {
		return errors.New("car name is required")
	}
	if carReq.Model == "" {
		return errors.New("car model is required")
	}
	if carReq.Year < 1900 || carReq.Year > 2030 {
		return errors.New("invalid car year")
	}
	if carReq.Brand == "" {
		return errors.New("car brand is required")
	}
	if carReq.FuelType == "" {
		return errors.New("fuel type is required")
	}
	if carReq.LocationCity == "" {
		return errors.New("location city is required")
	}
	if carReq.LocationState == "" {
		return errors.New("location state is required")
	}
	if carReq.LocationCountry == "" {
		return errors.New("location country is required")
	}
	if carReq.Status == "" {
		return errors.New("car status is required")
	}
	if carReq.AvailabilityType == "" {
		return errors.New("availability type is required")
	}

	// Validate engine data
	if carReq.Engine.EngineSize <= 0 {
		return errors.New("engine size must be greater than 0")
	}
	if carReq.Engine.Cylinders <= 0 {
		return errors.New("number of cylinders must be greater than 0")
	}
	if carReq.Engine.Horsepower <= 0 {
		return errors.New("engine horsepower must be greater than 0")
	}
	if carReq.Engine.Transmission == "" {
		return errors.New("transmission type is required")
	}

	// Validate price data
	if carReq.AvailabilityType == "rental" || carReq.AvailabilityType == "both" {
		if carReq.Price.RentalPriceDaily <= 0 {
			return errors.New("rental price daily must be specified for rental cars")
		}
	}
	if carReq.AvailabilityType == "sale" || carReq.AvailabilityType == "both" {
		if carReq.Price.SalePrice == nil || *carReq.Price.SalePrice <= 0 {
			return errors.New("sale price must be specified for cars for sale")
		}
	}

	return nil
}
