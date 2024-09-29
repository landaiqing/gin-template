package impl

import (
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/model"
)

type PermissionDaoImpl struct{}

// GetPermissionsByIds 通过权限ID列表获取权限列表
func (PermissionDaoImpl) GetPermissionsByIds(ids []int64) ([]model.ScaAuthPermission, error) {
	var permissions []model.ScaAuthPermission
	if err := global.DB.Where("id IN ? and deleted = 0", ids).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// CreatePermissions 批量创建权限
func (PermissionDaoImpl) CreatePermissions(permissions []model.ScaAuthPermission) error {
	if err := global.DB.Model(&model.ScaAuthPermission{}).CreateInBatches(&permissions, len(permissions)).Error; err != nil {
		return err
	}
	return nil
}
