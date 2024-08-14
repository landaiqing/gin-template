package role_service

import (
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/model"
)

// GetRoleListByIds : 通过Id列表获取角色信息列表
func (RoleService) GetRoleListByIds(id []*int64) ([]model.ScaAuthRole, error) {
	var roles []model.ScaAuthRole
	if err := global.DB.Where("id IN ?", id).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}
