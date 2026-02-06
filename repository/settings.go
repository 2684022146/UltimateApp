package repository

import (
	"context"

	"gorm.io/gorm"
)

type SettingRepository interface {
	CreateAddress(ctx context.Context) error
}
type settingRepository struct {
	db gorm.DB
}

func NewSettingsRepository(db *gorm.DB) SettingRepository {
	return &settingRepository{
		db: *db,
	}
}
func (r *settingRepository) CreateAddress(ctx context.Context) error {
	return nil
}
