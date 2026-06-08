package handler

import (
	"encoding/json"
	"net/http"

	"github.com/cubenet-cms/backend/internal/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ManifestHandler struct {
	pool    *pgxpool.Pool
	storage *storage.Storage
}

func NewManifestHandler(pool *pgxpool.Pool, s3 *storage.Storage) *ManifestHandler {
	return &ManifestHandler{pool: pool, storage: s3}
}

func (h *ManifestHandler) Launcher(w http.ResponseWriter, r *http.Request) {
	type serverItem struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		Slug       string `json:"slug"`
		Version    string `json:"version"`
		Status     string `json:"status"`
		Players    int    `json:"players"`
		MaxPlayers int    `json:"max_players"`
	}

	type buildItem struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		Version   string `json:"version"`
		ModLoader string `json:"mod_loader"`
		MCVersion string `json:"mc_version"`
	}

	servers := make([]serverItem, 0)
	rows, err := h.pool.Query(r.Context(),
		"SELECT id, name, slug, version, status, players, max_players FROM servers ORDER BY name")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var s serverItem
			rows.Scan(&s.ID, &s.Name, &s.Slug, &s.Version, &s.Status, &s.Players, &s.MaxPlayers)
			servers = append(servers, s)
		}
	}

	builds := make([]buildItem, 0)
	rows2, err := h.pool.Query(r.Context(),
		"SELECT id, name, version, mod_loader, mc_version FROM builds ORDER BY created_at DESC LIMIT 50")
	if err == nil {
		defer rows2.Close()
		for rows2.Next() {
			var b buildItem
			rows2.Scan(&b.ID, &b.Name, &b.Version, &b.ModLoader, &b.MCVersion)
			builds = append(builds, b)
		}
	}

	json.NewEncoder(w).Encode(map[string]any{
		"version":     "1.0.0",
		"servers":     servers,
		"builds":      builds,
		"launcher_url": "", // будет заполнено позже
	})
}
