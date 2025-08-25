package engine

import (
	"context"
	"github.com/PrateekKumar15/CarZone/models"
	"github.com/PrateekKumar15/CarZone/store"
	"go.opentelemetry.io/otel"
)

type EngineService struct {
	store store.EngineStoreInterface
}
func NewEngineService(store store.EngineStoreInterface) *EngineService {
	return &EngineService{store: store}
}
func (s *EngineService) GetEngineByID(ctx context.Context, id string) (*models.Engine, error) {
	tracer := otel.Tracer("EngineService")
	ctx, span := tracer.Start(ctx, "GetEngineByID-Service")
	defer span.End()
	engine, err := s.store.GetEngineByID(ctx, id)
	if err != nil {
		return nil, err // Return error if fetching engine fails
	}
	return &engine, nil // Return the found engine
}
func (s *EngineService) CreateEngine(ctx context.Context, engineReq models.EngineRequest) (*models.Engine, error) {
	tracer := otel.Tracer("EngineService")
	ctx, span := tracer.Start(ctx, "CreateEngine-Service")	
	defer span.End()
	if err := models.ValidateEngineRequest(engineReq); err != nil {
		return nil, err // Return error if validation fails
	}
	createdEngine, err := s.store.CreateEngine(ctx, engineReq)
	if err != nil {
		return nil, err // Return error if creating engine fails
	}
	return &createdEngine, nil // Return the created engine
}
func (s *EngineService) UpdateEngine(ctx context.Context, id string, engineReq models.EngineRequest) (*models.Engine, error) {
	tracer := otel.Tracer("EngineService")
	ctx, span := tracer.Start(ctx, "UpdateEngine-Service")
	defer span.End()
	if err := models.ValidateEngineRequest(engineReq); err != nil {
		return nil, err // Return error if validation fails
	}
	updatedEngine, err := s.store.UpdateEngine(ctx, id, engineReq)
	if err != nil {
		return nil, err // Return error if updating engine fails
	}
	return &updatedEngine, nil // Return the updated engine
}
func (s *EngineService) DeleteEngine(ctx context.Context, id string) (*models.Engine, error) {
	tracer := otel.Tracer("EngineService")
	ctx, span := tracer.Start(ctx, "DeleteEngine-Service")
	defer span.End()
	deletedEngine, err := s.store.DeleteEngine(ctx, id)
	if err != nil {
		return nil, err // Return error if deleting engine fails
	}
	return &deletedEngine, nil // Return the deleted engine
}

