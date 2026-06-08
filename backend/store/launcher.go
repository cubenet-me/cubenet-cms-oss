package store

import (
	"context"

	"github.com/cubenet-cms/backend/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BuildRepo struct {
	pool *pgxpool.Pool
}

func NewBuildRepo(pool *pgxpool.Pool) *BuildRepo {
	return &BuildRepo{pool: pool}
}

func (r *BuildRepo) List(ctx context.Context, limit int) ([]model.Build, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, version, mod_loader, mc_version, file_hash, file_size
		FROM builds ORDER BY created_at DESC LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.Build
	for rows.Next() {
		var b model.Build
		if err := rows.Scan(&b.ID, &b.Name, &b.Version, &b.ModLoader, &b.MCVersion, &b.FileHash, &b.FileSize); err != nil {
			continue
		}
		out = append(out, b)
	}
	return out, nil
}

func (r *BuildRepo) GetByID(ctx context.Context, id string) (*model.Build, error) {
	b := &model.Build{}
	err := r.pool.QueryRow(ctx, `
		SELECT id, name, version, mod_loader, mc_version, file_hash, file_size
		FROM builds WHERE id = $1
	`, id).Scan(&b.ID, &b.Name, &b.Version, &b.ModLoader, &b.MCVersion, &b.FileHash, &b.FileSize)
	if err != nil {
		return nil, err
	}
	return b, nil
}
