package service

import (
	"context"
	"fmt"
	"webdemo/model"
	"webdemo/repository"
)

type DeliveryOrdersService interface {
	DeliveryOrderList(ctx context.Context, deliveryUserId uint) ([]*model.OrderResponse, error)
}

type deliveryOrdersService struct {
	repo repository.DeliveryOrdersRepository
}

func NewDeliveryOrdersService(repo repository.DeliveryOrdersRepository) DeliveryOrdersService {
	return &deliveryOrdersService{
		repo: repo,
	}
}

func (s *deliveryOrdersService) DeliveryOrderList(ctx context.Context, deliveryUserId uint) ([]*model.OrderResponse, error) {
	orders, err := s.repo.DeliveryOrderList(ctx, deliveryUserId)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	return orders, nil
}
