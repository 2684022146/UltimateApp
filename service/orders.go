package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"webdemo/model"
	"webdemo/repository"
	"webdemo/util"
)

type OrderService interface {
	CreateOrder(ctx context.Context, req *model.CreateOrderRequest, userId uint) error
	OrderDetailBasic(ctx context.Context, orderNo string) (*model.OrderResponse, error)
	//
	SenderFinishedOrder(ctx context.Context, userId uint, pageInfo *model.Page) ([]*model.OrderResponse, int64, error)
	SenderInTransitOrder(ctx context.Context, userId uint, pageInfo *model.Page) ([]*model.OrderResponse, int64, error)
	SenderWaitingOrder(ctx context.Context, userId uint, pageInfo *model.Page) ([]*model.OrderResponse, int64, error)
	SenderCencelOrder(ctx context.Context, orderId string) error
	//

	ReceiverFinishedOrder(ctx context.Context, userId uint, pageInfo *model.Page) ([]*model.OrderResponse, int64, error)
	ReceiverInTransitOrder(ctx context.Context, userId uint, pageInfo *model.Page) ([]*model.OrderResponse, int64, error)
	//
	RiderOrderList(ctx context.Context, status int, pageInfo *model.Page) ([]*model.OrderResponse, int64, error)
	// 骑手相关
	AcceptOrder(ctx context.Context, orderNo string, deliveryUserID uint) error
	StartDelivery(ctx context.Context, orderNo string) error
	PickupOrder(ctx context.Context, orderNo string) error
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
	log.Println("开始创建订单")
	// 构建完整的收件地址
	completeAddress := fmt.Sprintf("%s%s%s%s%s", req.ReceiverProvince, req.ReceiverCity, req.ReceiverDistrict, req.ReceiverStreet, req.ReceiverDetail)
	log.Printf("完整地址: %s", completeAddress)
	encodedAddress := url.QueryEscape(completeAddress)
	url := fmt.Sprintf(apiUrl, encodedAddress)
	log.Printf("API URL: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("API请求失败: %v", err)
		return fmt.Errorf("api失败:%w", err)
	}
	defer resp.Body.Close()
	log.Printf("API响应状态码: %d", resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("读取API响应失败: %v", err)
		return fmt.Errorf("读取地址信息失败:%w", err)
	}
	log.Printf("API响应内容: %s", string(body))

	var geocodeResp model.GeocodeResponse
	err = json.Unmarshal(body, &geocodeResp)
	if err != nil {
		log.Printf("解析API响应失败: %v", err)
		return fmt.Errorf("解析地址信息失败:%w", err)
	}
	if geocodeResp.Status != "1" || len(geocodeResp.Geocodes) == 0 {
		log.Printf("API获取geo失败, status: %s, geocodes长度: %d", geocodeResp.Status, len(geocodeResp.Geocodes))
		return fmt.Errorf("api获取geo失败:%s", geocodeResp.Status)
	}
	coordinates := geocodeResp.Geocodes[0].Location
	log.Printf("获取到的坐标: %s", coordinates)
	parts := strings.Split(coordinates, ",")
	longitudeStr := parts[0]
	longitude, err := strconv.ParseFloat(longitudeStr, 64)
	if err != nil {
		log.Printf("解析经度失败: %v", err)
		return fmt.Errorf("解析经度失败:%w", err)
	}
	latitudeStr := parts[1]
	latitude, err := strconv.ParseFloat(latitudeStr, 64)
	if err != nil {
		log.Printf("解析纬度失败: %v", err)
		return fmt.Errorf("解析纬度失败:%w", err)
	}
	log.Printf("解析后的坐标: 经度=%f, 纬度=%f", longitude, latitude)

	orderId := util.GenerateOrderNo()
	log.Printf("生成的订单号: %s", orderId)
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
		ReceiverLatitude:  latitude,
		ReceiverLongitude: longitude,
	}

	log.Println("准备创建订单到数据库")
	if err := s.repo.CreateOrder(ctx, orderDetail); err != nil {
		log.Printf("数据库创建订单失败: %v", err)
		return fmt.Errorf("创建新订单失败:%w", err)
	}
	log.Println("订单创建成功")
	return nil
}
func (s *orderService) OrderDetailBasic(ctx context.Context, orderNo string) (*model.OrderResponse, error) {

	orderDetailBasic, err := s.repo.OrderDetailBasic(ctx, orderNo)
	if err != nil {
		return nil, fmt.Errorf("获取订单基本详情失败:%w", err)
	}
	return s.convertToOrderResponse(orderDetailBasic), nil
}

//

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
func (s *orderService) RiderOrderList(ctx context.Context, status int, pageInfo *model.Page) ([]*model.OrderResponse, int64, error) {
	orderList, total, err := s.repo.RiderOrderList(ctx, status, pageInfo)
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

// 取件
func (s *orderService) PickupOrder(ctx context.Context, orderNo string) error {
	if err := s.repo.PickupOrder(ctx, orderNo); err != nil {
		return fmt.Errorf("取件失败:%w", err)
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
