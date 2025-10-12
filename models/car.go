// Package models contains the data structures and validation logic for the CarZone application
package models

import (
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// Engine represents the engine specifications embedded within a car
type Engine struct {
	EngineSize   float64 `json:"engine_size"`  // Engine displacement in liters
	Cylinders    int     `json:"cylinders"`    // Number of cylinders
	Horsepower   int     `json:"horsepower"`   // Engine horsepower
	Transmission string  `json:"transmission"` // Manual, Automatic, CVT, Semi-Automatic
}

// Price represents the pricing information for a car rental

// Car represents a vehicle entity in the CarZone rental and sales system
// It contains all necessary information for rental management including
// ownership, pricing, availability, and specifications
type Car struct {
	ID       uuid.UUID  `json:"id"`              // Unique identifier for the car
	OwnerID  *uuid.UUID `json:"owner_id"`        // ID of the user who owns this car
	Owner    *User      `json:"owner,omitempty"` // Owner user information (populated when needed)
	Name     string     `json:"name"`            // Display name/model of the car
	Brand    string     `json:"brand"`           // Manufacturer brand name
	Model    string     `json:"model"`           // Specific model name
	Year     int        `json:"year"`            // Manufacturing year
	FuelType string     `json:"fuel_type"`       // Type of fuel (Petrol, Diesel, Electric, Hybrid)

	// Engine specifications (embedded struct)
	Engine Engine `json:"engine"` // Engine specifications

	// Location information
	LocationCity    string `json:"location_city"`    // City where car is located
	LocationState   string `json:"location_state"`   // State/province where car is located
	LocationCountry string `json:"location_country"` // Country where car is located

	// Pricing (embedded struct)
	Price float64 `json:"rental_price"` // Pricing information

	// Status and availability
	Status      string `json:"status"`       // active, maintenance, inactive
	IsAvailable bool   `json:"is_available"` // Current availability status

	// Additional information
	Features    map[string]interface{} `json:"features"`    // Car features as JSON (GPS, AC, etc.)
	Description string                 `json:"description"` // Detailed description
	Images      []string               `json:"images"`      // Array of image URLs
	Mileage     int                    `json:"mileage"`     // Current mileage

	// Timestamps
	CreatedAt time.Time `json:"created_at"` // When the car record was created
	UpdatedAt time.Time `json:"updated_at"` // When the car record was last updated
}

// CarRequest represents the data structure for creating or updating a car
// It contains all necessary fields for car creation/update but excludes system-generated fields
type CarRequest struct {
	OwnerID  *uuid.UUID `json:"owner_id"`  // ID of the user who owns this car
	Name     string     `json:"name"`      // Display name/model of the car
	Brand    string     `json:"brand"`     // Manufacturer brand name
	Model    string     `json:"model"`     // Specific model name
	Year     int        `json:"year"`      // Manufacturing year
	FuelType string     `json:"fuel_type"` // Type of fuel

	// Engine specifications (embedded struct)
	Engine Engine `json:"engine"` // Engine specifications

	// Location information
	LocationCity    string `json:"location_city"`    // City where car is located
	LocationState   string `json:"location_state"`   // State/province
	LocationCountry string `json:"location_country"` // Country

	// Pricing (embedded struct)
	Price float64 `json:"rental_price"` // Pricing information

	// Status and availability
	Status      string `json:"status"`       // active, maintenance, inactive
	IsAvailable bool   `json:"is_available"` // Current availability

	// Additional information
	Features    map[string]interface{} `json:"features"`    // Car features as JSON
	Description string                 `json:"description"` // Detailed description
	Images      []string               `json:"images"`      // Array of image URLs
	Mileage     int                    `json:"mileage"`     // Current mileage
}

// ValidateRequest performs comprehensive validation on a CarRequest
// It validates all fields including name, year, brand, fuel type, engine specs, and pricing
// Returns an error if any validation fails, nil if all validations pass
func ValidateRequest(carRequest CarRequest) error {
	if err := validateName(carRequest.Name); err != nil {
		return err
	}
	if err := validateYear(strconv.Itoa(carRequest.Year)); err != nil {
		return err
	}
	if err := validateBrand(carRequest.Brand); err != nil {
		return err
	}
	if err := validateModel(carRequest.Model); err != nil {
		return err
	}
	if err := validateFuelType(carRequest.FuelType); err != nil {
		return err
	}
	if err := validateEngine(carRequest.Engine); err != nil {
		return err
	}
	if err := validateLocation(carRequest.LocationCity, carRequest.LocationState, carRequest.LocationCountry); err != nil {
		return err
	}
	if err := validatePrice(carRequest.Price); err != nil {
		return err
	}
	if err := validateStatus(carRequest.Status); err != nil {
		return err
	}
	if err := validateMileage(carRequest.Mileage); err != nil {
		return err
	}
	return nil
}

// validateName checks if the car name meets the minimum length requirement
func validateName(name string) error {
	if len(name) < 3 {
		return errors.New("name must be at least 3 characters long")
	}
	return nil
}

// validateBrand checks if the car brand name meets the minimum length requirement
func validateBrand(brand string) error {
	if len(brand) < 2 {
		return errors.New("brand must be at least 2 characters long")
	}
	return nil
}

// validateModel checks if the car model name is valid
func validateModel(model string) error {
	if len(model) < 1 {
		return errors.New("model cannot be empty")
	}
	return nil
}

// validateYear validates the manufacturing year of the car
func validateYear(year string) error {
	if year == "" {
		return errors.New("year cannot be empty")
	}

	yearInt, err := strconv.Atoi(year)
	if err != nil {
		return errors.New("year must be a valid number")
	}

	currentYear := time.Now().Year()
	if yearInt < 1886 || yearInt > currentYear {
		return errors.New("year must be between 1886 and the current year")
	}

	return nil
}

// validateFuelType ensures the fuel type is one of the accepted values
func validateFuelType(fuelType string) error {
	validFuelTypes := []string{"Petrol", "Diesel", "Electric", "Hybrid", "CNG", "LPG"}
	for _, validType := range validFuelTypes {
		if fuelType == validType {
			return nil
		}
	}
	return errors.New("fuel type must be one of: Petrol, Diesel, Electric, Hybrid, CNG, LPG")
}

// validateTransmission ensures the transmission type is valid
func validateTransmission(transmission string) error {
	validTransmissions := []string{"Manual", "Automatic", "CVT", "Semi-Automatic"}
	for _, validType := range validTransmissions {
		if transmission == validType {
			return nil
		}
	}
	return errors.New("transmission must be one of: Manual, Automatic, CVT, Semi-Automatic")
}

// validateEngine validates the engine struct and all its fields
func validateEngine(engine Engine) error {
	if err := validateTransmission(engine.Transmission); err != nil {
		return err
	}
	if err := validateEngineSpecs(engine.EngineSize, engine.Cylinders, engine.Horsepower); err != nil {
		return err
	}
	return nil
}

// validateEngineSpecs validates engine specifications
func validateEngineSpecs(engineSize float64, cylinders, horsepower int) error {
	if engineSize <= 0 || engineSize > 12.0 {
		return errors.New("engine size must be between 0.1 and 12.0 liters")
	}
	if cylinders <= 0 || cylinders > 16 {
		return errors.New("number of cylinders must be between 1 and 16")
	}
	if horsepower < 0 || horsepower > 2000 {
		return errors.New("horsepower must be between 0 and 2000")
	}
	return nil
}

// validateLocation validates car location information
func validateLocation(city, state, country string) error {
	if len(city) < 2 {
		return errors.New("city must be at least 2 characters long")
	}
	if len(state) < 2 {
		return errors.New("state must be at least 2 characters long")
	}
	if len(country) < 2 {
		return errors.New("country must be at least 2 characters long")
	}
	return nil
}

// validatePrice validates the price struct and all its fields
func validatePrice(price float64) error {
	if price <= 0 {
		return errors.New("rental price must be greater than 0")
	}
	return nil
}

// validateStatus ensures the status is valid
func validateStatus(status string) error {
	validStatuses := []string{"active", "maintenance", "inactive"}
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return nil
		}
	}
	return errors.New("status must be one of: active, maintenance, inactive")
}

// validateMileage validates car mileage
func validateMileage(mileage int) error {
	if mileage < 0 || mileage > 1000000 {
		return errors.New("mileage must be between 0 and 1,000,000")
	}
	return nil
}
