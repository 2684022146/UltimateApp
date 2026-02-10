package repository

import (
	"fmt"
	"webdemo/model"

	"gorm.io/gorm"
)

type PermissionRepository interface {
	GetPermissionsByRoleID(roleId int8) ([]*model.Permission, error)
}
type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{
		db: db,
	}
}
func (r *permissionRepository) GetPermissionsByRoleID(roleId int8) ([]*model.Permission, error) {
	var permissions []*model.Permission
	err := r.db.Table("permissions").Model(&model.Permission{}).Joins("JOIN role_permissions ON permissions.id=role_permissions.perm_id").Where("role_permission.role_id=?", roleId).Find(&permissions).Error
	if err != nil {
		return nil, fmt.Errorf("获取权限校验失败")
	}
	return permissions, nil
}
