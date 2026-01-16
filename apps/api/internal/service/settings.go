package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/mokuhyo-driven-test/api/internal/model"
	"github.com/mokuhyo-driven-test/api/internal/repository"
)

type SettingsService struct {
	settingsRepo *repository.SettingsRepository
}

func NewSettingsService(settingsRepo *repository.SettingsRepository) *SettingsService {
	return &SettingsService{settingsRepo: settingsRepo}
}

func (s *SettingsService) GetSettings(ctx context.Context, userID uuid.UUID) (*model.UserSettings, error) {
	return s.settingsRepo.GetByUserID(ctx, userID)
}

func (s *SettingsService) UpdateSettings(ctx context.Context, userID uuid.UUID, req model.UpdateSettingsRequest) (*model.UserSettings, error) {
	return s.settingsRepo.Upsert(ctx, userID, req)
}
