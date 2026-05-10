package repository

import (
	"github.com/gabutlabs/devopin/internal/model"
	"gorm.io/gorm"
)

type SettingRepository interface {
	GetSettings() (*model.AppSetting, error)
	UpdateSettings(setting *model.AppSetting) error
}

type settingRepository struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) SettingRepository {
	return &settingRepository{db: db}
}

func (r *settingRepository) GetSettings() (*model.AppSetting, error) {
	var setting model.AppSetting
	err := r.db.First(&setting, 1).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *settingRepository) UpdateSettings(setting *model.AppSetting) error {
	setting.ID = 1 // Ensure we always update the singleton record
	return r.db.Save(setting).Error
}
