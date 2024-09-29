package impl

import (
	"schisandra-cloud-album/dao/impl"
	"schisandra-cloud-album/model"
)

var userDeviceDao = impl.UserDeviceImpl{}

type UserDeviceServiceImpl struct{}

// GetUserDeviceByUIDIPAgentService 获取用户设备信息 根据用户ID、IP、User-Agent
func (UserDeviceServiceImpl) GetUserDeviceByUIDIPAgentService(userId, ip, userAgent string) (*model.ScaAuthUserDevice, error) {
	return userDeviceDao.GetUserDeviceByUIDIPAgent(userId, ip, userAgent)
}

// AddUserDeviceService 新增用户设备信息
func (UserDeviceServiceImpl) AddUserDeviceService(userDevice *model.ScaAuthUserDevice) error {
	return userDeviceDao.AddUserDevice(userDevice)
}

// UpdateUserDeviceService 更新用户设备信息
func (UserDeviceServiceImpl) UpdateUserDeviceService(userId int64, userDevice *model.ScaAuthUserDevice) error {
	return userDeviceDao.UpdateUserDevice(userId, userDevice)
}
