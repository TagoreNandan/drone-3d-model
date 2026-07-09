package missions

import (
	"context"

	"gcs/internal/db"
	"gcs/internal/pgutil"
)

// Service holds mission business logic on top of the sqlc-generated queries.
type Service struct {
	queries *db.Queries
}

func NewService(queries *db.Queries) *Service {
	return &Service{queries: queries}
}

func (s *Service) ListByVehicle(ctx context.Context, vehicleID string) ([]MissionResponse, error) {
	uid, err := pgutil.ParseUUID(vehicleID)
	if err != nil {
		return nil, err
	}
	rows, err := s.queries.ListMissionsByVehicle(ctx, uid)
	if err != nil {
		return nil, err
	}
	out := make([]MissionResponse, 0, len(rows))
	for _, row := range rows {
		out = append(out, toResponse(row))
	}
	return out, nil
}

func (s *Service) Create(ctx context.Context, req CreateMissionRequest) (*MissionResponse, error) {
	vehicleID, err := pgutil.ParseUUID(req.VehicleID)
	if err != nil {
		return nil, err
	}
	row, err := s.queries.CreateMission(ctx, db.CreateMissionParams{
		VehicleID:       vehicleID,
		Name:            req.Name,
		PlanJson:        req.PlanJSON,
		GeofenceJson:    req.GeofenceJSON,
		RallyPointsJson: req.RallyPointsJSON,
	})
	if err != nil {
		return nil, err
	}
	resp := toResponse(row)
	return &resp, nil
}

func (s *Service) MarkUploaded(ctx context.Context, id string) error {
	uid, err := pgutil.ParseUUID(id)
	if err != nil {
		return err
	}
	return s.queries.MarkMissionUploaded(ctx, uid)
}

func toResponse(m *db.Mission) MissionResponse {
	return MissionResponse{
		ID:              pgutil.UUIDString(m.ID),
		VehicleID:       pgutil.UUIDString(m.VehicleID),
		Name:            m.Name,
		PlanJSON:        m.PlanJson,
		GeofenceJSON:    m.GeofenceJson,
		RallyPointsJSON: m.RallyPointsJson,
		CreatedAt:       m.CreatedAt.Time,
		UploadedAt:      pgutil.TimeOrNil(m.UploadedAt),
	}
}
