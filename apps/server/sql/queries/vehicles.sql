-- name: ListVehicles :many
SELECT * FROM vehicles ORDER BY created_at DESC;

-- name: GetVehicle :one
SELECT * FROM vehicles WHERE id = $1;

-- name: CreateVehicle :one
INSERT INTO vehicles (name, type, firmware, mavlink_sysid)
VALUES ($1, $2, $3, $4)
RETURNING *;
