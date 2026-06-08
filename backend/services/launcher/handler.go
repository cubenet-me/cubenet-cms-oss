package launcher

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	pool *pgxpool.Pool
}

func NewHandler(pool *pgxpool.Pool) *Handler {
	return &Handler{pool: pool}
}

type buildItem struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Version   string `json:"version"`
	ModLoader string `json:"mod_loader"`
	MCVersion string `json:"mc_version"`
	FileHash  string `json:"file_hash"`
	FileSize  int64  `json:"file_size"`
}

func (h *Handler) ListBuilds(w http.ResponseWriter, r *http.Request) {
	rows, err := h.pool.Query(r.Context(), `
		SELECT id, name, version, mod_loader, mc_version, file_hash, file_size
		FROM builds ORDER BY created_at DESC LIMIT 50
	`)
	if err != nil {
		http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	builds := make([]buildItem, 0)
	for rows.Next() {
		var b buildItem
		if err := rows.Scan(&b.ID, &b.Name, &b.Version, &b.ModLoader, &b.MCVersion, &b.FileHash, &b.FileSize); err != nil {
			continue
		}
		builds = append(builds, b)
	}

	json.NewEncoder(w).Encode(builds)
}

func (h *Handler) Manifest(w http.ResponseWriter, r *http.Request) {
	rows, err := h.pool.Query(r.Context(), `
		SELECT id, name, version, mod_loader, mc_version, file_hash, file_size
		FROM builds ORDER BY created_at DESC LIMIT 50
	`)
	if err != nil {
		http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	builds := make([]buildItem, 0)
	for rows.Next() {
		var b buildItem
		if err := rows.Scan(&b.ID, &b.Name, &b.Version, &b.ModLoader, &b.MCVersion, &b.FileHash, &b.FileSize); err != nil {
			continue
		}
		builds = append(builds, b)
	}

	json.NewEncoder(w).Encode(map[string]any{
		"version": "1.0.0",
		"builds":  builds,
	})
}

func (h *Handler) GetBuild(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var b buildItem
	err := h.pool.QueryRow(r.Context(), `
		SELECT id, name, version, mod_loader, mc_version, file_hash, file_size
		FROM builds WHERE id = $1
	`, id).Scan(&b.ID, &b.Name, &b.Version, &b.ModLoader, &b.MCVersion, &b.FileHash, &b.FileSize)
	if err != nil {
		http.Error(w, `{"error":"build not found"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(b)
}
