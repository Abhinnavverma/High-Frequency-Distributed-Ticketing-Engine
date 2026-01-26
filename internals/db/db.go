package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB holds the connection pool to the database
type DB struct {
	Pool *pgxpool.Pool
}

// New initializes the connection to the database
func NewDatabase(dsn string) (*DB, error) {

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database config: %w", err)
	}

	// Optimization: Connection Pooling settings
	config.MaxConns = 50
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Verify connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	log.Println("✅ Connected to Postgres via pgxPool")

	return &DB{Pool: pool}, nil
}

// Close gracefully shuts down the pool
func (d *DB) Close() {
	d.Pool.Close()
}
