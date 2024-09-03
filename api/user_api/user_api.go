package user_api

import (
	"errors"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/mssola/useragent"
	"github.com/yitter/idgenerator-go/idgen"
	"gorm.io/gorm"
	"reflect"
	"schisandra-cloud-album/api/user_api/dto"
	"schisandra-cloud-album/common/constant"
	"schisandra-cloud-album/common/enum"
	"schisandra-cloud-album/common/redis"
	"schisandra-cloud-album/common/result"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/model"
	"schisandra-cloud-album/service"
	"schisandra-cloud-album/utils"
	"strconv"
	"time"
)

var userService = service.Service.UserService
var userDeviceService = service.Service.UserDeviceService

// GetUserList
// @Summary 获取所有用户列表
// @Tags 用户模块
// @Success 200 {string} json
// @Router /api/auth/user/List [get]
func (UserAPI) GetUserList(c *gin.Context) {
	userList := userService.GetUserList()
	result.OkWithData(userList, c)
}

// QueryUserByUsername
// @Summary 根据用户名查询用户
// @Tags 用户模块
// @Param username query string true "用户名"
// @Success 200 {string} json
// @Router /api/auth/user/query_by_username [get]
func (UserAPI) QueryUserByUsername(c *gin.Context) {
	username := c.Query("username")
	user := userService.QueryUserByUsername(username)
	if reflect.DeepEqual(user, model.ScaAuthUser{}) {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "NotFoundUser"), c)
		return
	}
	result.OkWithData(user, c)
}

// QueryUserByUuid
// @Summary 根据uuid查询用户
// @Tags 用户模块
// @Param uuid query string true "用户uuid"
// @Success 200 {string} json
// @Router /api/auth/user/query_by_uuid [get]
func (UserAPI) QueryUserByUuid(c *gin.Context) {
	uuid := c.Query("uuid")
	user, err := userService.QueryUserByUuid(&uuid)
	if err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "NotFoundUser"), c)
		return
	}
	if reflect.DeepEqual(user, model.ScaAuthUser{}) {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "NotFoundUser"), c)
		return
	}
	result.OkWithData(user, c)
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Tags 用户模块
// @Param uuid query string true "用户uuid"
// @Success 200 {string} json
// @Router /api/auth/user/delete [delete]
func (UserAPI) DeleteUser(c *gin.Context) {
	uuid := c.Query("uuid")
	err := userService.DeleteUser(uuid)
	if err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "DeletedFailed"), c)
		return
	}
	result.OkWithMessage(ginI18n.MustGetMessage(c, "DeletedSuccess"), c)
}

// QueryUserByPhone 根据手机号查询用户
// @Summary 根据手机号查询用户
// @Tags 用户模块
// @Param phone query string true "手机号"
// @Success 200 {string} json
// @Router /api/auth/user/query_by_phone [get]
func (UserAPI) QueryUserByPhone(c *gin.Context) {
	phone := c.Query("phone")
	user := userService.QueryUserByPhone(phone)
	if reflect.DeepEqual(user, model.ScaAuthUser{}) {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "NotFoundUser"), c)
		return
	}
	result.OkWithData(user, c)
}

// AccountLogin 账号登录
// @Summary 账号登录
// @Tags 用户模块
// @Param user body dto.AccountLoginRequest true "用户信息"
// @Success 200 {string} json
// @Router /api/user/login [post]
func (UserAPI) AccountLogin(c *gin.Context) {
	accountLoginRequest := dto.AccountLoginRequest{}
	err := c.ShouldBindJSON(&accountLoginRequest)
	if err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ParamsError"), c)
		return
	}
	account := accountLoginRequest.Account
	password := accountLoginRequest.Password
	if account == "" || password == "" {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "AccountAndPasswordNotEmpty"), c)
		return
	}

	var user model.ScaAuthUser
	if utils.IsPhone(account) {
		user = userService.QueryUserByPhone(account)
	} else if utils.IsEmail(account) {
		user = userService.QueryUserByEmail(account)
	} else if utils.IsUsername(account) {
		user = userService.QueryUserByUsername(account)
	} else {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "AccountErrorFormat"), c)
		return
	}

	if reflect.DeepEqual(user, model.ScaAuthUser{}) {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "NotFoundUser"), c)
		return
	}

	if !utils.Verify(*user.Password, password) {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "PasswordError"), c)
		return
	}
	handelUserLogin(user, accountLoginRequest.AutoLogin, c)
}

// PhoneLogin 手机号登录/注册
// @Summary 手机号登录/注册
// @Tags 用户模块
// @Param user body dto.PhoneLoginRequest true "用户信息"
// @Success 200 {string} json
// @Router /api/user/phone_login [post]
func (UserAPI) PhoneLogin(c *gin.Context) {
	request := dto.PhoneLoginRequest{}
	err := c.ShouldBindJSON(&request)
	if err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ParamsError"), c)
		return
	}
	phone := request.Phone
	captcha := request.Captcha
	if phone == "" || captcha == "" {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "PhoneAndCaptchaNotEmpty"), c)
		return
	}
	if !utils.IsPhone(phone) {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "PhoneErrorFormat"), c)
		return
	}

	userChan := make(chan model.ScaAuthUser)
	go func() {
		user := userService.QueryUserByPhone(phone)
		userChan <- user
	}()

	user := <-userChan
	close(userChan)

	if reflect.DeepEqual(user, model.ScaAuthUser{}) {
		// 未注册
		codeChan := make(chan *string)
		go func() {
			code := redis.Get(constant.UserLoginSmsRedisKey + phone).Val()
			codeChan <- &code
		}()

		code := <-codeChan
		close(codeChan)

		if code == nil {
			result.FailWithMessage(ginI18n.MustGetMessage(c, "CaptchaExpired"), c)
			return
		}

		uid := idgen.NextId()
		uidStr := strconv.FormatInt(uid, 10)
		createUser := model.ScaAuthUser{
			UID:   &uidStr,
			Phone: &phone,
		}

		errChan := make(chan error)
		go func() {
			err := global.DB.Transaction(func(tx *gorm.DB) error {
				addUser, err := userService.AddUser(createUser)
				if err != nil {
					return err
				}
				_, err = global.Casbin.AddRoleForUser(uidStr, enum.User)
				if err != nil {
					return err
				}
				handelUserLogin(addUser, request.AutoLogin, c)
				return nil
			})
			errChan <- err
		}()

		err := <-errChan
		close(errChan)

		if err != nil {
			result.FailWithMessage(ginI18n.MustGetMessage(c, "RegisterUserError"), c)
			return
		}
	} else {
		codeChan := make(chan *string)
		go func() {
			code := redis.Get(constant.UserLoginSmsRedisKey + phone).Val()
			codeChan <- &code
		}()

		code := <-codeChan
		close(codeChan)

		if code == nil {
			result.FailWithMessage(ginI18n.MustGetMessage(c, "CaptchaExpired"), c)
			return
		}
		if &captcha != code {
			result.FailWithMessage(ginI18n.MustGetMessage(c, "CaptchaError"), c)
			return
		}
		handelUserLogin(user, request.AutoLogin, c)
	}
}

// RefreshHandler 刷新token
// @Summary 刷新token
// @Tags 用户模块
// @Param refresh_token query string true "刷新token"
// @Success 200 {string} json
// @Router /api/token/refresh [post]
func (UserAPI) RefreshHandler(c *gin.Context) {
	request := dto.RefreshTokenRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ParamsError"), c)
		return
	}
	refreshToken := request.RefreshToken
	if refreshToken == "" {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ParamsError"), c)
		return
	}
	parseRefreshToken, isUpd, err := utils.ParseRefreshToken(refreshToken)
	if err != nil || !isUpd {
		global.LOG.Errorln(err)
		result.FailWithMessage(ginI18n.MustGetMessage(c, "LoginExpired"), c)
		return
	}
	accessTokenString, err := utils.GenerateAccessToken(utils.AccessJWTPayload{UserID: parseRefreshToken.UserID})
	if err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "LoginExpired"), c)
		return
	}
	tokenKey := constant.UserLoginTokenRedisKey + *parseRefreshToken.UserID
	token, err := redis.Get(tokenKey).Result()
	if err != nil || token == "" {
		global.LOG.Errorln(err)
		result.FailWithMessage(ginI18n.MustGetMessage(c, "LoginExpired"), c)
		return
	}
	data := dto.ResponseData{
		AccessToken:  accessTokenString,
		RefreshToken: refreshToken,
		UID:          parseRefreshToken.UserID,
	}
	if err := redis.Set(tokenKey, data, time.Hour*24*7).Err(); err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "LoginExpired"), c)
		return
	}
	result.OkWithData(data, c)
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
		days = 24 * time.Hour
	}

	refreshToken, expiresAt := utils.GenerateRefreshToken(utils.RefreshJWTPayload{UserID: user.UID}, days)
	data := dto.ResponseData{
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

	result.OkWithData(data, c)
}

// ResetPassword 重置密码
// @Summary 重置密码
// @Tags 用户模块
// @Param user body dto.ResetPasswordRequest true "用户信息"
// @Success 200 {string} json
// @Router /api/user/reset_password [post]
func (UserAPI) ResetPassword(c *gin.Context) {
	var resetPasswordRequest dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&resetPasswordRequest); err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ParamsError"), c)
		return
	}

	phone := resetPasswordRequest.Phone
	captcha := resetPasswordRequest.Captcha
	password := resetPasswordRequest.Password
	repassword := resetPasswordRequest.Repassword

	if phone == "" || captcha == "" || password == "" || repassword == "" {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ParamsError"), c)
		return
	}

	if !utils.IsPhone(phone) {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "PhoneError"), c)
		return
	}

	if password != repassword {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "PasswordNotSame"), c)
		return
	}

	if !utils.IsPassword(password) {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "PasswordError"), c)
		return
	}

	// 使用事务确保验证码检查和密码更新的原子性
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "DatabaseError"), c)
		return
	}

	code := redis.Get(constant.UserLoginSmsRedisKey + phone).Val()
	if code == "" {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CaptchaExpired"), c)
		return
	}

	if captcha != code {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CaptchaError"), c)
		return
	}

	// 验证码检查通过后立即删除或标记为已使用
	if err := redis.Del(constant.UserLoginSmsRedisKey + phone).Err(); err != nil {
		tx.Rollback()
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ResetPasswordError"), c)
		return
	}

	user := userService.QueryUserByPhone(phone)
	if reflect.DeepEqual(user, model.ScaAuthUser{}) {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "PhoneNotRegister"), c)
		return
	}

	encrypt, err := utils.Encrypt(password)
	if err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ResetPasswordError")+": "+err.Error(), c)
		return
	}

	if err := userService.UpdateUser(phone, encrypt); err != nil {
		tx.Rollback()
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ResetPasswordError"), c)
		return
	}

	tx.Commit()
	result.OkWithMessage(ginI18n.MustGetMessage(c, "ResetPasswordSuccess"), c)
}

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
	m := ua.Model()
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
		Model:           &m,
		Platform:        &platform,
		EngineName:      &engine,
		EngineVersion:   &engineVersion,
	}

	mu.Lock()
	defer mu.Unlock()

	userDevice, err := userDeviceService.GetUserDeviceByUIDIPAgent(*user.UID, ip, userAgent)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = userDeviceService.AddUserDevice(&device)
		if err != nil {
			global.LOG.Errorln(err)
			return false
		}
	} else if err != nil {
		global.LOG.Errorln(err)
		return false
	} else {
		err := userDeviceService.UpdateUserDevice(userDevice.ID, &device)
		if err != nil {
			global.LOG.Errorln(err)
			return false
		}
	}

	return true
}
