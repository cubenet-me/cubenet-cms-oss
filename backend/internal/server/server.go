package server

import (
	"net/http"

	"github.com/cubenet-cms/backend/internal/config"
	"github.com/cubenet-cms/backend/internal/handler"
	"github.com/cubenet-cms/backend/internal/storage"
	"github.com/cubenet-cms/backend/internal/ws"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	cfg    *config.Config
	router chi.Router
	pool   *pgxpool.Pool
}

func New(cfg *config.Config, pool *pgxpool.Pool, s3 *storage.Storage, hub *ws.Hub) *Server {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	authH := handler.NewAuthHandler(pool, cfg.JWTSecret)
	serverH := handler.NewServerHandler(pool)
	manifestH := handler.NewManifestHandler(pool, s3)
	wsH := handler.NewWSHandler(hub, cfg.WS)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/login", authH.Login)

		r.Get("/servers", serverH.List)
		r.Get("/servers/{slug}", serverH.GetBySlug)

		r.Get("/manifest", manifestH.Launcher)

		r.Get("/ws", wsH.ServeHTTP)

		r.Route("/admin", func(r chi.Router) {
			r.Use(handler.AuthMiddleware(cfg.JWTSecret))
		})
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	return &Server{cfg: cfg, router: r, pool: pool}
}

func (s *Server) Start() error {
	return http.ListenAndServe(s.cfg.Addr, s.router)
}
