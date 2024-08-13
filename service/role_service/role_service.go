package role_service

import (
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/model"
)

// GetRoleById : 通过Id获取角色信息
func (RoleService) GetRoleById(id int64) (model.ScaAuthRole, error) {
	var role model.ScaAuthRole
	if err := global.DB.Where("id = ? and deleted = 0", id).First(&role).Error; err != nil {
		return model.ScaAuthRole{}, err
	}
	return role, nil
}
