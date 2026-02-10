package model

type Permission struct {
	ID         int    `gorm:"column:id"`
	PermCode   string `gorm:"column:perm_code"`
	PermName   string `gorm:"column:perm_name" json:"perm_name"`
	ApiPath    string `gorm:"column:api_path" json:"api_path"`
	Method     string `gorm:"column:method" json:"method"`
	CreateTime string `gorm:"column:create_time" json:"create_time"`
}

func (permission Permission) TableName() string {
	return "permissions"
}

// RolePermission 角色-权限关联结构体
type RolePermission struct {
	ID         uint   `gorm:"primaryKey;column:id;autoIncrement" json:"id"`
	RoleID     int8   `gorm:"column:role_id" json:"role_id"`
	PermID     uint   `gorm:"column:perm_id" json:"perm_id"`
	CreateTime string `gorm:"column:create_time" json:"create_time"`
}

// TableName 设置表名
func (RolePermission) TableName() string {
	return "role_permissions"
}
