package store

import (
	"context"

	"github.com/cubenet-cms/cms/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServerRepo struct {
	pool *pgxpool.Pool
}

func NewServerRepo(pool *pgxpool.Pool) *ServerRepo {
	return &ServerRepo{pool: pool}
}

func (r *ServerRepo) List(ctx context.Context) ([]model.Server, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, slug, description, address, version,
		       status, tps, players, max_players, mods
		FROM servers ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.Server
	for rows.Next() {
		var s model.Server
		if err := rows.Scan(&s.ID, &s.Name, &s.Slug, &s.Description, &s.Address,
			&s.Version, &s.Status, &s.TPS, &s.Players, &s.MaxPlayers, &s.Mods); err != nil {
			continue
		}
		out = append(out, s)
	}
	return out, nil
}

func (r *ServerRepo) GetBySlug(ctx context.Context, slug string) (*model.Server, error) {
	s := &model.Server{}
	err := r.pool.QueryRow(ctx, `
		SELECT id, name, slug, description, address, version,
		       status, tps, players, max_players, mods
		FROM servers WHERE slug = $1
	`, slug).Scan(&s.ID, &s.Name, &s.Slug, &s.Description, &s.Address,
		&s.Version, &s.Status, &s.TPS, &s.Players, &s.MaxPlayers, &s.Mods)
	if err != nil {
		return nil, err
	}
	return s, nil
}
