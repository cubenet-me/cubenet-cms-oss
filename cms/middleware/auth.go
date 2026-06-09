package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/cubenet-cms/cms/pkg/jwt"
)

type ctxKey string

const claimsKey ctxKey = "claims"

func Auth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
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

			claims, err := jwt.Validate(secret, token)
			if err != nil {
				http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), claimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetClaims(r *http.Request) *jwt.Claims {
	c, _ := r.Context().Value(claimsKey).(*jwt.Claims)
	return c
}
