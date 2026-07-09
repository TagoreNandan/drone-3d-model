package alerts

import "time"

// AlertResponse is returned from all alert endpoints.
type AlertResponse struct {
	ID        int64     `json:"id"`
	VehicleID string    `json:"vehicle_id"`
	Ts        time.Time `json:"ts"`
	Severity  *int16    `json:"severity,omitempty"`
	Text      *string   `json:"text,omitempty"`
}

// CreateAlertRequest is the body for POST /api/alerts.
type CreateAlertRequest struct {
	VehicleID string  `json:"vehicle_id"`
	Severity  *int16  `json:"severity"`
	Text      *string `json:"text"`
}
