package repository

import (
	"context"
	"webdemo/model"

	"gorm.io/gorm"
)

type OrdersRepository interface {
	CreateOrder(ctx context.Context, order *model.Order) error
}
type ordersRepository struct {
	db *gorm.DB
}

func NewOrdersRepository(db *gorm.DB) OrdersRepository {
	return &ordersRepository{
		db: db,
	}
}
func (r *ordersRepository) CreateOrder(ctx context.Context, order *model.Order) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(order).Error
		if err != nil {
			return err
		}
		return nil
	})
}
