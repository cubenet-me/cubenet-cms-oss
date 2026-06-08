package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/cubenet-cms/backend/internal/config"
	"github.com/cubenet-cms/backend/internal/db"
	"github.com/cubenet-cms/backend/internal/storage"
	"github.com/cubenet-cms/backend/internal/ws"
	authsvc "github.com/cubenet-cms/backend/services/auth"
	launcher "github.com/cubenet-cms/backend/services/launcher"
	newssvc "github.com/cubenet-cms/backend/services/news"
	serverssvc "github.com/cubenet-cms/backend/services/servers"
	wssvc "github.com/cubenet-cms/backend/services/ws"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.Load()
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel})))

	pool, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		slog.Error("connect", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := db.Migrate(pool, authsvc.MigrationsFS, "migrations"); err != nil {
		slog.Error("migrate auth", "error", err)
		os.Exit(1)
	}
	if err := db.Migrate(pool, serverssvc.MigrationsFS, "migrations"); err != nil {
		slog.Error("migrate servers", "error", err)
		os.Exit(1)
	}
	if err := db.Migrate(pool, launcher.MigrationsFS, "migrations"); err != nil {
		slog.Error("migrate launcher", "error", err)
		os.Exit(1)
	}
	if err := db.Migrate(pool, newssvc.MigrationsFS, "migrations"); err != nil {
		slog.Error("migrate news", "error", err)
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

	hub := ws.NewHub()

	authH := authsvc.NewHandler(pool, cfg.JWTSecret)
	serversH := serverssvc.NewHandler(pool)
	launcherH := launcher.NewHandler(pool)
	newsH := newssvc.NewHandler(pool)

	wsCfg := wssvc.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		PingInterval: 30 * time.Second,
		MaxMessageSz: 4096,
	}
	wsH := wssvc.NewHandler(hub, wsCfg)

	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Post("/api/v1/auth/register", authH.Register)
	r.Post("/api/v1/auth/login", authH.Login)

	r.Get("/api/v1/servers", serversH.List)
	r.Get("/api/v1/servers/{slug}", serversH.GetBySlug)

	r.Get("/api/v1/launcher/manifest", launcherH.Manifest)
	r.Get("/api/v1/launcher/builds", launcherH.ListBuilds)
	r.Get("/api/v1/launcher/builds/{id}", launcherH.GetBuild)

	r.Get("/api/v1/news", newsH.List)

	r.Handle("/api/v1/ws", wsH)

	slog.Info("server started", "addr", cfg.Addr)
	if err := http.ListenAndServe(cfg.Addr, r); err != nil {
		slog.Error("server", "error", err)
		os.Exit(1)
	}
}
