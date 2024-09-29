package impl

import (
	"schisandra-cloud-album/dao/impl"
	"schisandra-cloud-album/model"
)

var userDao = impl.UserDaoImpl{}

type UserServiceImpl struct{}

// GetUserListService 返回用户列表
func (UserServiceImpl) GetUserListService() []*model.ScaAuthUser {
	return userDao.GetUserList()
}

// QueryUserByUsernameService 根据用户名查询用户
func (UserServiceImpl) QueryUserByUsernameService(username string) model.ScaAuthUser {
	return userDao.QueryUserByUsername(username)
}

// QueryUserByUuidService 根据uid查询用户
func (UserServiceImpl) QueryUserByUuidService(uid *string) model.ScaAuthUser {
	user, err := userDao.QueryUserByUuid(uid)
	if err != nil {
		return model.ScaAuthUser{}
	}
	return user
}

// DeleteUserService 根据uid删除用户
func (UserServiceImpl) DeleteUserService(uid string) error {
	return userDao.DeleteUser(uid)
}

// QueryUserByPhoneService 根据手机号查询用户
func (UserServiceImpl) QueryUserByPhoneService(phone string) model.ScaAuthUser {
	return userDao.QueryUserByPhone(phone)
}

// QueryUserByEmailService 根据邮箱查询用户
func (UserServiceImpl) QueryUserByEmailService(email string) model.ScaAuthUser {
	return userDao.QueryUserByEmail(email)
}

// AddUserService 新增用户
func (UserServiceImpl) AddUserService(user model.ScaAuthUser) (*model.ScaAuthUser, error) {
	return userDao.AddUser(user)
}

// UpdateUserService 更新用户信息
func (UserServiceImpl) UpdateUserService(phone, encrypt string) error {
	return userDao.UpdateUser(phone, encrypt)
}
