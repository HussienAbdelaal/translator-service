package db

import (
	"context"
	"log"
	"time"
	"translator/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDBPool(config config.DBConfig) *pgxpool.Pool {
	dsn := config.GetDSN()
	if dsn == "" {
		log.Fatal("DATABASE_CONFIG is required but not found in environment variables")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	// Ping to ensure it's valid
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}

	return pool
}
