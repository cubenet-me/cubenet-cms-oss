package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/cubenet-cms/backend/internal/config"
	"github.com/cubenet-cms/backend/internal/db"
	"github.com/cubenet-cms/backend/internal/storage"
	"github.com/cubenet-cms/backend/services/launcher"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.Load()
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel})))

	rootURL := cfg.DatabaseURL
	if err := db.CreateDatabaseIfNotExists(rootURL, cfg.DBName); err != nil {
		slog.Error("create database", "error", err)
		os.Exit(1)
	}

	pool, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		slog.Error("connect", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := db.Migrate(pool, launcher.MigrationsFS, "migrations"); err != nil {
		slog.Error("migrate", "error", err)
		os.Exit(1)
	}

	s3Cfg := storage.Config{
		Endpoint:  os.Getenv("S3_ENDPOINT"),
		AccessKey: os.Getenv("S3_ACCESS_KEY"),
		SecretKey: os.Getenv("S3_SECRET_KEY"),
		Bucket:    os.Getenv("S3_BUCKET"),
	}
	s3, err := storage.New(s3Cfg)
	if err != nil {
		slog.Error("s3", "error", err)
		os.Exit(1)
	}
	_ = s3

	h := launcher.NewHandler(pool)
	r := chi.NewRouter()
	r.Get("/api/v1/launcher/manifest", h.Manifest)
	r.Get("/api/v1/launcher/builds", h.ListBuilds)
	r.Get("/api/v1/launcher/builds/{id}", h.GetBuild)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	slog.Info("launcher service started", "addr", cfg.Addr)
	http.ListenAndServe(cfg.Addr, r)
}
