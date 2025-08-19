// Package models contains the data structures and validation logic for the CarZone application
package models

import (
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// Car represents a vehicle entity in the CarZone system
// It contains all the necessary information about a car including its specifications,
// pricing, and metadata for tracking creation and updates
type Car struct {
	ID        uuid.UUID `json:"id"`         // Unique identifier for the car
	Name      string    `json:"name"`       // Display name/model of the car
	Year      int       `json:"year"`       // Manufacturing year of the car
	Brand     string    `json:"brand"`      // Manufacturer brand name
	FuelType  string    `json:"fuel_type"`  // Type of fuel (Petrol, Diesel, Electric, Hybrid)
	Engine    Engine    `json:"engine"`     // Engine specifications
	Price     float64   `json:"price"`      // Price of the car in the system's currency
	CreatedAt time.Time `json:"created_at"` // Timestamp when the car record was created
	UpdatedAt time.Time `json:"updated_at"` // Timestamp when the car record was last updated
}

// CarRequest represents the data structure for creating or updating a car
// It contains all the necessary fields for car creation but excludes system-generated
// fields like ID, CreatedAt, and UpdatedAt which are managed internally
type CarRequest struct {
	Name     string  `json:"name"`      // Display name/model of the car
	Year     int     `json:"year"`      // Manufacturing year of the car
	Brand    string  `json:"brand"`     // Manufacturer brand name
	FuelType string  `json:"fuel_type"` // Type of fuel (Petrol, Diesel, Electric, Hybrid)
	Engine   Engine  `json:"engine"`    // Engine specifications
	Price    float64 `json:"price"`     // Price of the car in the system's currency
}

// ValidateRequest performs comprehensive validation on a CarRequest
// It validates all fields including name, year, brand, fuel type, engine, and price
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
	if err := validateFuelType(carRequest.FuelType); err != nil {
		return err
	}
	if err := validateEngine(carRequest.Engine); err != nil {
		return err
	}
	if err := validatePrice(carRequest.Price); err != nil {
		return err
	}
	return nil
}

// validateName checks if the car name meets the minimum length requirement
// Car names must be at least 3 characters long to be considered valid
// Returns an error if the name is too short, nil otherwise
func validateName(name string) error {
	if len(name) < 3 {
		return errors.New("name must be at least 3 characters long")
	}
	return nil
}

// validateYear validates the manufacturing year of the car
// Accepts the year as a string and performs the following checks:
// - Ensures the year is not empty
// - Verifies the year is a valid integer
// - Checks that the year is between 1886 (first automobile) and the current year
// Returns an error if any validation fails, nil otherwise
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

// validateBrand checks if the car brand name meets the minimum length requirement
// Brand names must be at least 2 characters long to be considered valid
// Returns an error if the brand name is too short, nil otherwise
func validateBrand(brand string) error {
	if len(brand) < 2 {
		return errors.New("brand must be at least 2 characters long")
	}
	return nil
}

// validateFuelType ensures the fuel type is one of the accepted values
// Valid fuel types are: Petrol, Diesel, Electric, Hybrid
// The validation is case-sensitive and exact match is required
// Returns an error if the fuel type is not in the valid list, nil otherwise
func validateFuelType(fuelType string) error {
	validFuelTypes := []string{"Petrol", "Diesel", "Electric", "Hybrid"}
	for _, validType := range validFuelTypes {
		if fuelType == validType {
			return nil
		}
	}
	return errors.New("fuel type must be one of: Petrol, Diesel, Electric, Hybrid")
}

// validateEngine performs validation on the engine specifications
// Checks the following engine properties:
// - Displacement: Must be a positive number (greater than 0)
// - Number of cylinders: Must be a positive integer (greater than 0)
// - Car range: Must be a positive number (greater than 0)
// Returns an error if any engine property is invalid, nil otherwise
func validateEngine(engine Engine) error {
	if engine.Displacement <= 0 {
		return errors.New("displacement must be a positive number")
	}
	if engine.NoOfCylinders <= 0 {
		return errors.New("number of cylinders must be a positive number")
	}
	if engine.CarRange <= 0 {
		return errors.New("car range must be a positive number")
	}
	return nil
}

// validatePrice ensures the car price is a positive value
// Price must be greater than 0 to be considered valid
// Returns an error if the price is zero or negative, nil otherwise
func validatePrice(price float64) error {
	if price <= 0 {
		return errors.New("price must be a positive number")
	}
	return nil
}
