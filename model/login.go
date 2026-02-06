package model

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	RoleID   int8   `json:"role_id"`
}
type User struct {
	Username string `gorm:"column:username" json:"username"`
	Password string `gorm:"column:password" json:"password"`
	Id       uint   `gorm:"primaryKey;column:id;autoIncrement" json:"id"`
	RoleID   int8   `gorm:"column:role_id" json:"role_id"`
}

func (u User) TableName() string {
	return "users"
}
