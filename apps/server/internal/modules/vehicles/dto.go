package vehicles

import "time"

// VehicleResponse is returned from all vehicle endpoints.
type VehicleResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Type         *string   `json:"type,omitempty"`
	Firmware     *string   `json:"firmware,omitempty"`
	MavlinkSysID *int32    `json:"mavlink_sysid,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// CreateVehicleRequest is the body for POST /api/vehicles.
type CreateVehicleRequest struct {
	Name         string  `json:"name"`
	Type         *string `json:"type"`
	Firmware     *string `json:"firmware"`
	MavlinkSysID *int32  `json:"mavlink_sysid"`
}
