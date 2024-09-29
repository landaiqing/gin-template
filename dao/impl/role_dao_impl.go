package impl

import (
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/model"
)

type RoleDaoImpl struct{}

// GetRoleListByIds : 通过Id列表获取角色信息列表
func (RoleDaoImpl) GetRoleListByIds(id []*int64) ([]model.ScaAuthRole, error) {
	var roles []model.ScaAuthRole
	if err := global.DB.Where("id IN ?", id).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// AddRole 新增角色
func (RoleDaoImpl) AddRole(role model.ScaAuthRole) error {
	if err := global.DB.Create(&role).Error; err != nil {
		return err
	}
	return nil
}
