package service

import (
	"context"
	"time"

	"github.com/cubenet-cms/cms/model"
	"github.com/cubenet-cms/cms/pkg/cache"
	"github.com/cubenet-cms/cms/store"
)

type NewsService struct {
	repo *store.NewsRepo
	c    *cache.Cache[[]model.News]
}

func NewNewsService(repo *store.NewsRepo) *NewsService {
	return &NewsService{
		repo: repo,
		c:    cache.New[[]model.News](5*time.Minute, 0),
	}
}

func (s *NewsService) CacheStats() cache.Stats {
	return s.c.Stats()
}

func (s *NewsService) List(ctx context.Context) ([]model.News, error) {
	if v, ok := s.c.Get("list"); ok {
		return v, nil
	}
	v, err := s.repo.List(ctx, 20)
	if err != nil {
		return nil, err
	}
	s.c.Set("list", v)
	return v, nil
}
