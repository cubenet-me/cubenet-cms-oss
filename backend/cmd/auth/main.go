package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/cubenet-cms/backend/internal/config"
	"github.com/cubenet-cms/backend/internal/db"
	authsvc "github.com/cubenet-cms/backend/services/auth"
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

	if err := db.Migrate(pool, authsvc.MigrationsFS, "migrations"); err != nil {
		slog.Error("migrate", "error", err)
		os.Exit(1)
	}

	h := authsvc.NewHandler(pool, cfg.JWTSecret)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/auth/register", h.Register)
	mux.HandleFunc("POST /api/v1/auth/login", h.Login)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	slog.Info("auth service started", "addr", cfg.Addr)
	http.ListenAndServe(cfg.Addr, mux)
}
