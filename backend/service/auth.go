package service

import (
	"context"
	"errors"

	"github.com/cubenet-cms/backend/pkg/jwt"
	"github.com/cubenet-cms/backend/store"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo   *store.AuthRepo
	secret string
}

func NewAuthService(repo *store.AuthRepo, secret string) *AuthService {
	return &AuthService{repo: repo, secret: secret}
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

	id, err := s.repo.Create(ctx, username, email, string(hash))
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
