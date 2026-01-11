package models

import (
	"errors"
	"time"
)

// VehicleEvent represents a vehicle event in the fleet management system
type VehicleEvent struct {
	EventID   string    `json:"event_id"`
	VehicleID string    `json:"vehicle_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Speed     float64   `json:"speed"`
	Heading   float64   `json:"heading"`
	Timestamp time.Time `json:"timestamp"`
	EventType string    `json:"event_type"`
	DriverID  string    `json:"driver_id,omitempty"`
}

// Validate checks if the VehicleEvent has all required fields and valid values
// parameters: none
//
// returns: error: An error if validation fails, nil otherwise.
func (e *VehicleEvent) Validate() error {
	if e.EventID == "" {
		return errors.New("event_id is required")
	}

	if e.VehicleID == "" {
		return errors.New("vehicle_id is required")
	}

	if e.Latitude < -90 || e.Latitude > 90 {
		return errors.New("latitude must be between -90 and 90")
	}

	if e.Longitude < -180 || e.Longitude > 180 {
		return errors.New("longitude must be between -180 and 180")
	}

	if e.Speed < 0 {
		return errors.New("speed cannot be negative")
	}

	if e.EventType == "" {
		return errors.New("event_type is required")
	}

	if e.Timestamp.IsZero() {
		return errors.New("timestamp is required")
	}
	return nil
}

// Constants for event types
const (
	EventTypePosition    = "position"
	EventTypeSpeedAlert  = "speed_alert"
	EventTypeGeofence    = "geofence"
	EventTypeEngineStart = "engine_start"
	EventTypeEngineStop  = "engine_stop"
)
