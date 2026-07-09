# sql/migrations

Versioned schema migrations, run with [golang-migrate](https://github.com/golang-migrate/migrate). This replaced the old single `sql/schema/001_schema.sql` file that was just re-applied with `psql` on every run.

Each migration is a numbered pair:

```
000001_init_schema.up.sql
000001_init_schema.down.sql
```

`sqlc.yaml` points `schema:` at this directory — sqlc reads every `*.up.sql` file in sequence to build the schema it type-checks queries against, and automatically ignores `*.down.sql` files.

## Creating a new migration

```sh
# from apps/server
migrate create -ext sql -dir sql/migrations -seq add_something
```

This creates `NNNNNN_add_something.up.sql` / `.down.sql` with the next sequence number. Write the forward change in `.up.sql` and its exact inverse in `.down.sql` (e.g. `CREATE TABLE` ↔ `DROP TABLE`, `ALTER TABLE ... ADD COLUMN` ↔ `... DROP COLUMN`).

After adding a migration, regenerate the sqlc client:

```sh
sqlc generate
```

## Applying migrations

```sh
# from apps/server, with DATABASE_URL exported or passed via -database
migrate -path sql/migrations -database "$DATABASE_URL" up      # apply all pending
migrate -path sql/migrations -database "$DATABASE_URL" down 1  # roll back the last one
migrate -path sql/migrations -database "$DATABASE_URL" version # show current version
```

`./run.sh` (repo root) runs `migrate ... up` automatically against the local TimescaleDB container, installing the `migrate` CLI first if it isn't already on `PATH` or in `$(go env GOPATH)/bin`.

Migration state is tracked in a `schema_migrations` table that `migrate` creates automatically in the target database — do not hand-edit it.

## Installing the CLI

```sh
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```
