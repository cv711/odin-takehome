package db

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, log *slog.Logger) *pgxpool.Pool {
	_, isProduction := os.LookupEnv("IS_PROD")

	log.Info("Connecting to database...")
	pgUser, found := os.LookupEnv("PG_USER")
	if !found || pgUser == "" {
		pgUser = "odin"
	}
	pgPassword, found := os.LookupEnv("PG_PASSWORD")
	if !found || pgPassword == "" {
		pgPassword = "exercise"
	}
	pgHost, hasPgHostEnv := os.LookupEnv("PG_HOST")
	if !hasPgHostEnv {
		pgHost = "localhost"
	}
	pgDatabase, found := os.LookupEnv("PG_DATABASE")
	if !found || pgDatabase == "" {
		pgDatabase = "odinexercise"
	}

	pingTimeout := 100 * time.Millisecond
	if isProduction {
		pingTimeout = 5 * time.Second
	}

	// Create a new connection pool to the database
	dbPool, err := pgxpool.New(ctx, fmt.Sprintf(`postgres://%s:%s@%s:5432/%s`, pgUser, pgPassword, pgHost, pgDatabase))
	if err != nil {
		log.ErrorContext(ctx, "Failed to connect to database: "+err.Error())
		return nil
	}

	pingCtx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()
	if dbPool.Ping(pingCtx) != nil {
		log.Error("Failed to connect to database at " + pgHost + ": ping failed")
		return nil
	}
	log.Info("Database connected.")

	return dbPool
}
