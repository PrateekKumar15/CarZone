package car

import (
	"log"
	"net/http"
	"encoding/json"
	"github.com/PrateekKumar15/CarZone/models"
	"github.com/PrateekKumar15/CarZone/service"
	"github.com/gorilla/mux"
	"io"
)

// CarHandler struct to handle car-related requests
type CarHandler struct {
	service service.CarServiceInterface
}
// NewCarHandler creates a new CarHandler with the provided service
func NewCarHandler(service service.CarServiceInterface) *CarHandler {
	return &CarHandler{service: service}
}

// GetCarByID retrieves a car by its ID
func (h *CarHandler) GetCarByID(w http.ResponseWriter, r *http.Request)  {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]
	resp, err := h.service.GetCarByID(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error retrieving car by ID:", err)
		return
	}
	if resp == nil {
		http.Error(w, "Car not found", http.StatusNotFound)
		return
	}
	body,err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling response:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_,err = w.Write(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error writing response:", err)
		return
	}
}

func (h *CarHandler) GetCarByBrand(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	brand := r.URL.Query().Get("brand")
	isEngine := r.URL.Query().Get("isEngine") == "true"
	
	resp,err := h.service.GetCarByBrand(ctx, brand, isEngine)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error retrieving car by brand:", err)
		return
	}
	body,err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling response:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_,err = w.Write(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error writing response:", err)
		return
	}
}

func (h *CarHandler) CreateCar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error reading request body:", err)
		return
	}
	var carRequest models.CarRequest
	err = json.Unmarshal(body, &carRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error unmarshalling request body:", err)
		return
	}
	if err := models.ValidateRequest(carRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Validation error:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	createdCar, err := h.service.CreateCar(ctx, carRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error creating car:", err)
		return
	}
	createdCarJSON, err := json.Marshal(createdCar)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling response:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// Write the created car JSON to the response
	_,_ = w.Write(createdCarJSON)
	
}

func (h *CarHandler) UpdateCar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error reading request body:", err)
		return
	}
	var carRequest models.CarRequest
	err = json.Unmarshal(body, &carRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error unmarshalling request body:", err)
		return
	}
	if err := models.ValidateRequest(carRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Validation error:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	updatedCar, err := h.service.UpdateCar(ctx, id, carRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error updating car:", err)
		return
	}
	updatedCarJSON, err := json.Marshal(updatedCar)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling response:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	// Write the updated car JSON to the response
	_,_ = w.Write(updatedCarJSON)
	
}

func (h *CarHandler) DeleteCar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]
	deletedCar,err := h.service.DeleteCar(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error deleting car:", err)
		return
	}
	if deletedCar == nil {
		http.Error(w, "Car not found", http.StatusNotFound)
		return
	}
	// No need to return the deleted car, just a success status
	body,err := json.Marshal(deletedCar)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling response:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) 
	_,_ = w.Write(body)
	
}