package database

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

// InitDB initializes the database connection pool
func InitDB() (*pgxpool.Pool, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://user:password@localhost:5432/dbname?sslmode=disable"
	}

	var err error
	Pool, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	// Verify connection
	if err := Pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return Pool, nil
}

// GetPool returns the database pool
func GetPool() *pgxpool.Pool {
	return Pool
}

// Close closes the database pool
func Close() {
	if Pool != nil {
		Pool.Close()
	}
}
