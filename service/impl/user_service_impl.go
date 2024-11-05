package impl

import (
	"errors"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mssola/useragent"
	"gorm.io/gorm"

	"schisandra-cloud-album/common/constant"
	"schisandra-cloud-album/common/redis"
	"schisandra-cloud-album/common/types"
	"schisandra-cloud-album/dao/impl"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/model"
	"schisandra-cloud-album/utils"
)

var userDao = impl.UserDaoImpl{}

type UserServiceImpl struct{}

var mu = &sync.Mutex{}

// GetUserListService 返回用户列表
func (UserServiceImpl) GetUserListService() []*model.ScaAuthUser {
	return userDao.GetUserList()
}

// QueryUserByUsernameService 根据用户名查询用户
func (UserServiceImpl) QueryUserByUsernameService(username string) model.ScaAuthUser {
	return userDao.QueryUserByUsername(username)
}

// QueryUserByUuidService 根据uid查询用户
func (UserServiceImpl) QueryUserByUuidService(uid *string) model.ScaAuthUser {
	user, err := userDao.QueryUserByUuid(uid)
	if err != nil {
		return model.ScaAuthUser{}
	}
	return user
}

// DeleteUserService 根据uid删除用户
func (UserServiceImpl) DeleteUserService(uid string) error {
	return userDao.DeleteUser(uid)
}

// QueryUserByPhoneService 根据手机号查询用户
func (UserServiceImpl) QueryUserByPhoneService(phone string) model.ScaAuthUser {
	return userDao.QueryUserByPhone(phone)
}

// QueryUserByEmailService 根据邮箱查询用户
func (UserServiceImpl) QueryUserByEmailService(email string) model.ScaAuthUser {
	return userDao.QueryUserByEmail(email)
}

// AddUserService 新增用户
func (UserServiceImpl) AddUserService(user model.ScaAuthUser) (*model.ScaAuthUser, error) {
	return userDao.AddUser(user)
}

// UpdateUserService 更新用户信息
func (UserServiceImpl) UpdateUserService(phone, encrypt string) error {
	return userDao.UpdateUser(phone, encrypt)
}

// RefreshTokenService 刷新用户token
func (UserServiceImpl) RefreshTokenService(c *gin.Context, refreshToken string) (string, bool) {
	parseRefreshToken, isUpd, err := utils.ParseRefreshToken(refreshToken)
	if err != nil || !isUpd {
		global.LOG.Errorln(err)
		return "", false
	}
	accessTokenString, err := utils.GenerateAccessToken(utils.AccessJWTPayload{UserID: parseRefreshToken.UserID})
	if err != nil {
		return "", false
	}
	tokenKey := constant.UserLoginTokenRedisKey + *parseRefreshToken.UserID
	session := utils.GetSession(c, constant.SessionKey)
	if session.RefreshToken == "" {
		return "", false
	}
	redisTokenData := types.RedisToken{
		AccessToken: accessTokenString,
		UID:         *parseRefreshToken.UserID,
	}
	if err = redis.Set(tokenKey, redisTokenData, time.Hour*24*7).Err(); err != nil {
		global.LOG.Errorln(err)
		return "", false
	}
	return accessTokenString, true
}

// HandelUserLogin 处理用户登录
func (UserServiceImpl) HandelUserLogin(user model.ScaAuthUser, autoLogin bool, c *gin.Context) (*types.ResponseData, bool) {
	// 检查 user.UID 是否为 nil
	if user.UID == "" {
		return nil, false
	}
	if !GetUserLoginDevice(user, c) {
		return nil, false
	}
	accessToken, err := utils.GenerateAccessToken(utils.AccessJWTPayload{UserID: &user.UID})
	if err != nil {
		return nil, false
	}
	var days time.Duration
	if autoLogin {
		days = 7 * 24 * time.Hour
	} else {
		days = time.Minute * 30
	}

	refreshToken := utils.GenerateRefreshToken(utils.RefreshJWTPayload{UserID: &user.UID}, days)
	data := types.ResponseData{
		AccessToken: accessToken,
		UID:         &user.UID,
		Username:    user.Username,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Status:      user.Status,
	}
	redisTokenData := types.RedisToken{
		AccessToken: accessToken,
		UID:         user.UID,
	}
	err = redis.Set(constant.UserLoginTokenRedisKey+user.UID, redisTokenData, days).Err()
	if err != nil {
		return nil, false
	}
	sessionData := utils.SessionData{
		RefreshToken: refreshToken,
		UID:          user.UID,
	}
	err = utils.SetSession(c, constant.SessionKey, sessionData)
	if err != nil {
		return nil, false
	}
	return &data, true
}

// GetUserLoginDevice 获取用户登录设备
func GetUserLoginDevice(user model.ScaAuthUser, c *gin.Context) bool {

	// 检查user.UID是否为空
	if user.UID == "" {
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
		IP:              ip,
		Location:        location,
		Agent:           userAgent,
		Browser:         browser,
		BrowserVersion:  browserVersion,
		OperatingSystem: os,
		Mobile:          mobile,
		Bot:             isBot,
		Mozilla:         mozilla,
		Platform:        platform,
		EngineName:      engine,
		EngineVersion:   engineVersion,
	}

	mu.Lock()
	defer mu.Unlock()

	userDevice, err := userDeviceDao.GetUserDeviceByUIDIPAgent(user.UID, ip, userAgent)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = userDeviceDao.AddUserDevice(&device)
		if err != nil {
			global.LOG.Errorln(err)
			return false
		}
	} else if err != nil {
		global.LOG.Errorln(err)
		return false
	} else {
		err := userDeviceDao.UpdateUserDevice(userDevice.ID, &device)
		if err != nil {
			global.LOG.Errorln(err)
			return false
		}
	}
	return true
}
