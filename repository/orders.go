package repository

import (
	"context"
	"errors"
	"fmt"
	"webdemo/model"

	"gorm.io/gorm"
)

type OrdersRepository interface {
	OrderList(ctx context.Context, userId uint) ([]*model.OrderResponse, error)
	CreateOrder(ctx context.Context, req *model.Order) error
}

type ordersRepository struct {
	db *gorm.DB
}

func NewOrdersRepository(db *gorm.DB) OrdersRepository {
	return &ordersRepository{
		db: db,
	}
}

func (r *ordersRepository) OrderList(ctx context.Context, userId uint) ([]*model.OrderResponse, error) {
	var orders []*model.Order
	var orderResponses []*model.OrderResponse

	// 查询订单列表
	err := r.db.WithContext(ctx).Model(&model.Order{}).Where("user_id=? AND status IN (1, 2)", userId).Find(&orders).Error
	if err != nil {
		return nil, errors.New("获取订单列表失败")
	}

	// 遍历订单，查询每个订单的地址信息
	for _, order := range orders {
		var address model.Address
		err := r.db.WithContext(ctx).Where("id=?", order.AddressID).First(&address).Error
		if err != nil {
			// 打印错误日志，以便排查问题
			fmt.Printf("查询地址信息失败，address_id=%d, error=%v\n", order.AddressID, err)
			// 地址不存在时，使用空地址信息
			address = model.Address{}
		}

		// 构建响应结构体
		orderResponse := &model.OrderResponse{
			ID:             order.ID,
			UserID:         order.UserID,
			DeliveryUserID: order.DeliveryUserID,
			AddressID:      order.AddressID,
			OrderNo:        order.OrderNo,
			Status:         order.Status,
			GoodsInfo:      order.GoodsInfo,
			CreateTime:     order.CreateTime,
			UpdateTime:     order.UpdateTime,
			Address: model.AddressInfo{
				ID:        address.ID,
				UserID:    address.UserID,
				Province:  address.Province,
				City:      address.City,
				District:  address.District,
				Street:    address.Street,
				Detail:    address.Detail,
				Receiver:  address.Receiver,
				Phone:     address.Phone,
				Latitude:  float64(address.Latitude),
				Longitude: float64(address.Longitude),
				IsDefault: address.IsDefault,
			},
		}
		orderResponses = append(orderResponses, orderResponse)
	}

	return orderResponses, nil
}

func (r *ordersRepository) CreateOrder(ctx context.Context, req *model.Order) error {
	err := r.db.WithContext(ctx).Create(req).Error
	if err != nil {
		return errors.New("创建订单失败")
	}
	return nil
}
