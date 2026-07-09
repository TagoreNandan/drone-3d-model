# gcs — Ground Control Station

A drone/UAV ground control station: MAVLink telemetry ingestion into TimescaleDB, a live map/mission web UI, and a WebSocket channel for realtime vehicle data. Originally scaffolded with [Better Fullstack](https://github.com/Marve10s/Better-Fullstack); the domain-specific pieces (schema, MAVLink handling, mapping UI) are being layered on top.

See [`apps/server/README.md`](apps/server/README.md) and [`apps/web/README.md`](apps/web/README.md) for stack and implementation detail per app.

## Stack

- **Frontend**: React 19 + Vite + react-router, Tailwind 4, shadcn/radix UI, Zustand — [`apps/web`](apps/web/README.md)
- **Backend**: Go + Echo, pgx, sqlc, gorilla/websocket, zap — [`apps/server`](apps/server/README.md)
- **Database**: TimescaleDB (Postgres 16 + hypertables), via Docker
- **Tooling**: pnpm workspaces + Turborepo for the JS side; plain Go modules for the backend

## Project Structure

```text
gcs/
├── apps/
│   ├── web/         # React/Vite frontend — see apps/web/README.md
│   └── server/      # Go/Echo backend — see apps/server/README.md
├── packages/
│   ├── config/      # Shared tsconfig base
│   ├── db/          # Placeholder — actual DB stack lives in docker-compose.yml + apps/server/sql
│   └── env/         # Typed env helpers (@t3-oss/env-core)
├── docker-compose.yml  # TimescaleDB container
├── run.sh              # One-command local launch (db + backend + frontend)
└── package.json        # Root scripts
```

## Local Development

The fastest path is the launch script at the repo root:

```sh
./run.sh
```

It starts TimescaleDB via Docker, waits for it to be ready, runs pending database migrations, regenerates Swagger docs, then starts the Go backend (`:8080`) and Vite frontend (`:5173`), and tears everything down cleanly on Ctrl+C.

### Manual steps (what `run.sh` automates)

1. **Database** — TimescaleDB via Docker, published on host port **5433** (not 5432 — see note below):

   ```sh
   docker compose up -d
   ```

2. **Migrations** — versioned with [golang-migrate](https://github.com/golang-migrate/migrate) (see `apps/server/sql/migrations/README.md`):

   ```sh
   go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
   migrate -path apps/server/sql/migrations -database "postgres://gcs:gcs@localhost:5433/gcs?sslmode=disable" up
   ```

3. **Backend** (Go 1.22+ required):

   ```sh
   cd apps/server
   go mod tidy
   go run cmd/server/main.go
   ```

   Serves on `http://localhost:8080`. Health check: `GET /health`. Swagger UI: `GET /swagger/index.html`.

4. **Frontend**:

   ```sh
   pnpm install
   pnpm dev:web
   ```

   Serves on `http://localhost:5173`, and shows a live connected/disconnected indicator against the backend's `/health` endpoint.

### Port note

`docker-compose.yml` publishes TimescaleDB on **5433**, not the default 5432 — this machine already runs a native Postgres 16 service on 5432. `apps/server/.env` / `.env.example` are set to match (`postgres://gcs:gcs@localhost:5433/gcs?sslmode=disable`). If you're on a machine with nothing on 5432, you can remap back to `5432:5432` and drop the `5433` overrides.

## Root Scripts

- `dev` / `dev:web` — start the frontend (`turbo -F web dev`)
- `dev:server` — `cd apps/server && go run cmd/server/main.go`
- `setup:server` — `cd apps/server && go mod tidy`
- `check:server` / `test:server` — backend compile check / `go test ./...`
- `build` / `check-types` — Turborepo tasks across the JS workspace
- `db:start` / `db:watch` / `db:stop` / `db:down` — **currently broken**: these reference a `@gcs/db` package script that doesn't exist (`packages/db` has no `package.json`). Use `docker compose up -d` / `docker compose down` directly, or `./run.sh`.

## Implementation Status

This is early-stage. Roughly:

- ✅ Working: Postgres/Timescale schema (versioned via golang-migrate), Echo server + health check, vehicles/missions/telemetry/alerts REST endpoints (controller/service/dto per module in `apps/server/internal/modules`), Swagger docs, websocket hub (implemented but not wired to a route), Vite/React shell with routing, theming, and a live backend-health widget.
- ❌ Not started: MAVLink ingestion (`apps/server/internal/mavlink` is stub/TODO), map/mission UI (`apps/web` has no Fly View or 3D View yet — see `apps/web/MAPPING_ENGINES.md` for the plan), CI, auth.

## Compatibility Notes (from the generator)

- Cross-ecosystem graph projects share a plain HTTP boundary — no tRPC or shared API client across the Go/TypeScript language boundary; the frontend just calls the backend's base URL and health endpoint.
