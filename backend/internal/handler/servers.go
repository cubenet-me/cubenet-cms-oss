package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServerHandler struct {
	pool *pgxpool.Pool
}

func NewServerHandler(pool *pgxpool.Pool) *ServerHandler {
	return &ServerHandler{pool: pool}
}

func (h *ServerHandler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.pool.Query(r.Context(), `
		SELECT id, name, slug, description, address, version,
		       status, tps, players, max_players, mods
		FROM servers ORDER BY name
	`)
	if err != nil {
		http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type serverItem struct {
		ID         string   `json:"id"`
		Name       string   `json:"name"`
		Slug       string   `json:"slug"`
		Description string  `json:"description"`
		Address    string   `json:"address"`
		Version    string   `json:"version"`
		Status     string   `json:"status"`
		TPS        float64  `json:"tps"`
		Players    int      `json:"players"`
		MaxPlayers int      `json:"max_players"`
		Mods       []string `json:"mods"`
	}

	servers := make([]serverItem, 0)
	for rows.Next() {
		var s serverItem
		if err := rows.Scan(&s.ID, &s.Name, &s.Slug, &s.Description, &s.Address,
			&s.Version, &s.Status, &s.TPS, &s.Players, &s.MaxPlayers, &s.Mods); err != nil {
			continue
		}
		servers = append(servers, s)
	}

	json.NewEncoder(w).Encode(servers)
}

func (h *ServerHandler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	var s struct {
		ID         string   `json:"id"`
		Name       string   `json:"name"`
		Slug       string   `json:"slug"`
		Description string  `json:"description"`
		Address    string   `json:"address"`
		Version    string   `json:"version"`
		Status     string   `json:"status"`
		TPS        float64  `json:"tps"`
		Players    int      `json:"players"`
		MaxPlayers int      `json:"max_players"`
		Mods       []string `json:"mods"`
	}

	err := h.pool.QueryRow(r.Context(), `
		SELECT id, name, slug, description, address, version,
		       status, tps, players, max_players, mods
		FROM servers WHERE slug = $1
	`, slug).Scan(&s.ID, &s.Name, &s.Slug, &s.Description, &s.Address,
		&s.Version, &s.Status, &s.TPS, &s.Players, &s.MaxPlayers, &s.Mods)
	if err != nil {
		http.Error(w, `{"error":"server not found"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(s)
}
