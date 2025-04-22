package db

import (
	"context"
	"fmt"
	"time"
	"translator/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDBPool(config config.DBConfig) (*pgxpool.Pool, error) {
	dsn := config.GetDSN()
	if dsn == "" {
		return nil, fmt.Errorf("database dsn is required but not found in config")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	// Ping to ensure it's valid
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}

	return pool, nil
}
