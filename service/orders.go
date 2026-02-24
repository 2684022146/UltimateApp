package service

import (
	"context"
	"fmt"
	"time"
	"webdemo/model"
	"webdemo/repository"
	"webdemo/util"
)

type OrdersService interface {
	OrderList(ctx context.Context, userId uint) ([]*model.OrderResponse, error)
	CreateOrder(ctx context.Context, req *model.CreateOrderRequest, userId uint) error
}

type ordersService struct {
	repo repository.OrdersRepository
}

func NewOrdersService(repo repository.OrdersRepository) OrdersService {
	return &ordersService{
		repo: repo,
	}
}

func (s *ordersService) OrderList(ctx context.Context, userId uint) ([]*model.OrderResponse, error) {
	orders, err := s.repo.OrderList(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	return orders, nil
}

func (s *ordersService) CreateOrder(ctx context.Context, req *model.CreateOrderRequest, userId uint) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	order := &model.Order{
		UserID:     userId,
		AddressID:  req.AddressID,
		OrderNo:    util.GenerateOrderNo(),
		Status:     1, // 默认状态为待配送
		GoodsInfo:  req.GoodsInfo,
		CreateTime: now,
		UpdateTime: now,
	}
	err := s.repo.CreateOrder(ctx, order)
	if err != nil {
		return fmt.Errorf("创建订单失败:%w", err)
	}
	return nil
}
