package store

import (
	"context"

	"github.com/cubenet-cms/cms/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepo struct {
	pool *pgxpool.Pool
}

func NewAuthRepo(pool *pgxpool.Pool) *AuthRepo {
	return &AuthRepo{pool: pool}
}

func (r *AuthRepo) Create(ctx context.Context, username, email, password, roleID string) (string, error) {
	var id string
	err := r.pool.QueryRow(ctx,
		`INSERT INTO users (username, email, password, role_id) VALUES ($1, $2, $3, $4) RETURNING id`,
		username, email, password, roleID,
	).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *AuthRepo) GetDefaultRoleID(ctx context.Context) (string, error) {
	var id string
	err := r.pool.QueryRow(ctx, `SELECT id FROM roles WHERE identifier = 'user'`).Scan(&id)
	return id, err
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

func (r *AuthRepo) GetProfile(ctx context.Context, userID string) (*model.User, error) {
	u := &model.User{RoleData: &model.Role{}}
	err := r.pool.QueryRow(ctx,
		`SELECT u.id, u.username, COALESCE(u.nickname, u.username), u.email,
		        u.role, u.roles, u.wallet,
		        r.identifier, r.name, r.color, r.permissions
		 FROM users u
		 LEFT JOIN roles r ON r.id = u.role_id
		 WHERE u.id = $1`,
		userID,
	).Scan(&u.ID, &u.Username, &u.Nickname, &u.Email,
		&u.Role, &u.Roles, &u.Wallet,
		&u.RoleData.Identifier, &u.RoleData.Name, &u.RoleData.Color, &u.RoleData.Permissions,
	)
	if err != nil {
		return nil, err
	}
	u.RoleID = u.RoleData.Identifier
	return u, nil
}

func (r *AuthRepo) GetAllRoles(ctx context.Context) ([]model.Role, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, identifier, name, color, permissions FROM roles ORDER BY created_at`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []model.Role
	for rows.Next() {
		var role model.Role
		if err := rows.Scan(&role.ID, &role.Identifier, &role.Name, &role.Color, &role.Permissions); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *AuthRepo) UpdateRoles(ctx context.Context, userID string, roles []model.UserRole) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE users SET roles = $1 WHERE id = $2`,
		roles, userID,
	)
	return err
}
