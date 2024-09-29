package dao

import "schisandra-cloud-album/model"

type UserDao interface {
	// GetUserList 获取用户列表
	GetUserList() []*model.ScaAuthUser
	// QueryUserByUsername 根据用户名查询用户
	QueryUserByUsername(username string) model.ScaAuthUser
	// QueryUserByUuid 根据uuid查询用户
	QueryUserByUuid(uuid *string) (model.ScaAuthUser, error)
	// QueryUserById 根据id查询用户
	QueryUserById(id *int64) (model.ScaAuthUser, error)
	// AddUser 新增用户
	AddUser(user model.ScaAuthUser) (*model.ScaAuthUser, error)
	// UpdateUser 更新用户
	UpdateUser(phone string, password string) error
	// DeleteUser 删除用户
	DeleteUser(uuid string) error
	// QueryUserByPhone 根据手机号查询用户
	QueryUserByPhone(phone string) model.ScaAuthUser
	// QueryUserByEmail 根据邮箱查询用户
	QueryUserByEmail(email string) model.ScaAuthUser
}
