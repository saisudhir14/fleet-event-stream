package processor

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/saisudhir14/fleet-event-stream/internal/models"
)

// EventProcessor handles the processing of vehicle events. It maintains
// per-vehicle event counts and provides thread-safe access to statistics.
type EventProcessor struct {
	mu         sync.RWMutex
	eventCount map[string]int64
	logger     *slog.Logger
}

// NewEventProcessor creates a new EventProcessor instance
// parameters: logger *slog.Logger: The logger instance.
//
// returns: *EventProcessor: A new EventProcessor instance.
func NewEventProcessor(logger *slog.Logger) *EventProcessor {
	return &EventProcessor{
		eventCount: make(map[string]int64),
		logger:     logger,
	}
}

// ProcessEvent processes a VehicleEvent based on its type
// parameters: ctx context.Context: The context for the processing operation.
//
//	event *models.VehicleEvent: The vehicle event to be processed.
//
// returns: error: An error if processing fails, nil otherwise.
func (p *EventProcessor) ProcessEvent(ctx context.Context, event *models.VehicleEvent) error {

	if err := event.Validate(); err != nil {
		p.logger.Error("event validation failed",
			"error", err,
			"event_id", event.EventID,
			"vehicle_id", event.VehicleID,
		)
		return fmt.Errorf("validation error: %w", err)
	}

	switch event.EventType {
	case models.EventTypeSpeedAlert:
		p.handleSpeedAlert(event)
	case models.EventTypeGeofence:
		p.handleGeofenceEvent(event)
	case models.EventTypePosition:
		p.handlePositionUpdate(event)
	default:
		p.handleGenericEvent(event)
	}
	p.incrementEventCount(event.VehicleID)
	p.logger.Info("event processed successfully",
		"event_id", event.EventID,
		"vehicle_id", event.VehicleID,
		"event_type", event.EventType,
	)
	return nil
}

// handleSpeedAlert processes speed alert events
// parameters: event *models.VehicleEvent: The vehicle event to be processed.
// returns: none
func (p *EventProcessor) handleSpeedAlert(event *models.VehicleEvent) {
	p.logger.Warn("speed alert detected",
		"vehicle_id", event.VehicleID,
		"speed", event.Speed,
		"timestamp", event.Timestamp,
	)
}

// handleGeofenceEvent processes geofence events
// parameters: event *models.VehicleEvent: The vehicle event to be processed.
// returns: none
func (p *EventProcessor) handleGeofenceEvent(event *models.VehicleEvent) {
	p.logger.Info("geofence event",
		"vehicle_id", event.VehicleID,
		"location", fmt.Sprintf("%.6f,%.6f", event.Latitude, event.Longitude),
	)
}

// handlePositionUpdate processes position update events
// parameters: event *models.VehicleEvent: The vehicle event to be processed.
// returns: none
func (p *EventProcessor) handlePositionUpdate(event *models.VehicleEvent) {
	p.logger.Debug("position updated",
		"vehicle_id", event.VehicleID,
		"location", fmt.Sprintf("%.6f,%.6f", event.Latitude, event.Longitude),
	)
}

// handleGenericEvent processes events of unknown types
// parameters: event *models.VehicleEvent: The vehicle event to be processed.
// returns: none
func (p *EventProcessor) handleGenericEvent(event *models.VehicleEvent) {
	p.logger.Info("generic event processed",
		"vehicle_id", event.VehicleID,
		"event_type", event.EventType,
	)
}

// incrementEventCount increments the event count for a given vehicle ID
// parameters: vehicleID string: The ID of the vehicle.
// returns: none
func (p *EventProcessor) incrementEventCount(vehicleID string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.eventCount[vehicleID]++
}

// GetEventCount retrieves the event count for a given vehicle ID
// parameters: vehicleID string: The ID of the vehicle.
//
// returns: int64: The number of events processed for the vehicle.
func (p *EventProcessor) GetEventCount(vehicleID string) int64 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.eventCount[vehicleID]
}

// GetTotalEventCount retrieves the total number of events processed across all vehicles
// parameters: none
//
// returns: int64: The total number of events processed.
func (p *EventProcessor) GetTotalEventCount() int64 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	var total int64

	for _, count := range p.eventCount {
		total += count
	}
	return total
}
