package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	v1 "github.com/cubenet-cms/cms/api/v1/auth"
	v1launcher "github.com/cubenet-cms/cms/api/v1/launcher"
	v1news "github.com/cubenet-cms/cms/api/v1/news"
	v1servers "github.com/cubenet-cms/cms/api/v1/servers"
	v1ws "github.com/cubenet-cms/cms/api/v1/ws"
	"github.com/cubenet-cms/cms/config"
	"github.com/cubenet-cms/cms/pkg/db"
	"github.com/cubenet-cms/cms/pkg/s3"
	"github.com/cubenet-cms/cms/pkg/ws"
	"github.com/cubenet-cms/cms/plugin"
	"github.com/cubenet-cms/cms/plugin/builtin"
	"github.com/cubenet-cms/cms/service"
	"github.com/cubenet-cms/cms/store"
	"github.com/cubenet-cms/cms/web"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	pipe := plugin.New()
	pipe.Register(builtin.NewSessionPlugin(authSvc))
	pipe.Register(builtin.NewFooterPlugin())

	webH := web.NewHandler(authSvc, serverSvc, newsSvc, pipe)

	r := chi.NewRouter()
	r.Use(middleware.RedirectSlashes)

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	// API
	r.Post("/api/v1/auth/register", authH.Register)
	r.Post("/api/v1/auth/login", authH.Login)

	r.Get("/api/v1/servers", serverH.List)
	r.Get("/api/v1/servers/{slug}", serverH.GetBySlug)

	r.Get("/api/v1/launcher/manifest", launcherH.Manifest)
	r.Get("/api/v1/launcher/builds", launcherH.ListBuilds)
	r.Get("/api/v1/launcher/builds/{id}", launcherH.GetBuild)

	r.Get("/api/v1/news", newsH.List)

	r.Handle("/api/v1/ws", wsH)

	// Web (Templ + htmx)
	r.Get("/static/*", webH.Static)
	r.Get("/assets/*", webH.Assets)
	r.Get("/", webH.Home)
	r.Get("/login", webH.LoginPage)
	r.Post("/auth/login", webH.Login)
	r.Get("/register", webH.RegisterPage)
	r.Post("/auth/register", webH.Register)
	r.Get("/servers", webH.Servers)
	r.Get("/admin", webH.Admin)

	slog.Info("server started", "addr", cfg.Addr)
	if err := http.ListenAndServe(cfg.Addr, r); err != nil {
		slog.Error("server", "error", err)
		os.Exit(1)
	}
}
