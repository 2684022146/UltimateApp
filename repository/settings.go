package repository

import (
	"context"
	"errors"
	"fmt"
	"webdemo/model"

	"gorm.io/gorm"
)

type SettingsRepository interface {
	IsAddressExists(ctx context.Context, longitude, latitude float64, userId int) (bool, error)
	CreateAddress(ctx context.Context, req *model.Address) error
	AddressList(ctx context.Context, userID int) ([]*model.Address, error)
}
type settingsRepository struct {
	db *gorm.DB
}

func NewSettingsRepository(db *gorm.DB) SettingsRepository {
	return &settingsRepository{
		db: db,
	}
}
func (r *settingsRepository) CreateAddress(ctx context.Context, req *model.Address) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return errors.New("开启事务失败")
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	exists, err := r.IsAddressExists(ctx, req.Longitude, req.Latitude, req.UserID)
	if err != nil {
		tx.Rollback()
		return err
	}
	if exists {
		tx.Rollback()
		return errors.New("已存在相同地址")
	} else {
		if req.IsDefault == 1 {
			if err := tx.Model(req).Where("user_id=? AND is_del=?", req.UserID, 0).Update("is_default", 0).Error; err != nil {
				tx.Rollback()
				return errors.New("更新默认地址失败")
			}
		}
		if err := tx.Create(req).Error; err != nil {
			tx.Rollback()
			return errors.New("新建地址失败")
		}
		if err := tx.Commit().Error; err != nil {
			return errors.New("提交事务失败")
		}
	}
	return nil
}
func (r *settingsRepository) AddressList(ctx context.Context, userId int) ([]*model.Address, error) {
	var addresses = make([]*model.Address, 0, 10)
	err := r.db.WithContext(ctx).Model(&model.Address{}).Where("user_id=?", userId).Scan(&addresses).Error
	if err != nil {
		return nil, errors.New("获取地址列表失败")
	}
	return addresses, nil
}
func (r *settingsRepository) IsAddressExists(ctx context.Context, longitude, latitude float64, userId int) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Address{}).Where("longitude=? AND latitude=? AND user_id=?", longitude, latitude, userId).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("查询地址失败")
	}
	return count > 0, nil
}
