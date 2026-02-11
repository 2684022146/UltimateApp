package service

import (
	"fmt"
	"webdemo/repository"
)

type PermissionService interface {
	CheckPermission(roleId int8, apiPath, method string) bool
	RefreshPermissions() error
}
type permissionService struct {
	repo repository.PermissionRepository
	//[roleId][method][apiPath]bool
	permissionMap map[int8]map[string]map[string]bool
}

func NewPermissionService(repo repository.PermissionRepository) PermissionService {
	s := &permissionService{
		repo:          repo,
		permissionMap: make(map[int8]map[string]map[string]bool),
	}
	s.RefreshPermissions()
	return s
}

// 校验角色权限 在map里每层查找 O(1)
func (s *permissionService) CheckPermission(roleId int8, apiPath, method string) bool {
	methodMap, roleIdExists := s.permissionMap[roleId]
	if !roleIdExists {
		return false
	}
	apiPathMap, mothodExists := methodMap[method]
	if !mothodExists {
		return false
	}
	_, apiPathExists := apiPathMap[apiPath]
	return apiPathExists

}

// 将repo获取到的角色-权限切片放入map里 每次启动只调用一次
func (s *permissionService) RefreshPermissions() error {
	s.permissionMap = make(map[int8]map[string]map[string]bool)
	roles, err := s.repo.GetRoles()
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	for _, role := range roles {
		permissions, err := s.repo.GetPermissionsByRoleID(role)
		if err != nil {
			return fmt.Errorf("%s", err)
		}
		s.permissionMap[role] = make(map[string]map[string]bool)
		for _, permission := range permissions {
			method := permission.Method
			apiPath := permission.ApiPath
			if _, exists := s.permissionMap[role][method]; !exists {
				s.permissionMap[role][method] = make(map[string]bool)
			}
			s.permissionMap[role][method][apiPath] = true
		}
	}
	return nil
}
