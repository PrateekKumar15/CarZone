package engine

import (
	"encoding/json"
	"fmt"
	"github.com/PrateekKumar15/CarZone/models"
	"github.com/PrateekKumar15/CarZone/service"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"io"
	"net/http"
)

type EngineHandler struct {
	service service.EngineServiceInterface
}

func NewEngineHandler(service service.EngineServiceInterface) *EngineHandler {
	return &EngineHandler{service: service}
}

func (h *EngineHandler) GetEngineByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer("EngineHandler")
	ctx, span := tracer.Start(ctx, "GetEngineByID-Handler")
	defer span.End()
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	engine, err := h.service.GetEngineByID(ctx, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching engine: %v", err), http.StatusInternalServerError)
		return
	}

	if engine == nil {
		http.Error(w, "Engine not found", http.StatusNotFound)
		return
	}
	body, err := json.Marshal(engine)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error marshalling response: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error writing response: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *EngineHandler) CreateEngine(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer("EngineHandler")
	ctx, span := tracer.Start(ctx, "CreateEngine-Handler")
	defer span.End()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading request body: %v", err), http.StatusBadRequest)
		return
	}
	var engineReq models.EngineRequest
	if err := json.Unmarshal(body, &engineReq); err != nil {
		http.Error(w, fmt.Sprintf("Error unmarshalling request body: %v", err), http.StatusBadRequest)
		return
	}
	createdEngine, err := h.service.CreateEngine(ctx, engineReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating engine: %v", err), http.StatusInternalServerError)
		return
	}
	createdEngineJSON, err := json.Marshal(createdEngine)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error marshalling response: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(createdEngineJSON)
}

func (h *EngineHandler) UpdateEngine(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer("EngineHandler")
	ctx, span := tracer.Start(ctx, "UpdateEngine-Handler")
	defer span.End()
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading request body: %v", err), http.StatusBadRequest)
		return
	}
	var engineReq models.EngineRequest
	if err := json.Unmarshal(body, &engineReq); err != nil {
		http.Error(w, fmt.Sprintf("Error unmarshalling request body: %v", err), http.StatusBadRequest)
		return
	}
	updatedEngine, err := h.service.UpdateEngine(ctx, id, engineReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating engine: %v", err), http.StatusInternalServerError)
		return
	}
	updatedEngineJSON, err := json.Marshal(updatedEngine)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error marshalling response: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(updatedEngineJSON)

}

func (h *EngineHandler) DeleteEngine(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer("EngineHandler")
	ctx, span := tracer.Start(ctx, "DeleteEngine-Handler")
	defer span.End()
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}
	deletedEngine, err := h.service.DeleteEngine(ctx, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting engine: %v", err), http.StatusInternalServerError)
		return
	}
	body, err := json.Marshal(deletedEngine)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error marshalling response: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)

}
