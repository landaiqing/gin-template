package impl

import (
	"schisandra-cloud-album/dao/impl"
	"schisandra-cloud-album/model"
)

var permissionDao = impl.PermissionDaoImpl{}

type PermissionServiceImpl struct{}

// CreatePermissionsService 创建权限
func (PermissionServiceImpl) CreatePermissionsService(permissions []model.ScaAuthPermission) error {
	return permissionDao.CreatePermissions(permissions)
}
