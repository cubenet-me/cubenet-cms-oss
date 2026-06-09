package service

import (
	"context"
	"time"

	"github.com/cubenet-cms/cms/model"
	"github.com/cubenet-cms/cms/pkg/cache"
	"github.com/cubenet-cms/cms/store"
)

type LauncherService struct {
	repo      *store.BuildRepo
	builds    *cache.Cache[[]model.Build]
	manifest  *cache.Cache[map[string]any]
	byID      *cache.Cache[*model.Build]
}

func NewLauncherService(repo *store.BuildRepo) *LauncherService {
	return &LauncherService{
		repo:     repo,
		builds:   cache.New[[]model.Build](60*time.Second, 0),
		manifest: cache.New[map[string]any](60*time.Second, 0),
		byID:     cache.New[*model.Build](60*time.Second, 0),
	}
}

func (s *LauncherService) CacheStats() cache.Stats {
	return s.builds.Stats()
}

func (s *LauncherService) ListBuilds(ctx context.Context) ([]model.Build, error) {
	if v, ok := s.builds.Get("list"); ok {
		return v, nil
	}
	v, err := s.repo.List(ctx, 50)
	if err != nil {
		return nil, err
	}
	s.builds.Set("list", v)
	return v, nil
}

func (s *LauncherService) Manifest(ctx context.Context) (map[string]any, error) {
	if v, ok := s.manifest.Get("manifest"); ok {
		return v, nil
	}
	builds, err := s.repo.List(ctx, 50)
	if err != nil {
		return nil, err
	}
	v := map[string]any{
		"version": "1.0.0",
		"builds":  builds,
	}
	s.manifest.Set("manifest", v)
	return v, nil
}

func (s *LauncherService) GetBuild(ctx context.Context, id string) (*model.Build, error) {
	if v, ok := s.byID.Get(id); ok {
		return v, nil
	}
	v, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	s.byID.Set(id, v)
	return v, nil
}
