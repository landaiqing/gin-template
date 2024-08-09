package auth_service

import (
	"gorm.io/gorm"
	"schisandra-cloud-album/common/enum"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/model"
)

// GetUserList 获取所有用户列表
func (AuthService) GetUserList() []*model.ScaAuthUser {
	data := make([]*model.ScaAuthUser, 0)
	global.DB.Where("deleted = 0 ").Find(&data)
	return data
}

// QueryUserByUsername 根据用户名查询用户
func (AuthService) QueryUserByUsername(username string) model.ScaAuthUser {
	authUser := model.ScaAuthUser{}
	global.DB.Where("username = ? and deleted = 0", username).First(&authUser)
	return authUser
}

// QueryUserByUuid 根据用户uuid查询用户
func (AuthService) QueryUserByUuid(uuid string) model.ScaAuthUser {
	authUser := model.ScaAuthUser{}
	global.DB.Where("uuid = ? and deleted = 0", uuid).First(&authUser)
	return authUser
}

// AddUser 添加用户
func (AuthService) AddUser(user model.ScaAuthUser) error {
	return global.DB.Create(&user).Error
}

// UpdateUser 更新用户
func (AuthService) UpdateUser(user model.ScaAuthUser) *gorm.DB {
	authUser := model.ScaAuthUser{}
	return global.DB.Model(&authUser).Where("uuid = ?", user.UUID).Updates(user)
}

// DeleteUser 删除用户
func (AuthService) DeleteUser(uuid string) error {
	authUser := model.ScaAuthUser{}
	return global.DB.Model(&authUser).Where("uuid = ?", uuid).Updates(&model.ScaAuthUser{Deleted: &enum.DELETED}).Error
}

// QueryUserByPhone 根据手机号查询用户
func (AuthService) QueryUserByPhone(phone string) model.ScaAuthUser {
	authUser := model.ScaAuthUser{}
	global.DB.Where("phone = ? and deleted = 0", phone).First(&authUser)
	return authUser
}

// QueryUserByEmail 根据邮箱查询用户
func (AuthService) QueryUserByEmail(email string) model.ScaAuthUser {
	authUser := model.ScaAuthUser{}
	global.DB.Where("email = ? and deleted = 0", email).First(&authUser)
	return authUser
}
