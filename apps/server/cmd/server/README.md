# cmd/server

Backend entrypoint (`main.go`). Responsibilities, in order:

1. Load `.env` (via `godotenv`)
2. Init zap logger (`LOG_LEVEL=production` → JSON prod logger, otherwise dev logger)
3. Open the pgx pool (`internal/database.InitDB`) — fatals on failure
4. Build the Echo instance: `Logger`/`Recover`/`CORS` middleware, then routes
5. Register routes: `/health`, `/`, `/swagger/*` (Swagger UI), and `/api/*` from each module's `Controller.RegisterRoutes` (see `internal/modules/{vehicles,missions,telemetry,alerts}`)
6. Start listening on `HOST:PORT` (default `0.0.0.0:8080`)

General Swagger API metadata (`@title`, `@version`, `@description`, `@host`, `@BasePath`) lives in the comment block directly above `func main()` — `swag init` reads it from there, so keep it attached to `main`, not moved elsewhere.

Routes not yet registered here: a `/ws` WebSocket endpoint for `internal/realtime.Hub` (implemented but unwired).
