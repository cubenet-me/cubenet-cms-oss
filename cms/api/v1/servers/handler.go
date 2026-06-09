package v1

import (
	"encoding/json"
	"net/http"

	"github.com/cubenet-cms/cms/model"
	"github.com/cubenet-cms/cms/service"
	"github.com/go-chi/chi/v5"
)

type ServerHandler struct {
	svc *service.ServerService
}

func NewServerHandler(svc *service.ServerService) *ServerHandler {
	return &ServerHandler{svc: svc}
}

func (h *ServerHandler) List(w http.ResponseWriter, r *http.Request) {
	servers, err := h.svc.List(r.Context())
	if err != nil {
		http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
		return
	}
	if servers == nil {
		servers = []model.Server{}
	}
	json.NewEncoder(w).Encode(servers)
}

func (h *ServerHandler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	s, err := h.svc.GetBySlug(r.Context(), slug)
	if err != nil {
		http.Error(w, `{"error":"server not found"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(s)
}
