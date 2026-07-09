-- name: InsertTelemetryFrame :exec
INSERT INTO telemetry_frames (
    vehicle_id, ts, lat, lon, alt, heading, groundspeed,
    battery_pct, voltage, current, flight_mode, armed,
    cpu_pct, ram_pct, disk_pct, cpu_temp, node_status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
);

-- name: GetTelemetryRange :many
SELECT * FROM telemetry_frames
WHERE vehicle_id = $1 AND ts BETWEEN $2 AND $3
ORDER BY ts ASC;
