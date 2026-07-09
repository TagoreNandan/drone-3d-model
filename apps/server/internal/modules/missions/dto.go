package missions

import (
	"encoding/json"
	"time"
)

// MissionResponse is returned from all mission endpoints.
type MissionResponse struct {
	ID              string          `json:"id"`
	VehicleID       string          `json:"vehicle_id"`
	Name            string          `json:"name"`
	PlanJSON        json.RawMessage `json:"plan_json" swaggertype:"object"`
	GeofenceJSON    json.RawMessage `json:"geofence_json,omitempty" swaggertype:"object"`
	RallyPointsJSON json.RawMessage `json:"rally_points_json,omitempty" swaggertype:"object"`
	CreatedAt       time.Time       `json:"created_at"`
	UploadedAt      *time.Time      `json:"uploaded_at,omitempty"`
}

// CreateMissionRequest is the body for POST /api/missions.
type CreateMissionRequest struct {
	VehicleID       string          `json:"vehicle_id"`
	Name            string          `json:"name"`
	PlanJSON        json.RawMessage `json:"plan_json" swaggertype:"object"`
	GeofenceJSON    json.RawMessage `json:"geofence_json" swaggertype:"object"`
	RallyPointsJSON json.RawMessage `json:"rally_points_json" swaggertype:"object"`
}
