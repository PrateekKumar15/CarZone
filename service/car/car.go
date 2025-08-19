package car

import (
	"context"

	"github.com/PrateekKumar15/CarZone/models"
	"github.com/PrateekKumar15/CarZone/store"
)

type CarService struct {
	store store.CarStoreInterface
}

func NewCarService(store store.CarStoreInterface) *CarService {
	return &CarService{store: store}
}

func (s *CarService) GetCarByID(ctx context.Context, id string) (*models.Car, error) {
	car, err := s.store.GetCarByID(ctx, id)
	if err != nil {
		return nil, err // Return error if fetching car fails
	}
	return &car, nil // Return the found car
}

func (s *CarService) GetCarByBrand(ctx context.Context, brand string, isEngine bool) (*[]models.Car, error) {
	cars, err := s.store.GetCarByBrand(ctx, brand, isEngine)
	if err != nil {
		return nil, err // Return error if fetching cars by brand fails
	}
	return &cars, nil // Return the found cars
}

func (s *CarService) CreateCar(ctx context.Context, carReq models.CarRequest) (*models.Car, error) {
	if err := models.ValidateRequest(carReq); err != nil {
		return nil, err // Return error if validation fails
	}
	createdCar, err := s.store.CreateCar(ctx, carReq)
	if err != nil {
		return nil, err // Return error if creating car fails
	}
	return &createdCar, nil // Return the created car
}

func (s *CarService) UpdateCar(ctx context.Context, id string, carReq models.CarRequest) (*models.Car, error) {
	if err := models.ValidateRequest(carReq); err != nil {
		return nil, err // Return error if validation fails
	}
	updatedCar, err := s.store.UpdateCar(ctx, id, carReq)
	if err != nil {
		return nil, err // Return error if updating car fails
	}
	return &updatedCar, nil // Return the updated car
}
func (s *CarService) DeleteCar(ctx context.Context, id string) (*models.Car, error) {
	deletedCar, err := s.store.DeleteCar(ctx, id)
	if err != nil {
		return nil, err // Return error if deleting car fails
	}
	return &deletedCar, nil // Return the deleted car
}
