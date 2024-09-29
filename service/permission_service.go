package service

import "schisandra-cloud-album/model"

type PermissionService interface {
	// CreatePermissionsService 创建权限
	CreatePermissionsService(permissions []model.ScaAuthPermission) error
}
