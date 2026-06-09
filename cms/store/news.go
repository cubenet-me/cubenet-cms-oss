package store

import (
	"context"

	"github.com/cubenet-cms/cms/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NewsRepo struct {
	pool *pgxpool.Pool
}

func NewNewsRepo(pool *pgxpool.Pool) *NewsRepo {
	return &NewsRepo{pool: pool}
}

func (r *NewsRepo) List(ctx context.Context, limit int) ([]model.News, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, title, content, created_at FROM news ORDER BY created_at DESC LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.News
	for rows.Next() {
		var n model.News
		if err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.CreatedAt); err != nil {
			continue
		}
		out = append(out, n)
	}
	return out, nil
}
