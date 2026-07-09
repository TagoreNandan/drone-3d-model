# sql/queries

`sqlc` query annotations (`-- name: X :one|:many|:exec`) against the schema in `../schema/001_schema.sql`. Already code-generated into `internal/db` (per `sqlc.yaml`: `package: db`, `sql_package: pgx/v5`, JSON tags emitted, empty slices instead of nil).

| File | Queries | Used by |
|---|---|---|
| `vehicles.sql` | `ListVehicles`, `GetVehicle`, `CreateVehicle` | `internal/modules/vehicles` |
| `missions.sql` | `ListMissionsByVehicle`, `CreateMission`, `MarkMissionUploaded` | `internal/modules/missions` |
| `telemetry.sql` | `InsertTelemetryFrame`, `GetTelemetryRange` | `internal/modules/telemetry` |
| `alerts.sql` | `InsertAlert`, `ListRecentAlerts` | `internal/modules/alerts` |

## Regenerating after a schema/query change

```sh
# from apps/server
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest   # once, if not installed
sqlc generate
```

`internal/db` is generated code — never hand-edit it. Each domain module's `service.go` converts between `db.Queries` results and its own DTOs (via `internal/pgutil`), so a schema change only requires updating the module's `service.go`/`dto.go`, not its `controller.go`.
