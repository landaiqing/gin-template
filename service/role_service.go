package service

import "schisandra-cloud-album/model"

type RoleService interface {
	// AddRoleService 添加角色
	AddRoleService(role model.ScaAuthRole) error
}
