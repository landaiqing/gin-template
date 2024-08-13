package user_role_service

import (
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/model"
)

// GetUserRoleIdsByUserId 通过用户ID获取用户角色ID列表
func (UserRoleService) GetUserRoleIdsByUserId(userId int64) ([]*int64, error) {
	var roleIds []*int64
	if err := global.DB.Table("sca_auth_user_role").Where("user_id = ?", userId).Pluck("role_id", &roleIds).Error; err != nil {
		return nil, err
	}
	return roleIds, nil
}

// AddUserRole 新增用户角色
func (UserRoleService) AddUserRole(userRole model.ScaAuthUserRole) error {
	if err := global.DB.Create(&userRole).Error; err != nil {
		return err
	}
	return nil
}
