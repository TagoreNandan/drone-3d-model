-- name: InsertAlert :exec
INSERT INTO alerts (vehicle_id, severity, text) VALUES ($1, $2, $3);

-- name: ListRecentAlerts :many
SELECT * FROM alerts WHERE vehicle_id = $1 ORDER BY ts DESC LIMIT 50;
