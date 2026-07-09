package vehicles

import (
	"context"

	"gcs/internal/db"
	"gcs/internal/pgutil"
)

// Service holds vehicle business logic on top of the sqlc-generated queries.
type Service struct {
	queries *db.Queries
}

func NewService(queries *db.Queries) *Service {
	return &Service{queries: queries}
}

func (s *Service) List(ctx context.Context) ([]VehicleResponse, error) {
	rows, err := s.queries.ListVehicles(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]VehicleResponse, 0, len(rows))
	for _, row := range rows {
		out = append(out, toResponse(row))
	}
	return out, nil
}

func (s *Service) Get(ctx context.Context, id string) (*VehicleResponse, error) {
	uid, err := pgutil.ParseUUID(id)
	if err != nil {
		return nil, err
	}
	row, err := s.queries.GetVehicle(ctx, uid)
	if err != nil {
		return nil, err
	}
	resp := toResponse(row)
	return &resp, nil
}

func (s *Service) Create(ctx context.Context, req CreateVehicleRequest) (*VehicleResponse, error) {
	row, err := s.queries.CreateVehicle(ctx, db.CreateVehicleParams{
		Name:         req.Name,
		Type:         pgutil.NewText(req.Type),
		Firmware:     pgutil.NewText(req.Firmware),
		MavlinkSysid: pgutil.NewInt4(req.MavlinkSysID),
	})
	if err != nil {
		return nil, err
	}
	resp := toResponse(row)
	return &resp, nil
}

func toResponse(v *db.Vehicle) VehicleResponse {
	return VehicleResponse{
		ID:           pgutil.UUIDString(v.ID),
		Name:         v.Name,
		Type:         pgutil.TextOrNil(v.Type),
		Firmware:     pgutil.TextOrNil(v.Firmware),
		MavlinkSysID: pgutil.Int4OrNil(v.MavlinkSysid),
		CreatedAt:    v.CreatedAt.Time,
	}
}
