package mavlink

// TODO: MAVSDK-Go session lifecycle per vehicle (per TRD §2.3).
//
// Standard telemetry/action/mission plugins should come from MAVSDK-Go —
// confirm the current import path/version from https://mavsdk.mavlink.io
// before pinning it, as this can change across releases.
//
// The custom TUNNEL payload (see tunnel.go) needs raw MAVLink parsing —
// go get github.com/... a maintained Go MAVLink dialect library — since
// it is not one of MAVSDK's standard plugins.
//
// Session should feed decoded frames into internal/realtime.Hub.broadcast
// (see internal/realtime/websocket.go) and into a batched writer against
// the telemetry_frames hypertable via the generated sqlc queries
// (sql/queries/telemetry.sql).
