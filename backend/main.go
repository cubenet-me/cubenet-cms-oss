package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := LoadConfig()
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel})))

	pool, err := Connect(cfg.DatabaseURL)
	if err != nil {
		slog.Error("connect", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := Migrate(pool, MigrationsFS, "migrations"); err != nil {
		slog.Error("migrate", "error", err)
		os.Exit(1)
	}

	s3, err := NewStorage(StorageConfig{
		Endpoint:  cfg.S3Endpoint,
		AccessKey: cfg.S3AccessKey,
		SecretKey: cfg.S3SecretKey,
		Bucket:    cfg.S3Bucket,
		UseSSL:    cfg.S3UseSSL,
	})
	if err != nil {
		slog.Error("s3", "error", err)
		os.Exit(1)
	}
	_ = s3

	hub := NewHub()

	a := &App{
		pool:   pool,
		s3:     s3,
		hub:    hub,
		secret: cfg.JWTSecret,
		wsCfg: WSConfig{
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			PingInterval: 30 * time.Second,
			MaxMessageSz: 4096,
		},
	}

	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Post("/api/v1/auth/register", a.handleRegister)
	r.Post("/api/v1/auth/login", a.handleLogin)

	r.Get("/api/v1/servers", a.handleListServers)
	r.Get("/api/v1/servers/{slug}", a.handleGetServerBySlug)

	r.Get("/api/v1/launcher/manifest", a.handleManifest)
	r.Get("/api/v1/launcher/builds", a.handleListBuilds)
	r.Get("/api/v1/launcher/builds/{id}", a.handleGetBuild)

	r.Get("/api/v1/news", a.handleListNews)

	r.Handle("/api/v1/ws", http.HandlerFunc(a.handleWS))

	slog.Info("server started", "addr", cfg.Addr)
	if err := http.ListenAndServe(cfg.Addr, r); err != nil {
		slog.Error("server", "error", err)
		os.Exit(1)
	}
}
