package model

import (
	"time"
)

// Order 订单结构体
type Order struct {
	ID                uint      `gorm:"primaryKey;column:id;autoIncrement" json:"id"`
	SenderUserID      uint      `gorm:"column:sender_user_id;not null" json:"sender_user_id"`
	SenderAddressID   uint      `gorm:"column:sender_address_id;not null" json:"sender_address_id"`
	OrderNo           string    `gorm:"column:order_no;not null;unique" json:"order_no"`
	Status            int8      `gorm:"column:status;not null;default:1" json:"status"`
	GoodsInfo         string    `gorm:"column:goods_info" json:"goods_info"`
	ReceiverName      string    `gorm:"column:receiver_name;not null" json:"receiver_name"`
	ReceiverPhone     string    `gorm:"column:receiver_phone;not null" json:"receiver_phone"`
	ReceiverProvince  string    `gorm:"column:receiver_province;not null" json:"receiver_province"`
	ReceiverCity      string    `gorm:"column:receiver_city;not null" json:"receiver_city"`
	ReceiverDistrict  string    `gorm:"column:receiver_district" json:"receiver_district"`
	ReceiverStreet    string    `gorm:"column:receiver_street" json:"receiver_street"`
	ReceiverDetail    string    `gorm:"column:receiver_detail;not null" json:"receiver_detail"`
	ReceiverLatitude  float64   `gorm:"column:receiver_latitude;not null" json:"receiver_latitude"`
	ReceiverLongitude float64   `gorm:"column:receiver_longitude;not null" json:"receiver_longitude"`
	DeliveryUserID    *uint     `gorm:"column:delivery_user_id" json:"delivery_user_id"`
	CreateTime        time.Time `gorm:"column:create_time;not null;default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime        time.Time `gorm:"column:update_time;not null;default:CURRENT_TIMESTAMP;autoUpdateTime" json:"update_time"`
	// 寄件人地址信息（通过 JOIN 查询填充）
	SenderProvince  string  `gorm:"-" json:"sender_province"`
	SenderCity      string  `gorm:"-" json:"sender_city"`
	SenderDistrict  string  `gorm:"-" json:"sender_district"`
	SenderStreet    string  `gorm:"-" json:"sender_street"`
	SenderDetail    string  `gorm:"-" json:"sender_detail"`
	SenderReceiver  string  `gorm:"-" json:"sender_receiver"`
	SenderPhone     string  `gorm:"-" json:"sender_phone"`
	SenderLatitude  float64 `gorm:"-" json:"sender_latitude"`
	SenderLongitude float64 `gorm:"-" json:"sender_longitude"`
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	SenderAddressID   uint    `json:"sender_address_id" binding:"required"`
	DeliveryUserID    *uint   `json:"delivery_user_id"`
	GoodsInfo         string  `json:"goods_info" binding:"required"`
	ReceiverName      string  `json:"receiver_name" binding:"required"`
	ReceiverPhone     string  `json:"receiver_phone" binding:"required"`
	ReceiverProvince  string  `json:"receiver_province" binding:"required"`
	ReceiverCity      string  `json:"receiver_city" binding:"required"`
	ReceiverDistrict  string  `json:"receiver_district"`
	ReceiverStreet    string  `json:"receiver_street"`
	ReceiverDetail    string  `json:"receiver_detail" binding:"required"`
	ReceiverLatitude  float64 `json:"receiver_latitude" binding:"required"`
	ReceiverLongitude float64 `json:"receiver_longitude" binding:"required"`
}

// OrderResponse 订单响应结构体（包含地址信息）
type OrderResponse struct {
	ID                uint        `json:"id"`
	SenderUserID      uint        `json:"sender_user_id"`
	SenderAddressID   uint        `json:"sender_address_id"`
	DeliveryUserID    *uint       `json:"delivery_user_id"`
	OrderNo           string      `json:"order_no"`
	Status            int8        `json:"status"`
	GoodsInfo         string      `json:"goods_info"`
	ReceiverName      string      `json:"receiver_name"`
	ReceiverPhone     string      `json:"receiver_phone"`
	ReceiverProvince  string      `json:"receiver_province"`
	ReceiverCity      string      `json:"receiver_city"`
	ReceiverDistrict  string      `json:"receiver_district"`
	ReceiverStreet    string      `json:"receiver_street"`
	ReceiverDetail    string      `json:"receiver_detail"`
	ReceiverLatitude  float64     `json:"receiver_latitude"`
	ReceiverLongitude float64     `json:"receiver_longitude"`
	CreateTime        string      `json:"create_time"`
	UpdateTime        string      `json:"update_time"`
	SenderAddress     AddressInfo `json:"sender_address"` // 寄件人地址信息
}

// AddressInfo 地址信息结构体（用于响应）
type AddressInfo struct {
	ID        uint    `json:"id"`
	UserID    uint    `json:"user_id"`
	Province  string  `json:"province"`
	City      string  `json:"city"`
	District  string  `json:"district"`
	Street    string  `json:"street"`
	Detail    string  `json:"detail"`
	Receiver  string  `json:"receiver"`
	Phone     string  `json:"phone"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	IsDefault int8    `json:"is_default"`
}

func (order Order) TableName() string {
	return "orders"
}

// DeliveryAssign 配送分配表结构体
type DeliveryAssign struct {
	ID             uint      `gorm:"primaryKey;column:id;autoIncrement" json:"id"`
	OrderID        uint      `gorm:"column:order_id;not null" json:"order_id"`
	DeliveryUserID uint      `gorm:"column:delivery_user_id;not null" json:"delivery_user_id"`
	AssignUserID   uint      `gorm:"column:assign_user_id;not null" json:"assign_user_id"`
	AssignTime     time.Time `gorm:"column:assign_time;not null;default:CURRENT_TIMESTAMP" json:"assign_time"`
}

func (da DeliveryAssign) TableName() string {
	return "delivery_assign"
}

// LocationTrace 位置轨迹表结构体
type LocationTrace struct {
	ID             uint      `gorm:"primaryKey;column:id;autoIncrement" json:"id"`
	OrderID        uint      `gorm:"column:order_id;not null" json:"order_id"`
	DeliveryUserID uint      `gorm:"column:delivery_user_id;not null" json:"delivery_user_id"`
	Longitude      float64   `gorm:"column:longitude;not null" json:"longitude"`
	Latitude       float64   `gorm:"column:latitude;not null" json:"latitude"`
	UploadTime     time.Time `gorm:"column:upload_time;not null;default:CURRENT_TIMESTAMP" json:"upload_time"`
}

func (lt LocationTrace) TableName() string {
	return "location_traces"
}
