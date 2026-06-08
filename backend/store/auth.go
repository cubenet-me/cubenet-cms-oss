package store

import (
	"context"

	"github.com/cubenet-cms/backend/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepo struct {
	pool *pgxpool.Pool
}

func NewAuthRepo(pool *pgxpool.Pool) *AuthRepo {
	return &AuthRepo{pool: pool}
}

func (r *AuthRepo) Create(ctx context.Context, username, email, password string) (string, error) {
	var id string
	err := r.pool.QueryRow(ctx,
		`INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id`,
		username, email, password,
	).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *AuthRepo) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	u := &model.User{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, username, email, password, role FROM users WHERE username = $1`,
		username,
	).Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Role)
	if err != nil {
		return nil, err
	}
	return u, nil
}
