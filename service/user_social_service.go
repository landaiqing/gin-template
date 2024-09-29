package service

import "schisandra-cloud-album/model"

type UserSocialService interface {
	QueryUserSocialByOpenIDService(openID string, source string) (model.ScaAuthUserSocial, error)
	AddUserSocialService(userSocial model.ScaAuthUserSocial) error
}
