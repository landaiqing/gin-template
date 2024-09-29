package user_controller

import (
	"encoding/gob"
	"errors"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/mssola/useragent"
	"gorm.io/gorm"
	"schisandra-cloud-album/common/constant"
	"schisandra-cloud-album/common/redis"
	"schisandra-cloud-album/common/result"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/model"
	"schisandra-cloud-album/utils"
	"time"
)

// getUserLoginDevice 获取用户登录设备
func getUserLoginDevice(user model.ScaAuthUser, c *gin.Context) bool {

	// 检查user.UID是否为空
	if user.UID == nil {
		global.LOG.Errorln("user.UID is nil")
		return false
	}
	userAgent := c.GetHeader("User-Agent")
	if userAgent == "" {
		global.LOG.Errorln("user-agent is empty")
		return false
	}
	ua := useragent.New(userAgent)

	ip := utils.GetClientIP(c)
	location, err := global.IP2Location.SearchByStr(ip)
	if err != nil {
		global.LOG.Errorln(err)
		return false
	}
	location = utils.RemoveZeroAndAdjust(location)

	isBot := ua.Bot()
	browser, browserVersion := ua.Browser()
	os := ua.OS()
	mobile := ua.Mobile()
	mozilla := ua.Mozilla()
	platform := ua.Platform()
	engine, engineVersion := ua.Engine()

	device := model.ScaAuthUserDevice{
		UserID:          user.UID,
		IP:              &ip,
		Location:        &location,
		Agent:           userAgent,
		Browser:         &browser,
		BrowserVersion:  &browserVersion,
		OperatingSystem: &os,
		Mobile:          &mobile,
		Bot:             &isBot,
		Mozilla:         &mozilla,
		Platform:        &platform,
		EngineName:      &engine,
		EngineVersion:   &engineVersion,
	}

	mu.Lock()
	defer mu.Unlock()

	userDevice, err := userDeviceService.GetUserDeviceByUIDIPAgentService(*user.UID, ip, userAgent)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = userDeviceService.AddUserDeviceService(&device)
		if err != nil {
			global.LOG.Errorln(err)
			return false
		}
	} else if err != nil {
		global.LOG.Errorln(err)
		return false
	} else {
		err := userDeviceService.UpdateUserDeviceService(userDevice.ID, &device)
		if err != nil {
			global.LOG.Errorln(err)
			return false
		}
	}

	return true
}

// handelUserLogin 处理用户登录
func handelUserLogin(user model.ScaAuthUser, autoLogin bool, c *gin.Context) {
	// 检查 user.UID 是否为 nil
	if user.UID == nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ParamsError"), c)
		return
	}
	if !getUserLoginDevice(user, c) {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "LoginFailed"), c)
		return
	}
	accessToken, err := utils.GenerateAccessToken(utils.AccessJWTPayload{UserID: user.UID})
	if err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "LoginFailed"), c)
		return
	}

	var days time.Duration
	if autoLogin {
		days = 7 * 24 * time.Hour
	} else {
		days = time.Minute * 30
	}

	refreshToken, expiresAt := utils.GenerateRefreshToken(utils.RefreshJWTPayload{UserID: user.UID}, days)
	data := ResponseData{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		UID:          user.UID,
	}

	err = redis.Set(constant.UserLoginTokenRedisKey+*user.UID, data, days).Err()
	if err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "LoginFailed"), c)
		return
	}
	gob.Register(ResponseData{})
	err = utils.SetSession(c, constant.SessionKey, data)
	if err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "LoginFailed"), c)
		return
	}
	result.OkWithData(data, c)
}
