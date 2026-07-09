package telemetry

import (
	"encoding/json"
	"time"
)

// TelemetryFrameResponse is returned from telemetry range queries.
type TelemetryFrameResponse struct {
	VehicleID   string          `json:"vehicle_id"`
	Ts          time.Time       `json:"ts"`
	Lat         *float64        `json:"lat,omitempty"`
	Lon         *float64        `json:"lon,omitempty"`
	Alt         *float64        `json:"alt,omitempty"`
	Heading     *float64        `json:"heading,omitempty"`
	Groundspeed *float64        `json:"groundspeed,omitempty"`
	BatteryPct  *int16          `json:"battery_pct,omitempty"`
	Voltage     *float64        `json:"voltage,omitempty"`
	Current     *float64        `json:"current,omitempty"`
	FlightMode  *string         `json:"flight_mode,omitempty"`
	Armed       *bool           `json:"armed,omitempty"`
	CpuPct      *float32        `json:"cpu_pct,omitempty"`
	RamPct      *float32        `json:"ram_pct,omitempty"`
	DiskPct     *float32        `json:"disk_pct,omitempty"`
	CpuTemp     *float32        `json:"cpu_temp,omitempty"`
	NodeStatus  json.RawMessage `json:"node_status,omitempty" swaggertype:"object"`
}

// IngestTelemetryRequest is the body for POST /api/telemetry.
type IngestTelemetryRequest struct {
	VehicleID   string          `json:"vehicle_id"`
	Ts          time.Time       `json:"ts"`
	Lat         *float64        `json:"lat"`
	Lon         *float64        `json:"lon"`
	Alt         *float64        `json:"alt"`
	Heading     *float64        `json:"heading"`
	Groundspeed *float64        `json:"groundspeed"`
	BatteryPct  *int16          `json:"battery_pct"`
	Voltage     *float64        `json:"voltage"`
	Current     *float64        `json:"current"`
	FlightMode  *string         `json:"flight_mode"`
	Armed       *bool           `json:"armed"`
	CpuPct      *float32        `json:"cpu_pct"`
	RamPct      *float32        `json:"ram_pct"`
	DiskPct     *float32        `json:"disk_pct"`
	CpuTemp     *float32        `json:"cpu_temp"`
	NodeStatus  json.RawMessage `json:"node_status" swaggertype:"object"`
}
