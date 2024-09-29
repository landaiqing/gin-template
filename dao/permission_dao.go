package dao

import "schisandra-cloud-album/model"

type PermissionDao interface {
	// GetPermissionsByIds 返回权限列表 根据权限ID列表
	GetPermissionsByIds(ids []int64) ([]model.ScaAuthPermission, error)
	// CreatePermissions 创建权限
	CreatePermissions(permissions []model.ScaAuthPermission) error
}
