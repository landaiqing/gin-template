package user_device_service

import (
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/model"
)

// AddUserDevice 新增用户设备信息
func (UserDeviceService) AddUserDevice(userDevice *model.ScaAuthUserDevice) error {
	if err := global.DB.Create(&userDevice).Error; err != nil {
		return err
	}
	return nil
}

// GetUserDeviceByUIDIPAgent 根据uid / IP / agent 查询用户设备信息
func (UserDeviceService) GetUserDeviceByUIDIPAgent(uid, ip, agent string) (*model.ScaAuthUserDevice, error) {
	var userDevice model.ScaAuthUserDevice
	if err := global.DB.Where("user_id =? AND ip =? AND agent =? AND deleted = 0 ", uid, ip, agent).First(&userDevice).Error; err != nil {
		return nil, err
	}
	return &userDevice, nil
}

// UpdateUserDevice 更新用户设备信息
func (UserDeviceService) UpdateUserDevice(id int64, userDevice *model.ScaAuthUserDevice) error {
	result := global.DB.Model(&userDevice).Where("id =? AND deleted = 0 ", id).Updates(model.ScaAuthUserDevice{
		IP:              userDevice.IP,
		Location:        userDevice.Location,
		Agent:           userDevice.Agent,
		Browser:         userDevice.Browser,
		BrowserVersion:  userDevice.BrowserVersion,
		OperatingSystem: userDevice.OperatingSystem,
		Mobile:          userDevice.Mobile,
		Bot:             userDevice.Bot,
		Mozilla:         userDevice.Mozilla,
		Model:           userDevice.Model,
		Platform:        userDevice.Platform,
		EngineName:      userDevice.EngineName,
		EngineVersion:   userDevice.EngineVersion,
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
