package telemetry

import (
	"context"
	"time"

	"gcs/internal/db"
	"gcs/internal/pgutil"
)

// Service holds telemetry business logic on top of the sqlc-generated queries.
type Service struct {
	queries *db.Queries
}

func NewService(queries *db.Queries) *Service {
	return &Service{queries: queries}
}

func (s *Service) Ingest(ctx context.Context, req IngestTelemetryRequest) error {
	vehicleID, err := pgutil.ParseUUID(req.VehicleID)
	if err != nil {
		return err
	}
	return s.queries.InsertTelemetryFrame(ctx, db.InsertTelemetryFrameParams{
		VehicleID:   vehicleID,
		Ts:          pgutil.NewTimestamptz(req.Ts),
		Lat:         pgutil.NewFloat8(req.Lat),
		Lon:         pgutil.NewFloat8(req.Lon),
		Alt:         pgutil.NewFloat8(req.Alt),
		Heading:     pgutil.NewFloat8(req.Heading),
		Groundspeed: pgutil.NewFloat8(req.Groundspeed),
		BatteryPct:  pgutil.NewInt2(req.BatteryPct),
		Voltage:     pgutil.NewFloat8(req.Voltage),
		Current:     pgutil.NewFloat8(req.Current),
		FlightMode:  pgutil.NewText(req.FlightMode),
		Armed:       pgutil.NewBool(req.Armed),
		CpuPct:      pgutil.NewFloat4(req.CpuPct),
		RamPct:      pgutil.NewFloat4(req.RamPct),
		DiskPct:     pgutil.NewFloat4(req.DiskPct),
		CpuTemp:     pgutil.NewFloat4(req.CpuTemp),
		NodeStatus:  req.NodeStatus,
	})
}

func (s *Service) GetRange(ctx context.Context, vehicleID string, from, to time.Time) ([]TelemetryFrameResponse, error) {
	uid, err := pgutil.ParseUUID(vehicleID)
	if err != nil {
		return nil, err
	}
	rows, err := s.queries.GetTelemetryRange(ctx, db.GetTelemetryRangeParams{
		VehicleID: uid,
		Ts:        pgutil.NewTimestamptz(from),
		Ts_2:      pgutil.NewTimestamptz(to),
	})
	if err != nil {
		return nil, err
	}
	out := make([]TelemetryFrameResponse, 0, len(rows))
	for _, row := range rows {
		out = append(out, toResponse(row))
	}
	return out, nil
}

func toResponse(f *db.TelemetryFrame) TelemetryFrameResponse {
	return TelemetryFrameResponse{
		VehicleID:   pgutil.UUIDString(f.VehicleID),
		Ts:          f.Ts.Time,
		Lat:         pgutil.Float8OrNil(f.Lat),
		Lon:         pgutil.Float8OrNil(f.Lon),
		Alt:         pgutil.Float8OrNil(f.Alt),
		Heading:     pgutil.Float8OrNil(f.Heading),
		Groundspeed: pgutil.Float8OrNil(f.Groundspeed),
		BatteryPct:  pgutil.Int2OrNil(f.BatteryPct),
		Voltage:     pgutil.Float8OrNil(f.Voltage),
		Current:     pgutil.Float8OrNil(f.Current),
		FlightMode:  pgutil.TextOrNil(f.FlightMode),
		Armed:       pgutil.BoolOrNil(f.Armed),
		CpuPct:      pgutil.Float4OrNil(f.CpuPct),
		RamPct:      pgutil.Float4OrNil(f.RamPct),
		DiskPct:     pgutil.Float4OrNil(f.DiskPct),
		CpuTemp:     pgutil.Float4OrNil(f.CpuTemp),
		NodeStatus:  f.NodeStatus,
	}
}
