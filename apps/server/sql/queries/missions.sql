-- name: ListMissionsByVehicle :many
SELECT * FROM missions WHERE vehicle_id = $1 ORDER BY created_at DESC;

-- name: CreateMission :one
INSERT INTO missions (vehicle_id, name, plan_json, geofence_json, rally_points_json)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: MarkMissionUploaded :exec
UPDATE missions SET uploaded_at = now() WHERE id = $1;
