package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type registerRequest struct {
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

func (a *App) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		http.Error(w, `{"error":"all fields required"}`, http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	var userID string
	err = a.pool.QueryRow(r.Context(),
		`INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id`,
		req.Username, req.Email, string(hash),
	).Scan(&userID)
	if err != nil {
		http.Error(w, `{"error":"username or email already exists"}`, http.StatusConflict)
		return
	}

	token, err := GenerateToken(a.secret, userID, req.Username, "user")
	if err != nil {
		http.Error(w, `{"error":"token generation failed"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(authResponse{
		Token:    token,
		UserID:   userID,
		Username: req.Username,
		Role:     "user",
	})
}

func (a *App) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	var userID, username, password, role string
	err := a.pool.QueryRow(r.Context(),
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

	token, err := GenerateToken(a.secret, userID, username, role)
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
