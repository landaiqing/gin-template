package service

import (
	"schisandra-cloud-album/service/permission_service"
	"schisandra-cloud-album/service/role_permission_service"
	"schisandra-cloud-album/service/role_service"
	"schisandra-cloud-album/service/user_role_service"
	"schisandra-cloud-album/service/user_service"
	"schisandra-cloud-album/service/user_social_service"
)

// Services 统一导出的service
type Services struct {
	UserService           user_service.UserService
	RoleService           role_service.RoleService
	UserRoleService       user_role_service.UserRoleService
	RolePermissionService role_permission_service.RolePermissionService
	PermissionService     permission_service.PermissionService
	UserSocialService     user_social_service.UserSocialService
}

// Service new函数实例化，实例化完成后会返回结构体地指针类型
var Service = new(Services)
