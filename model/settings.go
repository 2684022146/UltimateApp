package model

// Address 收货地址结构体
type Address struct {
	ID        uint    `gorm:"primaryKey;column:id;autoIncrement" json:"id"`
	UserID    uint    `gorm:"column:user_id" json:"user_id"`
	Province  string  `gorm:"column:province" json:"province"`
	City      string  `gorm:"column:city" json:"city"`
	District  string  `gorm:"column:district" json:"district"`
	Street    string  `gorm:"column:street" json:"street"`
	Detail    string  `gorm:"column:detail" json:"detail"`
	Receiver  string  `gorm:"column:receiver" json:"receiver"`
	Phone     string  `gorm:"column:phone" json:"phone"`
	Latitude  float64 `gorm:"column:latitude" json:"latitude"`
	Longitude float64 `gorm:"column:longitude" json:"longitude"`
	IsDefault int8    `gorm:"column:is_default" json:"is_default"`
}

// CreateAddressRequest 创建地址请求
type CreateAddressRequest struct {
	Province  string  `json:"province" binding:"required"`
	City      string  `json:"city" binding:"required"`
	District  string  `json:"district"`
	Street    string  `json:"street"`
	Detail    string  `json:"detail" binding:"required"`
	Receiver  string  `json:"receiver" binding:"required"`
	Phone     string  `json:"phone" binding:"required"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	IsDefault int8    `json:"is_default"`
}

// UpdateAddressRequest 更新地址请求
type UpdateAddressRequest struct {
	ID        uint    `json:"id" binding:"required"`
	Province  string  `json:"province" binding:"required"`
	City      string  `json:"city" binding:"required"`
	District  string  `json:"district"`
	Street    string  `json:"street"`
	Detail    string  `json:"detail" binding:"required"`
	Receiver  string  `json:"receiver" binding:"required"`
	Phone     string  `json:"phone" binding:"required"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	IsDefault int8    `json:"is_default"`
}

// DeleteAddressRequest 删除地址请求
type DeleteAddressRequest struct {
	ID uint `json:"id" binding:"required"`
}

func (address Address) TableName() string {
	return "addresses"
}
