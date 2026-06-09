package service

import (
	"context"
	"time"

	"github.com/cubenet-cms/cms/model"
	"github.com/cubenet-cms/cms/pkg/cache"
	"github.com/cubenet-cms/cms/store"
)

type ServerService struct {
	repo  *store.ServerRepo
	list  *cache.Cache[[]model.Server]
	bySlug *cache.Cache[*model.Server]
}

func NewServerService(repo *store.ServerRepo) *ServerService {
	return &ServerService{
		repo:   repo,
		list:   cache.New[[]model.Server](30*time.Second, 0),
		bySlug: cache.New[*model.Server](30*time.Second, 0),
	}
}

func (s *ServerService) List(ctx context.Context) ([]model.Server, error) {
	if v, ok := s.list.Get("list"); ok {
		return v, nil
	}
	v, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	s.list.Set("list", v)
	return v, nil
}

func (s *ServerService) CacheStats() cache.Stats {
	return s.list.Stats()
}

func (s *ServerService) GetBySlug(ctx context.Context, slug string) (*model.Server, error) {
	if v, ok := s.bySlug.Get(slug); ok {
		return v, nil
	}
	v, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	s.bySlug.Set(slug, v)
	return v, nil
}
