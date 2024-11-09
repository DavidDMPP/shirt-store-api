// internal/repository/setting_repository.go
package repository

import (
	"shirt-store-api/internal/domain"

	"gorm.io/gorm"
)

type SettingRepository struct {
    db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) *SettingRepository {
    return &SettingRepository{db}
}

func (r *SettingRepository) GetMidtransConfig() (*domain.MidtransConfig, error) {
    var config domain.MidtransConfig
    err := r.db.Where("is_active = ?", true).First(&config).Error
    if err != nil {
        return nil, err
    }
    return &config, nil
}

func (r *SettingRepository) UpdateMidtransConfig(config *domain.MidtransConfig) error {
    // Deactivate all existing configs
    r.db.Model(&domain.MidtransConfig{}).Update("is_active", false)
    
    // Create new config
    config.IsActive = true
    if config.ID > 0 {
        return r.db.Save(config).Error
    }
    return r.db.Create(config).Error
}