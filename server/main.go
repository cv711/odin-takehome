package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/cv711/odin-takehome/server/api"
	"github.com/cv711/odin-takehome/server/db"
	"github.com/dusted-go/logging/prettylog"
	"github.com/jackc/pgx/v5/stdlib"
)

func main() {
	_, isProduction := os.LookupEnv("IS_PROD")
	var log *slog.Logger
	if isProduction {
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:       slog.LevelInfo,
			AddSource:   false,
			ReplaceAttr: nil,
		}))
	} else {
		log = slog.New(prettylog.NewHandler(&slog.HandlerOptions{
			Level:       slog.LevelDebug,
			AddSource:   false,
			ReplaceAttr: nil,
		}))
	}

	ctx := context.Background()
	dbPool := db.NewPool(ctx, log)
	if dbPool == nil {
		log.Error("Failed to create database pool")
		os.Exit(1)
	}

	internalDb := db.New(dbPool)
	if err := internalDb.Migrate(stdlib.OpenDBFromPool(dbPool)); err != nil {
		log.Error("Failed to migrate database: " + err.Error())
		os.Exit(1)
	}

	api.NewAPI(log, internalDb).Serve()
}
