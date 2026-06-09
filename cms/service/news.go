package service

import (
	"context"

	"github.com/cubenet-cms/cms/model"
	"github.com/cubenet-cms/cms/store"
)

type NewsService struct {
	repo *store.NewsRepo
}

func NewNewsService(repo *store.NewsRepo) *NewsService {
	return &NewsService{repo: repo}
}

func (s *NewsService) List(ctx context.Context) ([]model.News, error) {
	return s.repo.List(ctx, 20)
}
