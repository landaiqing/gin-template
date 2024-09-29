package role_controller

type RoleRequest struct {
	RoleName string `json:"role_name" binding:"required"`
	RoleKey  string `json:"role_key" binding:"required"`
}

type AddRoleToUserRequest struct {
	Uid     string `json:"uid" binding:"required"`
	RoleKey string `json:"role_key" binding:"required"`
}
