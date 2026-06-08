package handler

import (
	"encoding/json"
	"net/http"

	"github.com/cubenet-cms/backend/internal/auth"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	pool    *pgxpool.Pool
	secret  string
}

func NewAuthHandler(pool *pgxpool.Pool, secret string) *AuthHandler {
	return &AuthHandler{pool: pool, secret: secret}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authResponse struct {
	Token    string `json:"token"`
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	var userID, username, password, role string
	err := h.pool.QueryRow(r.Context(),
		"SELECT id, username, password, role FROM users WHERE username = $1", req.Username,
	).Scan(&userID, &username, &password, &role)
	if err != nil {
		http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(req.Password)); err != nil {
		http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(h.secret, userID, username, role)
	if err != nil {
		http.Error(w, `{"error":"token generation failed"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(authResponse{
		Token:    token,
		UserID:   userID,
		Username: username,
		Role:     role,
	})
}
