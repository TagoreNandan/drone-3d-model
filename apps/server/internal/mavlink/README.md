# internal/mavlink — STUBS ONLY, not implemented

This package is where MAVLink telemetry ingestion is meant to live. Neither file has a working implementation yet.

## `session.go`

Empty except for a comment describing the intended design (referencing "TRD §2.3", a Technical Requirements Document not present in this repo):

- Use [MAVSDK-Go](https://mavsdk.mavlink.io) for standard telemetry/action/mission plugins per vehicle session — confirm the current import path/version before pinning, since it's not yet in `go.mod`.
- The custom TUNNEL payload (see `tunnel.go`) isn't one of MAVSDK's standard plugins — needs a separate raw MAVLink dialect parsing library.
- Decoded frames should feed into `internal/realtime.Hub.Broadcast` (for live WebSocket delivery) and into a batched writer against the `telemetry_frames` hypertable via sqlc-generated queries from `sql/queries/telemetry.sql`.

## `tunnel.go`

Defines `TunnelPayload` (cpu_pct, ram_pct, disk_pct, cpu_temp_c, node_mask — mirrors a custom MAVLink `TUNNEL` message, payload_type `0x8100`) and `WatchedNodeNames` (`obstacle_avoidance`, `path_planner`, `camera_driver`, `lidar_proc`).

`DecodeTunnelPayload(raw []byte)` is a literal `panic("not implemented")`. It needs to unpack a 17-byte little-endian `<ffffB` struct (per the comment, matching an existing Python prototype `mavlink_scraper.py` and a QGC `CustomInstrumentWidget` C++ fork — neither of which exist in this repo) using `encoding/binary`.

**To implement this package**: add a MAVLink Go library + MAVSDK-Go to `go.mod`, implement `DecodeTunnelPayload` with `encoding/binary.LittleEndian`, and build out session lifecycle management in `session.go` that ties into `internal/realtime` and the sqlc-generated telemetry queries.
