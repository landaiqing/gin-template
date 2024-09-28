package permission_api

import (
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/common/result"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/service"
)

var permissionService = service.Service.PermissionService

// AddPermissions 批量添加权限
// @Summary 批量添加权限
// @Description 批量添加权限
// @Tags 权限管理
// @Accept  json
// @Produce  json
// @Param permissions body AddPermissionRequest true "权限列表"
// @Router /api/auth/permission/add [post]
func (PermissionAPI) AddPermissions(c *gin.Context) {
	addPermissionRequest := AddPermissionRequest{}
	err := c.ShouldBind(&addPermissionRequest.Permissions)
	if err != nil {
		global.LOG.Error(err)
		return
	}
	err = permissionService.CreatePermissions(addPermissionRequest.Permissions)
	if err != nil {
		global.LOG.Error(err)
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CreatedFailed"), c)
		return
	}
	result.OkWithMessage(ginI18n.MustGetMessage(c, "CreatedSuccess"), c)
	return
}

// AssignPermissionsToRole 给指定角色分配权限
// @Summary 给指定角色分配权限
// @Description 给指定角色分配权限
// @Tags 权限管理
// @Accept  json
// @Produce  json
// @Param permissions body AddPermissionToRoleRequest true "权限列表"
// @Router /api/auth/permission/assign [post]
func (PermissionAPI) AssignPermissionsToRole(c *gin.Context) {
	permissionToRoleRequest := AddPermissionToRoleRequest{}

	err := c.ShouldBind(&permissionToRoleRequest)

	if err != nil {
		global.LOG.Error(err)
		result.FailWithMessage(ginI18n.MustGetMessage(c, "AssignFailed"), c)
		return
	}

	policy, err := global.Casbin.AddPolicy(permissionToRoleRequest.RoleKey, permissionToRoleRequest.Permission, permissionToRoleRequest.Method)
	if err != nil {
		global.LOG.Error(err)
		result.FailWithMessage(ginI18n.MustGetMessage(c, "AssignFailed"), c)
		return
	}
	if policy == false {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "AssignFailed"), c)
		return
	}
	result.OkWithMessage(ginI18n.MustGetMessage(c, "AssignSuccess"), c)
	return
}

// GetUserPermissions 获取用户角色权限
func (PermissionAPI) GetUserPermissions(c *gin.Context) {
	getPermissionRequest := GetPermissionRequest{}
	err := c.ShouldBindJSON(&getPermissionRequest)
	if err != nil {
		global.LOG.Error(err)
		return
	}
	data, err := global.Casbin.GetImplicitRolesForUser(getPermissionRequest.UserId)
	if err != nil {
		result.FailWithMessage("Get user permissions failed", c)
		return
	}
	result.OkWithData(data, c)
	return
}
