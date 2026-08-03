package config

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func DbConn(ctx context.Context) *pgxpool.Pool {
	cfg := GetConfig()

	pool, err := pgxpool.New(ctx, cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	return pool
}
