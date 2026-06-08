package main

import (
	"log/slog"
	"os"

	"github.com/cubenet-cms/backend/internal/config"
	"github.com/cubenet-cms/backend/internal/db"
	"github.com/cubenet-cms/backend/internal/server"
	"github.com/cubenet-cms/backend/internal/storage"
	"github.com/cubenet-cms/backend/internal/ws"
)

func main() {
	cfg := config.Load()
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel})))

	pool, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		slog.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := db.Migrate(pool); err != nil {
		slog.Error("migration failed", "error", err)
		os.Exit(1)
	}

	s3, err := storage.New(cfg.S3)
	if err != nil {
		slog.Error("s3 connection failed", "error", err)
		os.Exit(1)
	}

	hub := ws.NewHub()
	go hub.Run()

	srv := server.New(cfg, pool, s3, hub)
	slog.Info("starting server", "addr", cfg.Addr)
	if err := srv.Start(); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}
