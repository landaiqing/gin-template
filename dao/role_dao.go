package dao

import "schisandra-cloud-album/model"

type RoleDAO interface {
	// GetRoleListByIds 获取角色列表 by id
	GetRoleListByIds(id []*int64) ([]model.ScaAuthRole, error)
	// AddRole 新增角色
	AddRole(role model.ScaAuthRole) error
}
