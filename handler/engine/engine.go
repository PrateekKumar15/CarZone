package engine

import (
	"fmt"
	"net/http"
	"github.com/PrateekKumar15/CarZone/service"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/PrateekKumar15/CarZone/models"
	"io"
)	

type EngineHandler struct {
	service  service.EngineServiceInterface
}
func NewEngineHandler(service service.EngineServiceInterface) *EngineHandler {
	return &EngineHandler{service: service}
}

func (h *EngineHandler) GetEngineByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
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
	body,err := json.Marshal(engine)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error marshalling response: %v", err), http.StatusInternalServerError)
		return
	}
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)
	_,err = w.Write(body)	
	if err != nil {
		http.Error(w, fmt.Sprintf("Error writing response: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *EngineHandler) CreateEngine(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body,err := io.ReadAll(r.Body)
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

func (h* EngineHandler) GetEngineByBrand(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	brand := vars["brand"]
	if brand == "" {
		http.Error(w, "Brand is required", http.StatusBadRequest)
		return
	}
	engines, err := h.service.GetEngineByBrand(ctx, brand)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching engines: %v", err), http.StatusInternalServerError)
		return
	}
	body, err := json.Marshal(engines)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error marshalling response: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}


