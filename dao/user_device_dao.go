package dao

import "schisandra-cloud-album/model"

type UserDeviceDao interface {
	// AddUserDevice 添加用户设备信息
	AddUserDevice(userDevice *model.ScaAuthUserDevice) error
	// GetUserDeviceByUIDIPAgent 根据用户ID、IP、Agent获取用户设备信息
	GetUserDeviceByUIDIPAgent(uid, ip, agent string) (*model.ScaAuthUserDevice, error)
	// UpdateUserDevice 更新用户设备信息
	UpdateUserDevice(id int64, userDevice *model.ScaAuthUserDevice) error
}
