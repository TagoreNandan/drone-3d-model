# internal/realtime

A working gorilla/websocket hub — implemented, but **not currently wired into `main.go`**.

- `NewHub()` — constructs a `Hub` with client set, broadcast channel, register/unregister channels
- `(*Hub) Run()` — event loop; call in a goroutine (`go hub.Run()`)
- `(*Hub) ServeWS(w, r)` — upgrades an HTTP request to a WebSocket, registers the connection, and echoes any inbound message to all connected clients via `Broadcast`
- `(*Hub) Broadcast(message []byte)` — queues a message for delivery to every connected client

`upgrader.CheckOrigin` currently always returns `true` — fine for local dev, but tighten this before exposing the endpoint publicly.

## To use this

1. In `main.go`: `hub := realtime.NewHub(); go hub.Run()`
2. Register a route: `e.GET("/ws", func(c echo.Context) error { hub.ServeWS(c.Response(), c.Request()); return nil })`
3. Feed real telemetry into it from `internal/mavlink` (once implemented) via `hub.Broadcast(...)` instead of relying on the current echo-back-to-all-clients behavior.
