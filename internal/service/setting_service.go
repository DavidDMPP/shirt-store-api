// internal/service/setting_service.go
package service

import (
	"shirt-store-api/internal/domain"
	"shirt-store-api/internal/repository"
)

type SettingService struct {
    repo *repository.SettingRepository
}

func NewSettingService(repo *repository.SettingRepository) *SettingService {
    return &SettingService{repo}
}

func (s *SettingService) GetMidtransConfig() (*domain.MidtransConfig, error) {
    return s.repo.GetMidtransConfig()
}

func (s *SettingService) UpdateMidtransConfig(config *domain.MidtransConfig) error {
    return s.repo.UpdateMidtransConfig(config)
}