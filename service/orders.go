package service

import (
	"context"

	"time"
	"webdemo/model"
	"webdemo/repository"
	"webdemo/util"
)

type OrderService interface {
	CreateOrder(ctx context.Context, order *model.Order) error
}
type orderService struct {
	repo repository.OrdersRepository
}

func NewOrderService(repo repository.OrdersRepository) OrderService {
	return &orderService{
		repo: repo,
	}
}
func (s *orderService) CreateOrder(ctx context.Context, order *model.Order) error {
	// 设置默认值
	if order.Status == 0 {
		order.Status = 1 // 默认状态为待配送
	}

	// 设置时间字段
	now := time.Now().Format("2006-01-02 15:04:05")
	if order.CreateTime == "" {
		order.CreateTime = now
	}
	if order.UpdateTime == "" {
		order.UpdateTime = now
	}

	// 生成订单号
	if order.OrderNo == "" {
		order.OrderNo = util.GenerateOrderNo()
	}

	return s.repo.CreateOrder(ctx, order)
}
