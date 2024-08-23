package permission_service

import (
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/model"
)

// GetPermissionsByIds 通过权限ID列表获取权限列表
func (PermissionService) GetPermissionsByIds(ids []int64) ([]model.ScaAuthPermission, error) {
	var permissions []model.ScaAuthPermission
	if err := global.DB.Where("id IN ? and deleted = 0", ids).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}
