package engine

import (
	"fmt"
	"net/http"
	"github.com/PrateekKumar15/CarZone/service"
	"encoding/json"
	"github.com/gorilla/mux"
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
	_,_ = w.Write(body)	

}





