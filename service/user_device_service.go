package service

import "schisandra-cloud-album/model"

type UserDeviceService interface {
	GetUserDeviceByUIDIPAgentService(userId, ip, userAgent string) (*model.ScaAuthUserDevice, error)
	AddUserDeviceService(userDevice *model.ScaAuthUserDevice) error
	UpdateUserDeviceService(userId int64, userDevice *model.ScaAuthUserDevice) error
}
