package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool // shared global pool

func Init() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "<REDACTED>"
	}
	log.Printf("Connecting to database with DSN: %s", dsn)
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to DB: %v", err)
	}
	log.Println("Connected to the database successfully")

	DB = pool
}
