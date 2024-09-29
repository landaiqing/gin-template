package service

import "schisandra-cloud-album/model"

type UserService interface {
	// GetUserListService 返回用户列表
	GetUserListService() []*model.ScaAuthUser
	// QueryUserByUsernameService 根据用户名查询用户
	QueryUserByUsernameService(username string) model.ScaAuthUser
	// QueryUserByUuidService 根据用户ID查询用户
	QueryUserByUuidService(uid *string) model.ScaAuthUser
	// DeleteUserService 根据用户ID删除用户
	DeleteUserService(uid string) error
	// QueryUserByPhoneService 根据手机号查询用户
	QueryUserByPhoneService(phone string) model.ScaAuthUser
	// QueryUserByEmailService 根据邮箱查询用户
	QueryUserByEmailService(email string) model.ScaAuthUser
	// AddUserService 新增用户
	AddUserService(user model.ScaAuthUser) (*model.ScaAuthUser, error)
	// UpdateUserService 更新用户信息
	UpdateUserService(phone, encrypt string) error
}
