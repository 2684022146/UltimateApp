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
	SenderFinishedOrder(ctx context.Context, userId uint, pageInfo *model.Page) ([]*model.OrderResponse, int64, error)
	SenderInTransitOrder(ctx context.Context, userId uint, pageInfo *model.Page) ([]*model.OrderResponse, int64, error)
	SenderWaitingOrder(ctx context.Context, userId uint, pageInfo *model.Page) ([]*model.OrderResponse, int64, error)
	SenderCencelOrder(ctx context.Context, orderId string) error
	//

	ReceiverFinishedOrder(ctx context.Context, userId uint, pageInfo *model.Page) ([]*model.OrderResponse, int64, error)
	ReceiverInTransitOrder(ctx context.Context, userId uint, pageInfo *model.Page) ([]*model.OrderResponse, int64, error)
	//
	RiderOrderList(ctx context.Context, pageInfo *model.Page) ([]*model.OrderResponse, int64, error)
	// 骑手相关
	AcceptOrder(ctx context.Context, orderNo string, deliveryUserID uint) error
	StartDelivery(ctx context.Context, orderNo string) error
	UploadLocation(ctx context.Context, orderID uint, deliveryUserID uint, longitude, latitude float64) error
	CompleteOrder(ctx context.Context, orderNo string) error
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
		Status:            1,
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

func (s *orderService) convertToOrderResponse(order *model.Order) *model.OrderResponse {
	return &model.OrderResponse{
		ID:                order.ID,
		SenderUserID:      order.SenderUserID,
		SenderAddressID:   order.SenderAddressID,
		DeliveryUserID:    order.DeliveryUserID,
		OrderNo:           order.OrderNo,
		Status:            order.Status,
		GoodsInfo:         order.GoodsInfo,
		ReceiverName:      order.ReceiverName,
		ReceiverPhone:     order.ReceiverPhone,
		ReceiverProvince:  order.ReceiverProvince,
		ReceiverCity:      order.ReceiverCity,
		ReceiverDistrict:  order.ReceiverDistrict,
		ReceiverStreet:    order.ReceiverStreet,
		ReceiverDetail:    order.ReceiverDetail,
		ReceiverLatitude:  order.ReceiverLatitude,
		ReceiverLongitude: order.ReceiverLongitude,
		CreateTime:        order.CreateTime.Format("2006-01-02 15:04:05"),
		UpdateTime:        order.UpdateTime.Format("2006-01-02 15:04:05"),
		SenderAddress: model.AddressInfo{
			Province:  order.SenderProvince,
			City:      order.SenderCity,
			District:  order.SenderDistrict,
			Street:    order.SenderStreet,
			Detail:    order.SenderDetail,
			Receiver:  order.SenderReceiver,
			Phone:     order.SenderPhone,
			Latitude:  order.SenderLatitude,
			Longitude: order.SenderLongitude,
		},
	}
}

func (s *orderService) SenderFinishedOrder(ctx context.Context, userId uint, pageInfo *model.Page) ([]*model.OrderResponse, int64, error) {
	orderList, total, err := s.repo.SenderFinishedOrder(ctx, userId, pageInfo)
	if err != nil {
		return nil, 0, fmt.Errorf("获取寄件人已完成订单失败:%w", err)
	}
	responseList := make([]*model.OrderResponse, 0, len(orderList))
	for _, order := range orderList {
		responseList = append(responseList, s.convertToOrderResponse(order))
	}
	return responseList, total, nil
}

func (s *orderService) SenderInTransitOrder(ctx context.Context, userId uint, pageInfo *model.Page) ([]*model.OrderResponse, int64, error) {
	orderList, total, err := s.repo.SenderInTransitOrder(ctx, userId, pageInfo)
	if err != nil {
		return nil, 0, fmt.Errorf("获取寄件人在途中订单失败:%w", err)
	}
	responseList := make([]*model.OrderResponse, 0, len(orderList))
	for _, order := range orderList {
		responseList = append(responseList, s.convertToOrderResponse(order))
	}
	return responseList, total, nil
}
func (s *orderService) SenderWaitingOrder(ctx context.Context, userId uint, pageInfo *model.Page) ([]*model.OrderResponse, int64, error) {
	orderList, total, err := s.repo.SenderWaitingOrder(ctx, userId, pageInfo)
	if err != nil {
		return nil, 0, fmt.Errorf("获取待接单订单失败:%w", err)
	}
	responseList := make([]*model.OrderResponse, 0, len(orderList))
	for _, order := range orderList {
		responseList = append(responseList, s.convertToOrderResponse(order))
	}
	return responseList, total, nil
}
func (s *orderService) SenderCencelOrder(ctx context.Context, orderNo string) error {
	if err := s.repo.SenderCencelOrder(ctx, orderNo); err != nil {
		return fmt.Errorf("取消订单失败:%w", err)
	}
	return nil
}

// //////////////////////////
func (s *orderService) ReceiverFinishedOrder(ctx context.Context, userId uint, pageInfo *model.Page) ([]*model.OrderResponse, int64, error) {
	orderList, total, err := s.repo.ReceiverFinishedOrder(ctx, userId, pageInfo)
	if err != nil {
		return nil, 0, fmt.Errorf("获取收件人已完成订单失败:%w", err)
	}
	responseList := make([]*model.OrderResponse, 0, len(orderList))
	for _, order := range orderList {
		responseList = append(responseList, s.convertToOrderResponse(order))
	}
	return responseList, total, nil
}
func (s *orderService) ReceiverInTransitOrder(ctx context.Context, userId uint, pageInfo *model.Page) ([]*model.OrderResponse, int64, error) {
	orderList, total, err := s.repo.ReceiverInTransitOrder(ctx, userId, pageInfo)
	if err != nil {
		return nil, 0, fmt.Errorf("获取收件人途中订单失败:%w", err)
	}
	responseList := make([]*model.OrderResponse, 0, len(orderList))
	for _, order := range orderList {
		responseList = append(responseList, s.convertToOrderResponse(order))
	}
	return responseList, total, nil
}

// ///////////////////////
func (s *orderService) RiderOrderList(ctx context.Context, pageInfo *model.Page) ([]*model.OrderResponse, int64, error) {
	orderList, total, err := s.repo.RiderOrderList(ctx, pageInfo)
	if err != nil {
		return nil, 0, fmt.Errorf("获取待接单订单失败:%w", err)
	}
	responseList := make([]*model.OrderResponse, 0, len(orderList))
	for _, order := range orderList {
		responseList = append(responseList, s.convertToOrderResponse(order))
	}
	return responseList, total, nil
}

// 骑手接单
func (s *orderService) AcceptOrder(ctx context.Context, orderNo string, deliveryUserID uint) error {
	if err := s.repo.AcceptOrder(ctx, orderNo, deliveryUserID); err != nil {
		return fmt.Errorf("接单失败:%w", err)
	}
	return nil
}

// 开始配送
func (s *orderService) StartDelivery(ctx context.Context, orderNo string) error {
	if err := s.repo.StartDelivery(ctx, orderNo); err != nil {
		return fmt.Errorf("开始配送失败:%w", err)
	}
	return nil
}

// 上传位置
func (s *orderService) UploadLocation(ctx context.Context, orderID uint, deliveryUserID uint, longitude, latitude float64) error {
	if err := s.repo.UploadLocation(ctx, orderID, deliveryUserID, longitude, latitude); err != nil {
		return fmt.Errorf("上传位置失败:%w", err)
	}
	return nil
}

// 完成订单
func (s *orderService) CompleteOrder(ctx context.Context, orderNo string) error {
	if err := s.repo.CompleteOrder(ctx, orderNo); err != nil {
		return fmt.Errorf("完成订单失败:%w", err)
	}
	return nil
}
