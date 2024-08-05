package auth_service

import (
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/model"
)

func (AuthService) GetUserList() []*model.ScaAuthUser {
	data := make([]*model.ScaAuthUser, 0)
	global.DB.Find(&data)
	return data
}
