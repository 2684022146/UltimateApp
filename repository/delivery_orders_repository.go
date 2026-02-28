package repository

import (
	"context"
	"errors"
	"fmt"
	"webdemo/model"

	"gorm.io/gorm"
)

type DeliveryOrdersRepository interface {
	DeliveryOrderList(ctx context.Context, deliveryUserId uint) ([]*model.OrderResponse, error)
	//DeliveryStartTask(ctx context.Context,)
}

type deliveryOrdersRepository struct {
	db *gorm.DB
}

func NewDeliveryOrdersRepository(db *gorm.DB) DeliveryOrdersRepository {
	return &deliveryOrdersRepository{
		db: db,
	}
}

func (r *deliveryOrdersRepository) DeliveryOrderList(ctx context.Context, deliveryUserId uint) ([]*model.OrderResponse, error) {
	var orders []*model.Order
	var orderResponses []*model.OrderResponse

	// 查询配送员待配送订单列表
	err := r.db.WithContext(ctx).Model(&model.Order{}).Where("delivery_user_id=? AND status=1", deliveryUserId).Find(&orders).Error
	if err != nil {
		return nil, errors.New("获取待配送订单列表失败")
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

