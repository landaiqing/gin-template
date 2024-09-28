package role_api

import (
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/common/result"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/model"
	"schisandra-cloud-album/service"
)

var roleService = service.Service.RoleService

// CreateRole 创建角色
// @Summary 创建角色
// @Description 创建角色
// @Tags 角色
// @Accept  json
// @Produce  json
// @Param roleRequestDto body RoleRequest true "角色信息"
// @Router /api/auth/role/create [post]
func (RoleAPI) CreateRole(c *gin.Context) {
	roleRequest := RoleRequest{}
	err := c.ShouldBindJSON(&roleRequest)
	if err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CreatedFailed"), c)
		return
	}
	role := model.ScaAuthRole{
		RoleName: roleRequest.RoleName,
		RoleKey:  roleRequest.RoleKey,
	}
	err = roleService.AddRole(role)
	if err != nil {
		global.LOG.Error(err)
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CreatedFailed"), c)
		return
	}
	result.OkWithMessage(ginI18n.MustGetMessage(c, "CreatedSuccess"), c)
}

// AddRoleToUser 给指定用户添加角色
// @Summary 给指定用户添加角色
// @Description 给指定用户添加角色
// @Tags 角色
// @Accept  json
// @Produce  json
// @Param addRoleToUserRequestDto body AddRoleToUserRequest true "给指定用户添加角色"
// @Router /api/auth/role/add_role_to_user [post]
func (RoleAPI) AddRoleToUser(c *gin.Context) {
	addRoleToUserRequest := AddRoleToUserRequest{}
	err := c.ShouldBindJSON(&addRoleToUserRequest)
	if err != nil {
		global.LOG.Error(err)
		return
	}
	user, err := global.Casbin.AddRoleForUser(addRoleToUserRequest.Uid, addRoleToUserRequest.RoleKey)
	if err != nil {
		global.LOG.Error(err)
		return
	}
	result.OkWithData(user, c)
}
