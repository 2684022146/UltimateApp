package model

// Order 订单结构体
type Order struct {
	ID                uint    `gorm:"primaryKey;column:id;autoIncrement" json:"id"`
	SenderUserID      uint    `gorm:"column:sender_user_id" json:"sender_user_id"`
	SenderAddressID   uint    `gorm:"column:sender_address_id" json:"sender_address_id"`
	DeliveryUserID    uint    `gorm:"column:delivery_user_id" json:"delivery_user_id"`
	OrderNo           string  `gorm:"column:order_no" json:"order_no"`
	Status            int8    `gorm:"column:status" json:"status"`
	GoodsInfo         string  `gorm:"column:goods_info" json:"goods_info"`
	ReceiverName      string  `gorm:"column:receiver_name" json:"receiver_name"`
	ReceiverPhone     string  `gorm:"column:receiver_phone" json:"receiver_phone"`
	ReceiverProvince  string  `gorm:"column:receiver_province" json:"receiver_province"`
	ReceiverCity      string  `gorm:"column:receiver_city" json:"receiver_city"`
	ReceiverDistrict  string  `gorm:"column:receiver_district" json:"receiver_district"`
	ReceiverStreet    string  `gorm:"column:receiver_street" json:"receiver_street"`
	ReceiverDetail    string  `gorm:"column:receiver_detail" json:"receiver_detail"`
	ReceiverLatitude  float64 `gorm:"column:receiver_latitude" json:"receiver_latitude"`
	ReceiverLongitude float64 `gorm:"column:receiver_longitude" json:"receiver_longitude"`
	CreateTime        string  `gorm:"column:create_time" json:"create_time"`
	UpdateTime        string  `gorm:"column:update_time" json:"update_time"`
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	SenderAddressID   uint    `json:"sender_address_id" binding:"required"`
	DeliveryUserID    uint    `json:"delivery_user_id" binding:"required"`
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
	DeliveryUserID    uint        `json:"delivery_user_id"`
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
