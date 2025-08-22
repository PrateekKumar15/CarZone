// Package models contains the data structures and validation logic for the CarZone application.
// This file specifically defines the Engine entity and its related operations.
package models

import (
	"errors"

	"github.com/google/uuid"
)

// Engine represents an engine entity in the CarZone system.
// It contains all the technical specifications and characteristics of a vehicle engine.
// Each engine has a unique identifier and is associated with cars through foreign key relationships.
type Engine struct {
	ID            uuid.UUID `json:"engine_id"`       // Unique identifier for the engine (UUID format)
	Displacement  int64     `json:"displacement"`    // Engine displacement in cubic centimeters (cc)
	NoOfCylinders int64     `json:"no_of_cylinders"` // Number of cylinders in the engine
	CarRange      int64     `json:"car_range"`       // Maximum range of vehicles with this engine (in kilometers)
}

// EngineRequest represents the data structure for creating or updating an engine.
// It contains all the necessary fields for engine operations but excludes the system-generated
// ID field which is managed internally by the application.
// This struct is used for API requests where clients provide engine specifications.
type EngineRequest struct {
	Displacement  int64 `json:"displacement"`    // Engine displacement in cubic centimeters (cc)
	NoOfCylinders int64 `json:"no_of_cylinders"` // Number of cylinders in the engine
	CarRange      int64 `json:"car_range"`       // Maximum range of vehicles with this engine (in kilometers)
}

// ValidateEngineRequest performs comprehensive validation on an EngineRequest.
// It validates all fields including displacement, number of cylinders, and car range
// to ensure data integrity and business rule compliance.
//
// Parameters:
//   - engineRequest: The EngineRequest struct to validate
//
// Returns:
//   - error: An error if any validation fails, nil if all validations pass
//
// Validation Rules:
//   - Displacement must be a positive integer (> 0)
//   - Number of cylinders must be a positive integer (> 0)
//   - Car range must be a positive integer (> 0)
func ValidateEngineRequest(engineRequest EngineRequest) error {
	if err := validateDisplacement(engineRequest.Displacement); err != nil {
		return err
	}
	if err := validateNoOfCylinders(engineRequest.NoOfCylinders); err != nil {
		return err
	}
	if err := validateCarRange(engineRequest.CarRange); err != nil {
		return err
	}
	return nil
}

// validateDisplacement validates the engine displacement value.
// Engine displacement must be a positive value as it represents the engine's cubic capacity.
//
// Parameters:
//   - displacement: The displacement value in cubic centimeters to validate
//
// Returns:
//   - error: An error if displacement is invalid, nil if valid
//
// Business Rules:
//   - Must be greater than 0 (positive integer)
//   - Represents engine displacement in cubic centimeters (cc)
func validateDisplacement(displacement int64) error {
	if displacement <= 0 {
		return errors.New("displacement must be a positive integer representing engine capacity in cc")
	}
	return nil
}

// validateNoOfCylinders validates the number of cylinders in the engine.
// The number of cylinders must be a positive integer as engines cannot have zero or negative cylinders.
//
// Parameters:
//   - noOfCylinders: The number of cylinders to validate
//
// Returns:
//   - error: An error if the number of cylinders is invalid, nil if valid
//
// Business Rules:
//   - Must be greater than 0 (positive integer)
//   - Represents the actual count of engine cylinders
//   - Common values: 3, 4, 6, 8, 10, 12, etc.
func validateNoOfCylinders(noOfCylinders int64) error {
	if noOfCylinders <= 0 {
		return errors.New("number of cylinders must be a positive integer")
	}
	return nil
}

// validateCarRange validates the car range value.
// Car range represents the maximum distance a vehicle can travel and must be positive.
//
// Parameters:
//   - carRange: The car range value in kilometers to validate
//
// Returns:
//   - error: An error if car range is invalid, nil if valid
//
// Business Rules:
//   - Must be greater than 0 (positive integer)
//   - Represents maximum travel range in kilometers
//   - Used for fuel efficiency and performance metrics
func validateCarRange(carRange int64) error {
	if carRange <= 0 {
		return errors.New("car range must be a positive integer representing maximum travel distance in km")
	}
	return nil
}
