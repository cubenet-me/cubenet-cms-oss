package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	v1 "github.com/cubenet-cms/backend/api/v1/auth"
	v1launcher "github.com/cubenet-cms/backend/api/v1/launcher"
	v1news "github.com/cubenet-cms/backend/api/v1/news"
	v1servers "github.com/cubenet-cms/backend/api/v1/servers"
	v1ws "github.com/cubenet-cms/backend/api/v1/ws"
	"github.com/cubenet-cms/backend/config"
	"github.com/cubenet-cms/backend/pkg/db"
	"github.com/cubenet-cms/backend/pkg/s3"
	"github.com/cubenet-cms/backend/pkg/ws"
	"github.com/cubenet-cms/backend/service"
	"github.com/cubenet-cms/backend/store"
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

	if err := db.Migrate(pool, db.MigrationsFS, "migrations"); err != nil {
		slog.Error("migrate", "error", err)
		os.Exit(1)
	}

	st, err := s3.New(s3.Config{
		Endpoint: cfg.S3Endpoint, AccessKey: cfg.S3AccessKey,
		SecretKey: cfg.S3SecretKey, Bucket: cfg.S3Bucket, UseSSL: cfg.S3UseSSL,
	})
	if err != nil {
		slog.Error("s3", "error", err)
		os.Exit(1)
	}
	_ = st

	hub := ws.NewHub()

	authSvc := service.NewAuthService(store.NewAuthRepo(pool), cfg.JWTSecret)
	serverSvc := service.NewServerService(store.NewServerRepo(pool))
	launcherSvc := service.NewLauncherService(store.NewBuildRepo(pool))
	newsSvc := service.NewNewsService(store.NewNewsRepo(pool))

	authH := v1.NewAuthHandler(authSvc)
	serverH := v1servers.NewServerHandler(serverSvc)
	launcherH := v1launcher.NewLauncherHandler(launcherSvc)
	newsH := v1news.NewNewsHandler(newsSvc)
	wsH := v1ws.NewWSHandler(hub, 10*time.Second, 30*time.Second, 4096)

	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Post("/api/v1/auth/register", authH.Register)
	r.Post("/api/v1/auth/login", authH.Login)

	r.Get("/api/v1/servers", serverH.List)
	r.Get("/api/v1/servers/{slug}", serverH.GetBySlug)

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
