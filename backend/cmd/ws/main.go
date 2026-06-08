package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/cubenet-cms/backend/internal/config"
	"github.com/cubenet-cms/backend/internal/ws"
	wssvc "github.com/cubenet-cms/backend/services/ws"
)

func main() {
	cfg := config.Load()
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel})))

	hub := ws.NewHub()

	wsCfg := wssvc.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		PingInterval: 30 * time.Second,
		MaxMessageSz: 4096,
	}

	h := wssvc.NewHandler(hub, wsCfg)
	mux := http.NewServeMux()
	mux.Handle("/api/v1/ws", h)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	slog.Info("ws service started", "addr", cfg.Addr)
	http.ListenAndServe(cfg.Addr, mux)
}
