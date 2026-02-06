package model

type Address struct {
	UserID   int8   `gorm:"column:user_id" json:"user_id"`
	Province string `gorm:"column:province" json:"province"`
	City     string `gorm:"column:city" json:"city"`
	District string `gorm:"column:district" json:"district"`
	Street   string `gorm:"column:street" json:"street"`
	Detail   string `gorm:"column:detail" json:"detail"`
}
