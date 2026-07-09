package alerts

import (
	"context"

	"gcs/internal/db"
	"gcs/internal/pgutil"
)

// Service holds alert business logic on top of the sqlc-generated queries.
type Service struct {
	queries *db.Queries
}

func NewService(queries *db.Queries) *Service {
	return &Service{queries: queries}
}

func (s *Service) Create(ctx context.Context, req CreateAlertRequest) error {
	vehicleID, err := pgutil.ParseUUID(req.VehicleID)
	if err != nil {
		return err
	}
	return s.queries.InsertAlert(ctx, db.InsertAlertParams{
		VehicleID: vehicleID,
		Severity:  pgutil.NewInt2(req.Severity),
		Text:      pgutil.NewText(req.Text),
	})
}

func (s *Service) ListRecent(ctx context.Context, vehicleID string) ([]AlertResponse, error) {
	uid, err := pgutil.ParseUUID(vehicleID)
	if err != nil {
		return nil, err
	}
	rows, err := s.queries.ListRecentAlerts(ctx, uid)
	if err != nil {
		return nil, err
	}
	out := make([]AlertResponse, 0, len(rows))
	for _, row := range rows {
		out = append(out, toResponse(row))
	}
	return out, nil
}

func toResponse(a *db.Alert) AlertResponse {
	return AlertResponse{
		ID:        a.ID,
		VehicleID: pgutil.UUIDString(a.VehicleID),
		Ts:        a.Ts.Time,
		Severity:  pgutil.Int2OrNil(a.Severity),
		Text:      pgutil.TextOrNil(a.Text),
	}
}
