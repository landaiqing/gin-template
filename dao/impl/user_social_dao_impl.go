package impl

import (
	"errors"
	"gorm.io/gorm"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/model"
)

type UserSocialImpl struct{}

// AddUserSocial 添加社会化登录用户信息
func (UserSocialImpl) AddUserSocial(user model.ScaAuthUserSocial) error {
	result := global.DB.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// QueryUserSocialByOpenID 根据openID和source查询用户信息
func (UserSocialImpl) QueryUserSocialByOpenID(openID string, source string) (model.ScaAuthUserSocial, error) {
	var user model.ScaAuthUserSocial
	result := global.DB.Where("open_id = ? and source = ? and deleted = 0", openID, source).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return model.ScaAuthUserSocial{}, result.Error
		}
		return model.ScaAuthUserSocial{}, result.Error
	}
	return user, nil
}
