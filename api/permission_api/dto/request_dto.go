package dto

import "schisandra-cloud-album/model"

type AddPermissionRequestDto struct {
	Permissions []model.ScaAuthPermission `json:"permissions"`
}
type AddPermissionToRoleRequestDto struct {
	RoleKey     string                    `json:"role_key"`
	Permissions []model.ScaAuthPermission `json:"permissions"`
}
