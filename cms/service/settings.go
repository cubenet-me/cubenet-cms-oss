package service

import (
	"context"
	"sync"

	"github.com/cubenet-cms/cms/store"
)

type SettingsService struct {
	repo     *store.SettingsRepo
	mu       sync.RWMutex
	settings map[string]string
}

func NewSettingsService(repo *store.SettingsRepo) *SettingsService {
	return &SettingsService{repo: repo}
}

func (s *SettingsService) Load(ctx context.Context) error {
	settings, err := s.repo.GetAll(ctx)
	if err != nil {
		return err
	}
	s.mu.Lock()
	s.settings = settings
	s.mu.Unlock()
	return nil
}

func (s *SettingsService) GetAll() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.settings == nil {
		return map[string]string{}
	}
	out := make(map[string]string, len(s.settings))
	for k, v := range s.settings {
		out[k] = v
	}
	return out
}

func (s *SettingsService) Get(key, fallback string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.settings == nil {
		return fallback
	}
	v, ok := s.settings[key]
	if !ok {
		return fallback
	}
	return v
}

func (s *SettingsService) Set(ctx context.Context, key, value string) error {
	if err := s.repo.Set(ctx, key, value); err != nil {
		return err
	}
	s.mu.Lock()
	s.settings[key] = value
	s.mu.Unlock()
	return nil
}
