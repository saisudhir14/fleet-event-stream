package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics struct holds Prometheus metrics collectors
type Metrics struct {
	EventsProcessed   *prometheus.CounterVec
	EventsValidation  *prometheus.CounterVec
	ProcessingLatency *prometheus.HistogramVec
	ActiveVehicles    prometheus.Gauge
}

// NewMetrics initializes and returns a Metrics instance
// Returns a new Metrics instance with all Prometheus metrics initialized
// Each metric is defined with appropriate labels and help descriptions.
func NewMetrics() *Metrics {
	return &Metrics{
		EventsProcessed: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "fleet_events_processed_total",
				Help: "Total number of vehicle events processed",
			},
			[]string{"event_type", "vehicle_id"},
		),

		EventsValidation: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "fleet_events_validation_total",
				Help: "Total number of validation results",
			},

			[]string{"status"},
		),

		ProcessingLatency: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "fleet_event_processing_duration_seconds",
				Help:    "Event processing latency in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"event_type"},
		),

		ActiveVehicles: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "fleet_active_vehicles",
				Help: "Number of active vehicles currently tracked",
			},
		),
	}
}

// RecordEventProcessed increments the EventsProcessed counter
// for a given event type and vehicle ID.
// parameters: eventType string: The type of the event processed.
//
//	vehicleID string: The ID of the vehicle associated with the event.
//
// returns: none
func (m *Metrics) RecordEventProcessed(eventType, vehicleID string) {
	m.EventsProcessed.WithLabelValues(eventType, vehicleID).Inc()
}

// RecordValidation increments the EventsValidation counter
// based on the success or failure of event validation.
// parameters: success bool: Indicates whether the validation was successful.
// returns: none
func (m *Metrics) RecordValidation(success bool) {
	status := "success"
	if !success {
		status = "failed"
	}
	m.EventsValidation.WithLabelValues(status).Inc()

}

// RecordProcessingDuration observes the processing duration
// for a given event type.
// parameters: eventType string: The type of the event processed.
//
//	duration float64: The processing duration in seconds.
//
// returns: none
func (m *Metrics) RecordProcessingDuration(eventType string, duration float64) {
	m.ProcessingLatency.WithLabelValues(eventType).Observe(duration)
}

// UpdateActiveVehicles sets the ActiveVehicles gauge
// to the current count of active vehicles.
// parameters: count float64: The current number of active vehicles.
// returns: none
func (m *Metrics) UpdateActiveVehicles(count float64) {
	m.ActiveVehicles.Set(count)

}
