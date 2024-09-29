package impl

import (
	"schisandra-cloud-album/dao/impl"
	"schisandra-cloud-album/model"
)

var userSocialDao = impl.UserSocialImpl{}

type UserSocialServiceImpl struct{}

// QueryUserSocialByOpenIDService 查询用户第三方账号信息
func (UserSocialServiceImpl) QueryUserSocialByOpenIDService(openID string, source string) (model.ScaAuthUserSocial, error) {
	return userSocialDao.QueryUserSocialByOpenID(openID, source)
}

// AddUserSocialService 新增用户第三方账号信息
func (UserSocialServiceImpl) AddUserSocialService(userSocial model.ScaAuthUserSocial) error {
	return userSocialDao.AddUserSocial(userSocial)
}
