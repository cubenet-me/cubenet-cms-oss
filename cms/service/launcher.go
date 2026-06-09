package service

import (
	"context"

	"github.com/cubenet-cms/cms/model"
	"github.com/cubenet-cms/cms/store"
)

type LauncherService struct {
	repo *store.BuildRepo
}

func NewLauncherService(repo *store.BuildRepo) *LauncherService {
	return &LauncherService{repo: repo}
}

func (s *LauncherService) ListBuilds(ctx context.Context) ([]model.Build, error) {
	return s.repo.List(ctx, 50)
}

func (s *LauncherService) Manifest(ctx context.Context) (map[string]any, error) {
	builds, err := s.repo.List(ctx, 50)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"version": "1.0.0",
		"builds":  builds,
	}, nil
}

func (s *LauncherService) GetBuild(ctx context.Context, id string) (*model.Build, error) {
	return s.repo.GetByID(ctx, id)
}
