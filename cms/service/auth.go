package service

import (
	"context"
	"errors"
	"time"

	"github.com/cubenet-cms/cms/model"
	"github.com/cubenet-cms/cms/pkg/jwt"
	"github.com/cubenet-cms/cms/store"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo      *store.AuthRepo
	secret    string
	rolesByID map[string]*model.Role
}

func NewAuthService(repo *store.AuthRepo, secret string) *AuthService {
	return &AuthService{repo: repo, secret: secret}
}

func (s *AuthService) LoadRoles(ctx context.Context) error {
	roles, err := s.repo.GetAllRoles(ctx)
	if err != nil {
		return err
	}
	s.rolesByID = make(map[string]*model.Role, len(roles))
	for i := range roles {
		s.rolesByID[roles[i].Identifier] = &roles[i]
	}
	return nil
}

func (s *AuthService) GetPermissions(identifier string) []string {
	if s.rolesByID == nil {
		return nil
	}
	r, ok := s.rolesByID[identifier]
	if !ok {
		return nil
	}
	return r.Permissions
}

func (s *AuthService) GetRoleColor(identifier string) string {
	if s.rolesByID == nil {
		return "#94a3b8"
	}
	r, ok := s.rolesByID[identifier]
	if !ok {
		return "#94a3b8"
	}
	return r.Color
}

type RegisterResult struct {
	Token    string
	UserID   string
	Username string
	Role     string
}

func (s *AuthService) Register(ctx context.Context, username, email, password string) (*RegisterResult, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("internal error")
	}

	roleID, err := s.repo.GetDefaultRoleID(ctx)
	if err != nil {
		roleID = ""
	}

	id, err := s.repo.Create(ctx, username, email, string(hash), roleID)
	if err != nil {
		return nil, errors.New("username or email already exists")
	}

	token, err := jwt.Generate(s.secret, id, username, "user")
	if err != nil {
		return nil, errors.New("token generation failed")
	}

	return &RegisterResult{Token: token, UserID: id, Username: username, Role: "user"}, nil
}

type LoginResult struct {
	Token    string
	UserID   string
	Username string
	Role     string
}

func (s *AuthService) ValidateToken(tokenString string) (*jwt.Claims, error) {
	return jwt.Validate(s.secret, tokenString)
}

func (s *AuthService) GetProfile(ctx context.Context, userID string) (*model.User, error) {
	u, err := s.repo.GetProfile(ctx, userID)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	cleaned := make([]model.UserRole, 0, len(u.Roles))
	for _, r := range u.Roles {
		if r.ExpiresAt == "" {
			cleaned = append(cleaned, r)
			continue
		}
		t, err := time.Parse(time.RFC3339, r.ExpiresAt)
		if err != nil || t.After(now) {
			cleaned = append(cleaned, r)
		}
	}
	if len(cleaned) != len(u.Roles) {
		u.Roles = cleaned
		_ = s.repo.UpdateRoles(ctx, userID, cleaned)
	}
	return u, nil
}

func (s *AuthService) Login(ctx context.Context, username, password string) (*LoginResult, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("invalid credentials")
		}
		return nil, errors.New("internal error")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := jwt.Generate(s.secret, user.ID, user.Username, user.Role)
	if err != nil {
		return nil, errors.New("token generation failed")
	}

	return &LoginResult{Token: token, UserID: user.ID, Username: user.Username, Role: user.Role}, nil
}
