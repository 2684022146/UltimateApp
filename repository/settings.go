package repository

import (
	"context"
	"errors"
	"log"
	"webdemo/model"

	"gorm.io/gorm"
)

type SettingsRepository interface {
	CreateAddress(ctx context.Context, req *model.Address) error
	AddressList(ctx context.Context, userId uint) ([]*model.Address, error)
	AddressDetail(ctx context.Context, addressId, userId uint) (*model.Address, error)
	UpdateAddress(ctx context.Context, req *model.Address) error
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
func (r *settingsRepository) AddressDetail(ctx context.Context, addressId, userId uint) (*model.Address, error) {
	var address *model.Address
	err := r.db.Model(&model.Address{}).Where("id=? AND user_id=?", addressId, userId).Take(address).Error
	if err != nil {
		return nil, errors.New("获取地址详情失败")
	}
	return address, nil
}
func (r *settingsRepository) UpdateAddress(ctx context.Context, req *model.Address) error {
	tx := r.db.WithContext(ctx).Begin()
	var txErr error
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			return
		}
		if txErr == nil {
			if commitErr := tx.Commit(); commitErr != nil {
				tx.Rollback()
				return
			}
		}
	}()
	if req.IsDefault == 1 {
		err := tx.Model(&model.Address{}).Where("user_id=?", req.UserID).Update("is_default", 0).Error
		if err != nil {
			txErr = errors.New("更新默认地址失败")
			return txErr
		}
	}
	err := tx.Model(&model.Address{}).Where("id=? AND user_id=?", req.ID, req.UserID).Select("*").Updates(req).Error
	if err != nil {
		txErr = errors.New("更新地址失败")
		return txErr
	}
	return txErr
}
