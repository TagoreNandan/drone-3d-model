# internal/database

Thin wrapper around a `pgxpool.Pool`.

- `InitDB()` — reads `DATABASE_URL` from the environment (falls back to a dummy local DSN if unset, which will simply fail to connect), opens the pool, and `Ping`s it to verify connectivity before returning. Called once from `main()`; fatals the process on failure.
- `GetPool()` — returns the shared `*pgxpool.Pool`; `main()` wraps it with `db.New()` and passes the resulting `*db.Queries` into each module's `Service`.
- `Close()` — closes the pool; deferred in `main()`.

No retry/backoff logic — if Postgres isn't reachable at startup, the server exits immediately. No connection pool tuning (`MaxConns`, etc.) is configured — it uses pgx's defaults.
