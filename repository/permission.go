package repository

import (
	"errors"
	"webdemo/model"

	"gorm.io/gorm"
)

type PermissionRepository interface {
	GetPermissionsByRoleID(roleId int8) ([]*model.Permission, error)
	GetRoles() ([]int8, error)
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
	err := r.db.Model(&model.Permission{}).Joins("JOIN role_permissions ON permissions.id=role_permissions.perm_id").Where("role_permissions.role_id=?", roleId).Find(&permissions).Error
	if err != nil {
		return nil, errors.New("获取权限列表失败")
	}
	return permissions, nil
}
func (r *permissionRepository) GetRoles() ([]int8, error) {
	var roles = make([]int8, 0, 2)
	err := r.db.Table("roles").Select("id").Scan(&roles).Error
	if err != nil {
		return nil, errors.New("获取角色列表失败")
	}

	return roles, nil
}
