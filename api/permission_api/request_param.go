package permission_api

import "schisandra-cloud-album/model"

// AddPermissionToRoleRequest 添加权限请求
type AddPermissionRequest struct {
	Permissions []model.ScaAuthPermission `form:"permissions[]" json:"permissions"`
}

// AddPermissionToRoleRequest 添加权限到角色请求
type AddPermissionToRoleRequest struct {
	RoleKey    string `json:"role_key"`
	Permission string `json:"permission"`
	Method     string `json:"method"`
}
type GetPermissionRequest struct {
	UserId string `json:"user_id" binding:"required"`
}
