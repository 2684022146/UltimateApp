package model

// Order 订单结构体
type Order struct {
	ID             uint   `gorm:"primaryKey;column:id;autoIncrement" json:"id"`
	UserID         uint   `gorm:"column:user_id" json:"user_id"`
	DeliveryUserID uint   `gorm:"column:delivery_user_id" json:"delivery_user_id"`
	AddressID      uint   `gorm:"column:address_id" json:"address_id"`
	OrderNo        string `gorm:"column:order_no" json:"order_no"`
	Status         int8   `gorm:"column:status" json:"status"`
	GoodsInfo      string `gorm:"column:goods_info" json:"goods_info"`
	CreateTime     string `gorm:"column:create_time" json:"create_time"`
	UpdateTime     string `gorm:"column:update_time" json:"update_time"`
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	AddressID      uint   `json:"address_id" binding:"required"`
	DeliveryUserID uint   `json:"delivery_user_id" binding:"required"`
	GoodsInfo      string `json:"goods_info" binding:"required"`
}

// OrderResponse 订单响应结构体（包含地址信息）
type OrderResponse struct {
	ID             uint        `json:"id"`
	UserID         uint        `json:"user_id"`
	DeliveryUserID uint        `json:"delivery_user_id"`
	AddressID      uint        `json:"address_id"`
	OrderNo        string      `json:"order_no"`
	Status         int8        `json:"status"`
	GoodsInfo      string      `json:"goods_info"`
	CreateTime     string      `json:"create_time"`
	UpdateTime     string      `json:"update_time"`
	Address        AddressInfo `json:"address"` // 包含完整地址信息
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
