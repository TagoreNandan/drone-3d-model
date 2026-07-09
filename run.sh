#!/usr/bin/env bash
# Launches the GCS stack locally: TimescaleDB (docker), Go/Echo backend, React/Vite frontend.
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SERVER_DIR="$ROOT_DIR/apps/server"
WEB_DIR="$ROOT_DIR/apps/web"

log() { printf '\033[1;36m[run]\033[0m %s\n' "$1"; }

# --- 1. Database ---
log "Starting TimescaleDB via docker compose..."
(cd "$ROOT_DIR" && docker compose up -d)

# apps/server/.env ships with credentials that don't match docker-compose (gcs:gcs).
# Normalize it so the backend can actually connect.
ENV_FILE="$SERVER_DIR/.env"
EXPECTED_DB_URL="postgres://gcs:gcs@localhost:5433/gcs?sslmode=disable"
if [ -f "$ENV_FILE" ] && ! grep -q "^DATABASE_URL=${EXPECTED_DB_URL}$" "$ENV_FILE"; then
  log "Fixing DATABASE_URL in apps/server/.env to match docker-compose (host port 5433)..."
  sed -i.bak "s#^DATABASE_URL=.*#DATABASE_URL=${EXPECTED_DB_URL}#" "$ENV_FILE"
elif [ ! -f "$ENV_FILE" ]; then
  log "Creating apps/server/.env from .env.example..."
  cp "$SERVER_DIR/.env.example" "$ENV_FILE"
fi

log "Waiting for Postgres to accept connections..."
until docker exec gcs-timescaledb pg_isready -U gcs -d gcs >/dev/null 2>&1; do
  sleep 1
done
# The Postgres image runs a two-phase startup (initdb, then a temp server that
# creates POSTGRES_DB before the real restart) — pg_isready can return true
# before the "gcs" database itself exists, so also wait until it's queryable.
until docker exec gcs-timescaledb psql -U gcs -d gcs -c 'select 1' >/dev/null 2>&1; do
  sleep 1
done

log "Running database migrations (sql/migrations)..."
MIGRATE_BIN="$(command -v migrate || echo "$(go env GOPATH)/bin/migrate")"
if [ ! -x "$MIGRATE_BIN" ]; then
  log "Installing golang-migrate CLI..."
  GOBIN="$(go env GOPATH)/bin" go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
  MIGRATE_BIN="$(go env GOPATH)/bin/migrate"
fi
"$MIGRATE_BIN" -path "$SERVER_DIR/sql/migrations" -database "$EXPECTED_DB_URL" up

# --- 2. Backend (Go/Echo, :8080) ---
if ! command -v go >/dev/null 2>&1; then
  log "ERROR: Go is not installed. Install it first, e.g.: sudo apt-get install -y golang-go"
  exit 1
fi

log "Installing Go dependencies..."
(cd "$SERVER_DIR" && go mod tidy)

if command -v swag >/dev/null 2>&1; then
  log "Regenerating swagger docs..."
  (cd "$SERVER_DIR" && swag init -g cmd/server/main.go >/dev/null)
elif [ -f "$(go env GOPATH)/bin/swag" ]; then
  log "Regenerating swagger docs..."
  (cd "$SERVER_DIR" && "$(go env GOPATH)/bin/swag" init -g cmd/server/main.go >/dev/null)
else
  log "swag CLI not found — skipping doc regeneration (install: go install github.com/swaggo/swag/cmd/swag@latest)"
fi

log "Starting backend on :8080..."
(cd "$SERVER_DIR" && go run cmd/server/main.go) &
SERVER_PID=$!

# --- 3. Frontend (React/Vite, :5173) ---
log "Installing JS dependencies (pnpm)..."
(cd "$ROOT_DIR" && pnpm install)

log "Starting frontend on :5173..."
(cd "$WEB_DIR" && pnpm dev) &
WEB_PID=$!

cleanup() {
  log "Shutting down..."
  kill "$SERVER_PID" "$WEB_PID" 2>/dev/null || true
  wait "$SERVER_PID" "$WEB_PID" 2>/dev/null || true
}
trap cleanup EXIT INT TERM

log "Backend:  http://localhost:8080/health"
log "Swagger:  http://localhost:8080/swagger/index.html"
log "Frontend: http://localhost:5173"
log "Press Ctrl+C to stop everything."
wait
