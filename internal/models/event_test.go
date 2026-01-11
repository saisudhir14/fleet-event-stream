package models

import (
	"testing"
	"time"
)

func TestVehicleEvent_Validate(t *testing.T) {
	tests := []struct {
		name    string
		event   VehicleEvent
		wantErr bool
	}{

		{
			name: "valid event",
			event: VehicleEvent{
				EventID:   "evt-001",
				VehicleID: "vehicle-123",
				Latitude:  33.4484,
				Longitude: -112.0740,
				Speed:     65.5,
				Heading:   180.0,
				Timestamp: time.Now(),
				EventType: EventTypePosition,
			},
			wantErr: false,
		},

		{
			name: "missing event_id",
			event: VehicleEvent{
				VehicleID: "vehicle-123",
				Latitude:  33.4484,
				Longitude: -112.0740,
				Speed:     65.5,
				Timestamp: time.Now(),
				EventType: EventTypePosition,
			},
			wantErr: true,
		},

		{
			name: "invalid latitude",
			event: VehicleEvent{
				EventID:   "evt-001",
				VehicleID: "vehicle-123",
				Latitude:  91.0,
				Longitude: -112.0740,
				Speed:     65.5,
				Timestamp: time.Now(),
				EventType: EventTypePosition,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.event.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)

			}

		})
	}
}
