package v1

import (
	"encoding/json"
	"net/http"

	"github.com/cubenet-cms/cms/middleware"
	"github.com/cubenet-cms/cms/model"
	"github.com/cubenet-cms/cms/service"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

type authRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	Token    string `json:"token"`
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}
	if req.Username == "" || req.Password == "" {
		http.Error(w, `{"error":"username and password required"}`, http.StatusBadRequest)
		return
	}

	result, err := h.svc.Register(r.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(authResponse{
		Token: result.Token, UserID: result.UserID,
		Username: result.Username, Role: result.Role,
	})
}

type meResponse struct {
	UUID     string            `json:"uuid"`
	Nickname string            `json:"nickname"`
	Email    string            `json:"email"`
	Role     string            `json:"role"`
	RoleData *model.Role       `json:"role_data"`
	Roles    []model.UserRole  `json:"roles"`
	Wallet   model.UserWallet  `json:"wallet"`
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	if claims == nil {
		http.Error(w, `{"error":"not authenticated"}`, http.StatusUnauthorized)
		return
	}
	user, err := h.svc.GetProfile(r.Context(), claims.UserID)
	if err != nil {
		http.Error(w, `{"error":"user not found"}`, http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(meResponse{
		UUID:     user.ID,
		Nickname: user.Nickname,
		Email:    user.Email,
		Role:     user.Role,
		RoleData: user.RoleData,
		Roles:    user.Roles,
		Wallet:   user.Wallet,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}
	if req.Username == "" || req.Password == "" {
		http.Error(w, `{"error":"username and password required"}`, http.StatusBadRequest)
		return
	}

	result, err := h.svc.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(authResponse{
		Token: result.Token, UserID: result.UserID,
		Username: result.Username, Role: result.Role,
	})
}
