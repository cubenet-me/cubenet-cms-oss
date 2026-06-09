package v1

import (
	"encoding/json"
	"net/http"

	"github.com/cubenet-cms/cms/model"
	"github.com/cubenet-cms/cms/service"
	"github.com/go-chi/chi/v5"
)

type LauncherHandler struct {
	svc *service.LauncherService
}

func NewLauncherHandler(svc *service.LauncherService) *LauncherHandler {
	return &LauncherHandler{svc: svc}
}

func (h *LauncherHandler) Manifest(w http.ResponseWriter, r *http.Request) {
	result, err := h.svc.Manifest(r.Context())
	if err != nil {
		http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func (h *LauncherHandler) ListBuilds(w http.ResponseWriter, r *http.Request) {
	builds, err := h.svc.ListBuilds(r.Context())
	if err != nil {
		http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
		return
	}
	if builds == nil {
		builds = []model.Build{}
	}
	json.NewEncoder(w).Encode(builds)
}

func (h *LauncherHandler) GetBuild(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	b, err := h.svc.GetBuild(r.Context(), id)
	if err != nil {
		http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
		return
	}
	if b == nil {
		http.Error(w, `{"error":"build not found"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(b)
}
