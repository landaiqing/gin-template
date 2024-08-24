package dto

type RoleRequestDto struct {
	RoleName string `json:"role_name"`
	RoleKey  string `json:"role_key"`
}

type AddRoleToUserRequestDto struct {
	Uid     string `json:"uid"`
	RoleKey string `json:"role_key"`
}
