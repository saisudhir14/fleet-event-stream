package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/saisudhir14/fleet-event-stream/internal/metrics"
	"github.com/saisudhir14/fleet-event-stream/internal/models"
	"github.com/saisudhir14/fleet-event-stream/internal/processor"
)

// Handler struct encapsulates dependencies for HTTP handlers
type Handler struct {
	processor *processor.EventProcessor
	metrics   *metrics.Metrics
	logger    *slog.Logger
}

// NewHandler creates a new Handler instance
// parameters: proc *processor.EventProcessor: The event processor instance.
//
//	m *metrics.Metrics: The metrics collector instance.
//
//	logger *slog.Logger: The logger instance.
//
// returns: *Handler: A new Handler instance with dependencies injected.
func NewHandler(proc *processor.EventProcessor, m *metrics.Metrics, logger *slog.Logger) *Handler {
	return &Handler{
		processor: proc,
		metrics:   m,
		logger:    logger,
	}
}

// HealthCheck handles the /health endpoint
// parameters: w http.ResponseWriter: The HTTP response writer.
//
//	r *http.Request: The HTTP request.
//
// returns: none but writes a JSON health status to the response.
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	health := map[string]string{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"service":   "fleet-event-stream",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(health)
}

// ReadyCheck handles the /ready endpoint
// parameters: w http.ResponseWriter: The HTTP response writer.
//
//	r *http.Request: The HTTP request.
//
// returns: none but writes a JSON readiness status to the response.
func (h *Handler) ReadyCheck(w http.ResponseWriter, r *http.Request) {
	ready := map[string]string{
		"status": "ready",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ready)
}

// IngestEvent handles the /api/v1/events endpoint
// Accepts vehicle event data and processes it
// parameters: w http.ResponseWriter: The HTTP response writer.
//
//	r *http.Request: The HTTP request containing event data.
//
// returns: none but writes a JSON response indicating acceptance or error.
func (h *Handler) IngestEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var event models.VehicleEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		h.logger.Error("failed to decode event", "error", err)
		h.metrics.RecordValidation(false)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	start := time.Now()

	if err := h.processor.ProcessEvent(r.Context(), &event); err != nil {
		h.metrics.RecordValidation(false)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Record metrics
	duration := time.Since(start).Seconds()
	h.metrics.RecordValidation(true)
	h.metrics.RecordEventProcessed(event.EventType, event.VehicleID)
	h.metrics.RecordProcessingDuration(event.EventType, duration)
	h.logger.Info("event ingested",
		"event_id", event.EventID,
		"vehicle_id", event.VehicleID,
		"duration_ms", duration*1000,
	)

	response := map[string]string{
		"status":   "accepted",
		"event_id": event.EventID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)

}

// GetStats handles the /api/v1/stats endpoint
// Returns processing statistics
// parameters: w http.ResponseWriter: The HTTP response writer.
//
//	r *http.Request: The HTTP request.
//
// returns: none but writes a JSON response containing statistics.
func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats := map[string]interface{}{
		"total_events_processed": h.processor.GetTotalEventCount(),
		"timestamp":              time.Now().UTC().Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}
