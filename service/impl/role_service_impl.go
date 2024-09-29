package impl

import (
	"schisandra-cloud-album/dao/impl"
	"schisandra-cloud-album/model"
)

var roleDao = impl.RoleDaoImpl{}

type RoleServiceImpl struct{}

// AddRoleService 添加角色
func (RoleServiceImpl) AddRoleService(role model.ScaAuthRole) error {
	return roleDao.AddRole(role)

}
