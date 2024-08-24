package user_service

import (
	"schisandra-cloud-album/common/enum"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/model"
)

// GetUserList 获取所有用户列表
func (UserService) GetUserList() []*model.ScaAuthUser {
	data := make([]*model.ScaAuthUser, 0)
	global.DB.Where("deleted = 0 ").Find(&data)
	return data
}

// QueryUserByUsername 根据用户名查询用户
func (UserService) QueryUserByUsername(username string) model.ScaAuthUser {
	authUser := model.ScaAuthUser{}
	err := global.DB.Where("username = ? and deleted = 0", username).First(&authUser).Error
	if err != nil {
		return model.ScaAuthUser{}
	}
	return authUser
}

// QueryUserByUuid 根据用户uuid查询用户
func (UserService) QueryUserByUuid(uuid *string) (model.ScaAuthUser, error) {
	authUser := model.ScaAuthUser{}
	if err := global.DB.Where("uid = ? and deleted = 0", uuid).First(&authUser).Error; err != nil {
		return model.ScaAuthUser{}, err
	}
	return authUser, nil
}

// QueryUserById 根据用户id查询用户
func (UserService) QueryUserById(id *int64) (model.ScaAuthUser, error) {
	authUser := model.ScaAuthUser{}
	if err := global.DB.Where("id = ? and deleted = 0", id).First(&authUser).Error; err != nil {
		return model.ScaAuthUser{}, err
	}
	return authUser, nil
}

// AddUser 添加用户
func (UserService) AddUser(user model.ScaAuthUser) (model.ScaAuthUser, error) {
	if err := global.DB.Create(&user).Error; err != nil {
		return model.ScaAuthUser{}, err
	}
	// 查询创建后的用户信息
	var createdUser model.ScaAuthUser
	if err := global.DB.First(&createdUser, user.ID).Error; err != nil {
		return model.ScaAuthUser{}, err
	}
	return createdUser, nil
}

// UpdateUser 更新用户
func (UserService) UpdateUser(phone string, password string) error {
	return global.DB.Model(&model.ScaAuthUser{}).Where("phone = ? and deleted = 0", phone).Updates(&model.ScaAuthUser{Password: &password}).Error
}

// DeleteUser 删除用户
func (UserService) DeleteUser(uuid string) error {
	authUser := model.ScaAuthUser{}
	return global.DB.Model(&authUser).Where("uid = ?", uuid).Updates(&model.ScaAuthUser{Deleted: &enum.DELETED}).Error
}

// QueryUserByPhone 根据手机号查询用户
func (UserService) QueryUserByPhone(phone string) model.ScaAuthUser {
	authUser := model.ScaAuthUser{}
	global.DB.Where("phone = ? and deleted = 0", phone).First(&authUser)
	return authUser
}

// QueryUserByEmail 根据邮箱查询用户
func (UserService) QueryUserByEmail(email string) model.ScaAuthUser {
	authUser := model.ScaAuthUser{}
	global.DB.Where("email = ? and deleted = 0", email).First(&authUser)
	return authUser
}
