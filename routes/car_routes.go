package routes

import (
	"net/http"

	"github.com/PrateekKumar15/CarZone/middleware"
	"github.com/gorilla/mux"
)

// setupCarRoutes configures all car-related routes
func (r *Router) setupCarRoutes(router *mux.Router) {
	// Car CRUD operations

	// GET /cars - Retrieve all cars with optional filtering
	// Query parameters: ?brand=Toyota&fuel_type=Petrol&location=California
	router.HandleFunc("/cars", r.CarHandler.GetAllCars).Methods("GET", "OPTIONS")

	// GET /cars/{id} - Retrieve a specific car by its UUID
	// Path parameter: UUID of the car
	router.HandleFunc("/cars/{id}", r.CarHandler.GetCarByID).Methods("GET", "OPTIONS")

	// GET /cars/brand - Retrieve cars by brand with optional engine details
	// Query parameters: ?brand={brand}&engine={true/false}
	router.HandleFunc("/carsbybrand", r.CarHandler.GetCarByBrand).Methods("GET")

	// POST /cars - Create a new car record
	// Body: Car JSON data, supports multipart/form-data for image uploads
	router.Handle("/cars", middleware.ImageUploadMiddleware(http.HandlerFunc(r.CarHandler.CreateCar))).Methods("POST", "OPTIONS")

	// PUT /cars/{id} - Update an existing car by its UUID
	// Path parameter: UUID of the car to update
	// Body: Updated car JSON data, supports multipart/form-data for image uploads
	router.Handle("/cars/{id}", middleware.ImageUploadMiddleware(http.HandlerFunc(r.CarHandler.UpdateCar))).Methods("PUT", "OPTIONS")

	// DELETE /cars/{id} - Delete a car by its UUID
	// Path parameter: UUID of the car to delete
	router.HandleFunc("/cars/{id}", r.CarHandler.DeleteCar).Methods("DELETE", "OPTIONS")
}
