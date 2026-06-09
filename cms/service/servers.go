package service

import (
	"context"

	"github.com/cubenet-cms/cms/model"
	"github.com/cubenet-cms/cms/store"
)

type ServerService struct {
	repo *store.ServerRepo
}

func NewServerService(repo *store.ServerRepo) *ServerService {
	return &ServerService{repo: repo}
}

func (s *ServerService) List(ctx context.Context) ([]model.Server, error) {
	return s.repo.List(ctx)
}

func (s *ServerService) GetBySlug(ctx context.Context, slug string) (*model.Server, error) {
	return s.repo.GetBySlug(ctx, slug)
}
