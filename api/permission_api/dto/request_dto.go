package dto

import (
	"schisandra-cloud-album/model"
)

// AddPermissionRequestDto 添加权限请求dto
type AddPermissionRequestDto struct {
	Permissions []model.ScaAuthPermission `form:"permissions[]" json:"permissions"`
}

// AddPermissionToRoleRequestDto 添加权限到角色请求dto
type AddPermissionToRoleRequestDto struct {
	RoleKey    string `json:"role_key"`
	Permission string `json:"permission"`
	Method     string `json:"method"`
}
