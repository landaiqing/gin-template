package dto

type RoleRequestDto struct {
	RoleName string `json:"role_name" binding:"required"`
	RoleKey  string `json:"role_key" binding:"required"`
}

type AddRoleToUserRequestDto struct {
	Uid     string `json:"uid" binding:"required"`
	RoleKey string `json:"role_key" binding:"required"`
}
