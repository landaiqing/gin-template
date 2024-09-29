package dao

import "schisandra-cloud-album/model"

type UserSocialDao interface {
	// AddUserSocial 添加用户第三方登录信息
	AddUserSocial(user model.ScaAuthUserSocial) error
	// QueryUserSocialByOpenID 根据第三方登录的 openID 查询用户信息
	QueryUserSocialByOpenID(openID string, source string) (model.ScaAuthUserSocial, error)
}
