package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"webdemo/model"

	"gorm.io/gorm"
)

type SettingsRepository interface {
	IsAddressExists(ctx context.Context, longitude, latitude float64, userId uint) (bool, error)
	CreateAddress(ctx context.Context, req *model.Address) error
	AddressList(ctx context.Context, userId uint) ([]*model.Address, error)
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
	log.Println("is_default", req.IsDefault)
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return errors.New("开启事务失败")
	}
	var (
		txErr error
	)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			return
		}
		if txErr == nil {
			if commitErr := tx.Commit(); commitErr != nil {
				tx.Rollback()
				txErr = errors.New("提交事务失败")
				return
			}
		}
	}()
	exists, err := r.IsAddressExists(ctx, req.Longitude, req.Latitude, req.UserID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("已存在相同地址")
	} else {
		if req.IsDefault == 1 {
			if err := tx.Model(&model.Address{}).Where("user_id=?", req.UserID).Update("is_default", 0).Error; err != nil {
				tx.Rollback()
				txErr = errors.New("更新默认地址失败")
				return txErr
			}
		}
		log.Println("after update is_default", req.IsDefault)
		if err := tx.Create(req).Error; err != nil {
			tx.Rollback()
			txErr = errors.New("新建地址失败")
			return txErr
		}
	}
	return txErr
}

func (r *settingsRepository) AddressList(ctx context.Context, userId uint) ([]*model.Address, error) {
	var addresses = make([]*model.Address, 0, 4)
	err := r.db.WithContext(ctx).Model(&model.Address{}).Where("user_id=?", userId).Scan(&addresses).Error
	if err != nil {
		return nil, errors.New("获取地址列表失败")
	}

	return addresses, nil
}
func (r *settingsRepository) IsAddressExists(ctx context.Context, longitude, latitude float64, userId uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Address{}).Where("longitude=? AND latitude=? AND user_id=?", longitude, latitude, userId).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("查询地址失败")
	}
	return count > 0, nil
}
