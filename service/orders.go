package service

import (
	"context"
	"fmt"
	"time"

	"webdemo/model"
	"webdemo/repository"
	"webdemo/util"
)

type OrderService interface {
	CreateOrder(ctx context.Context, req *model.CreateOrderRequest, userId uint) error
}
type orderService struct {
	repo repository.OrdersRepository
}

func NewOrdersService(repo repository.OrdersRepository) OrderService {
	return &orderService{
		repo: repo,
	}
}
func (s *orderService) CreateOrder(ctx context.Context, req *model.CreateOrderRequest, userId uint) error {
	orderId := util.GenerateOrderNo()
	orderDetail := &model.Order{
		SenderUserID:      userId,
		SenderAddressID:   req.SenderAddressID,
		DeliveryUserID:    req.DeliveryUserID,
		OrderNo:           orderId,
		Status:            1, // 待接单状态
		GoodsInfo:         req.GoodsInfo,
		ReceiverName:      req.ReceiverName,
		ReceiverPhone:     req.ReceiverPhone,
		ReceiverProvince:  req.ReceiverProvince,
		ReceiverCity:      req.ReceiverCity,
		ReceiverDistrict:  req.ReceiverDistrict,
		ReceiverStreet:    req.ReceiverStreet,
		ReceiverDetail:    req.ReceiverDetail,
		ReceiverLatitude:  req.ReceiverLatitude,
		ReceiverLongitude: req.ReceiverLongitude,
		CreateTime:        time.Now(),
		UpdateTime:        time.Now(),
	}

	if err := s.repo.CreateOrder(ctx, orderDetail); err != nil {
		return fmt.Errorf("创建新订单失败:%w", err)
	}
	return nil
}
