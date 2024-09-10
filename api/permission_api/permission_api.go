package permission_api

import (
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/api/permission_api/dto"
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
// @Param permissions body dto.AddPermissionRequestDto true "权限列表"
// @Router /api/auth/permission/add [post]
func (PermissionAPI) AddPermissions(c *gin.Context) {
	addPermissionRequestDto := dto.AddPermissionRequestDto{}
	err := c.ShouldBind(&addPermissionRequestDto.Permissions)
	if err != nil {
		global.LOG.Error(err)
		return
	}
	err = permissionService.CreatePermissions(addPermissionRequestDto.Permissions)
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
// @Param permissions body dto.AddPermissionToRoleRequestDto true "权限列表"
// @Router /api/auth/permission/assign [post]
func (PermissionAPI) AssignPermissionsToRole(c *gin.Context) {
	permissionToRoleRequestDto := dto.AddPermissionToRoleRequestDto{}

	err := c.ShouldBind(&permissionToRoleRequestDto)

	if err != nil {
		global.LOG.Error(err)
		result.FailWithMessage(ginI18n.MustGetMessage(c, "AssignFailed"), c)
		return
	}

	policy, err := global.Casbin.AddPolicy(permissionToRoleRequestDto.RoleKey, permissionToRoleRequestDto.Permission, permissionToRoleRequestDto.Method)
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
	userId := c.Query("user_id")
	if userId == "" {
		result.FailWithMessage("user_id is required", c)
		return
	}
	data, err := global.Casbin.GetImplicitRolesForUser(userId)
	if err != nil {
		result.FailWithMessage("Get user permissions failed", c)
		return
	}
	result.OkWithData(data, c)
	return
}
