package main

import (
	"encoding/json"
	"net/http"
)

func (a *App) handleListNews(w http.ResponseWriter, r *http.Request) {
	rows, err := a.pool.Query(r.Context(), `
		SELECT id, title, content, created_at FROM news ORDER BY created_at DESC LIMIT 20
	`)
	if err != nil {
		http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type newsItem struct {
		ID        string `json:"id"`
		Title     string `json:"title"`
		Content   string `json:"content"`
		CreatedAt string `json:"created_at"`
	}

	news := make([]newsItem, 0)
	for rows.Next() {
		var n newsItem
		if err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.CreatedAt); err != nil {
			continue
		}
		news = append(news, n)
	}

	json.NewEncoder(w).Encode(news)
}
