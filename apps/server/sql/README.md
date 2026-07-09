# sql

- `migrations/` — versioned Postgres/TimescaleDB schema migrations, run with `golang-migrate`. See `migrations/README.md`.
- `queries/` — sqlc query annotations, code-generated into `../internal/db`. See `queries/README.md`.

Both are referenced by `../sqlc.yaml` (`schema: "sql/migrations/"`, `queries: "sql/queries/"`, `out: "internal/db"`).
