package v1

import (
	"encoding/json"
	"net/http"

	"github.com/cubenet-cms/backend/model"
	"github.com/cubenet-cms/backend/service"
)

type NewsHandler struct {
	svc *service.NewsService
}

func NewNewsHandler(svc *service.NewsService) *NewsHandler {
	return &NewsHandler{svc: svc}
}

func (h *NewsHandler) List(w http.ResponseWriter, r *http.Request) {
	news, err := h.svc.List(r.Context())
	if err != nil {
		http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
		return
	}
	if news == nil {
		news = []model.News{}
	}
	json.NewEncoder(w).Encode(news)
}
