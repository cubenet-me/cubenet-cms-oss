package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/cubenet-cms/backend/internal/config"
	"github.com/cubenet-cms/backend/internal/db"
	serverssvc "github.com/cubenet-cms/backend/services/servers"
	"github.com/go-chi/chi/v5"
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

	if err := db.Migrate(pool, serverssvc.MigrationsFS, "migrations"); err != nil {
		slog.Error("migrate", "error", err)
		os.Exit(1)
	}

	h := serverssvc.NewHandler(pool)
	r := chi.NewRouter()
	r.Get("/api/v1/servers", h.List)
	r.Get("/api/v1/servers/{slug}", h.GetBySlug)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	slog.Info("servers service started", "addr", cfg.Addr)
	http.ListenAndServe(cfg.Addr, r)
}
