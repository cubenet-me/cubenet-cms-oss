package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/cubenet-cms/backend/internal/config"
	"github.com/cubenet-cms/backend/internal/db"
	newssvc "github.com/cubenet-cms/backend/services/news"
)

func main() {
	cfg := config.Load()
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel})))

	connURL, err := db.CreateDatabaseIfNotExists(cfg.DatabaseURL, cfg.DBName)
	if err != nil {
		slog.Error("create database", "error", err)
		os.Exit(1)
	}

	pool, err := db.Connect(connURL)
	if err != nil {
		slog.Error("connect", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := db.Migrate(pool, newssvc.MigrationsFS, "migrations"); err != nil {
		slog.Error("migrate", "error", err)
		os.Exit(1)
	}

	h := newssvc.NewHandler(pool)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/news", h.List)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	slog.Info("news service started", "addr", cfg.Addr)
	http.ListenAndServe(cfg.Addr, mux)
}
