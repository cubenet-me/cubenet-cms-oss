package main

import (
	"context"
	"net/http"
	"strings"
)

type ctxKey string

const claimsKey ctxKey = "claims"

func (a *App) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		if token == header {
			http.Error(w, `{"error":"invalid token format"}`, http.StatusUnauthorized)
			return
		}

		claims, err := ValidateToken(a.secret, token)
		if err != nil {
			http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), claimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getClaims(r *http.Request) *Claims {
	c, _ := r.Context().Value(claimsKey).(*Claims)
	return c
}
