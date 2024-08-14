package role_permission_service

import (
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/model"
)

// QueryPermissionIdsByRoleId 通过角色ID列表查询权限ID列表
func (RolePermissionService) QueryPermissionIdsByRoleId(roleIds []*int64) []int64 {
	var permissionIds []int64
	rolePermission := model.ScaAuthRolePermission{}
	if err := global.DB.Model(&rolePermission).Where("role_id IN ?", roleIds).Pluck("permission_id", &permissionIds).Error; err != nil {
		global.LOG.Error(err)
		return nil
	}
	return permissionIds
}
