package user_controller

import (
	"errors"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/mssola/useragent"
	"github.com/yitter/idgenerator-go/idgen"
	"gorm.io/gorm"
	"reflect"
	"schisandra-cloud-album/common/constant"
	"schisandra-cloud-album/common/enum"
	"schisandra-cloud-album/common/randomname"
	"schisandra-cloud-album/common/redis"
	"schisandra-cloud-album/common/result"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/model"
	"schisandra-cloud-album/service/impl"
	"schisandra-cloud-album/utils"
	"strconv"
	"sync"
)

type UserController struct{}

var mu sync.Mutex
var userService = impl.UserServiceImpl{}
var userDeviceService = impl.UserDeviceServiceImpl{}

// GetUserList
// @Summary 获取所有用户列表
// @Tags 用户模块
// @Success 200 {string} json
// @Router /controller/auth/user/List [get]
func (UserController) GetUserList(c *gin.Context) {
	userList := userService.GetUserListService()
	result.OkWithData(userList, c)
}

// QueryUserByUsername
// @Summary 根据用户名查询用户
// @Tags 用户模块
// @Param username query string true "用户名"
// @Success 200 {string} json
// @Router /controller/auth/user/query_by_username [get]
func (UserController) QueryUserByUsername(c *gin.Context) {
	username := c.Query("username")
	user := userService.QueryUserByUsernameService(username)
	if reflect.DeepEqual(user, model.ScaAuthUser{}) {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "NotFoundUser"), c)
		return
	}
	result.OkWithData(user, c)
}

// QueryUserByUuid
// @Summary 根据uuid查询用户
// @Tags 用户模块
// @Param uid query string true "用户uid"
// @Success 200 {string} json
// @Router /controller/auth/user/query_by_uid [get]
func (UserController) QueryUserByUuid(c *gin.Context) {
	uid := c.Query("uid")
	user := userService.QueryUserByUuidService(&uid)
	if user.ID == 0 {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "NotFoundUser"), c)
		return
	}
	result.OkWithData(user, c)
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Tags 用户模块
// @Param uid query string true "用户uid"
// @Success 200 {string} json
// @Router /controller/auth/user/delete [delete]
func (UserController) DeleteUser(c *gin.Context) {
	uid := c.Query("uid")
	err := userService.DeleteUserService(uid)
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
// @Router /controller/auth/user/query_by_phone [get]
func (UserController) QueryUserByPhone(c *gin.Context) {
	phone := c.Query("phone")
	user := userService.QueryUserByPhoneService(phone)
	if user.ID == 0 {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "NotFoundUser"), c)
		return
	}
	result.OkWithData(user, c)
}

// AccountLogin 账号登录
// @Summary 账号登录
// @Tags 用户模块
// @Param user body AccountLoginRequest true "用户信息"
// @Success 200 {string} json
// @Router /controller/user/login [post]
func (UserController) AccountLogin(c *gin.Context) {
	accountLoginRequest := AccountLoginRequest{}
	err := c.ShouldBindJSON(&accountLoginRequest)
	if err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ParamsError"), c)
		return
	}
	rotateData := utils.CheckRotateData(accountLoginRequest.Angle, accountLoginRequest.Key)
	if !rotateData {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CaptchaVerifyError"), c)
		return
	}
	account := accountLoginRequest.Account
	password := accountLoginRequest.Password

	var user model.ScaAuthUser
	if utils.IsPhone(account) {
		user = userService.QueryUserByPhoneService(account)
	} else if utils.IsEmail(account) {
		user = userService.QueryUserByEmailService(account)
	} else if utils.IsUsername(account) {
		user = userService.QueryUserByUsernameService(account)
	} else {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "AccountErrorFormat"), c)
		return
	}

	if user.ID == 0 {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "NotFoundUser"), c)
		return
	}

	if !utils.Verify(user.Password, password) {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "PasswordError"), c)
		return
	}
	data, res := userService.HandelUserLogin(user, accountLoginRequest.AutoLogin, c)
	if !res {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "LoginFailed"), c)
		return
	}
	result.OkWithData(data, c)
}

// PhoneLogin 手机号登录/注册
// @Summary 手机号登录/注册
// @Tags 用户模块
// @Param user body PhoneLoginRequest true "用户信息"
// @Success 200 {string} json
// @Router /controller/user/phone_login [post]
func (UserController) PhoneLogin(c *gin.Context) {
	request := PhoneLoginRequest{}
	err := c.ShouldBind(&request)
	if err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ParamsError"), c)
		return
	}
	phone := request.Phone
	captcha := request.Captcha
	autoLogin := request.AutoLogin
	if !utils.IsPhone(phone) {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "PhoneErrorFormat"), c)
		return
	}
	// 获取验证码
	code := redis.Get(constant.UserLoginSmsRedisKey + phone).Val()
	if code == "" {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CaptchaExpired"), c)
		return
	}

	if captcha != code {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CaptchaError"), c)
		return
	}
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// 查询用户
	user := userService.QueryUserByPhoneService(phone)

	if user.ID == 0 {
		// 未注册
		uid := idgen.NextId()
		uidStr := strconv.FormatInt(uid, 10)
		avatar := utils.GenerateAvatar(uidStr)
		name := randomname.GenerateName()
		createUser := model.ScaAuthUser{
			UID:      uidStr,
			Phone:    phone,
			Avatar:   avatar,
			Nickname: name,
			Gender:   enum.Male,
		}

		addUser, w := userService.AddUserService(createUser)
		if w != nil {
			tx.Rollback()
			return
		}
		_, err = global.Casbin.AddRoleForUser(uidStr, enum.User)
		if err != nil {
			tx.Rollback()
			return
		}
		data, res := userService.HandelUserLogin(*addUser, autoLogin, c)
		if !res {
			tx.Rollback()
			result.FailWithMessage(ginI18n.MustGetMessage(c, "LoginFailed"), c)
			return
		}
		tx.Commit()
		result.OkWithData(data, c)
	} else {
		data, res := userService.HandelUserLogin(user, autoLogin, c)
		if !res {
			tx.Rollback()
			result.FailWithMessage(ginI18n.MustGetMessage(c, "LoginFailed"), c)
			return
		}
		tx.Commit()
		result.OkWithData(data, c)
	}
}

// RefreshHandler 刷新token
// @Summary 刷新token
// @Tags 用户模块
// @Param refresh_token query string true "刷新token"
// @Success 200 {string} json
// @Router /controller/token/refresh [post]
func (UserController) RefreshHandler(c *gin.Context) {
	request := RefreshTokenRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ParamsError"), c)
		return
	}
	data, res := userService.RefreshTokenService(request.RefreshToken)
	if !res {
		result.FailWithCodeAndMessage(403, ginI18n.MustGetMessage(c, "LoginExpired"), c)
		return
	}
	result.OkWithData(data, c)
}

// ResetPassword 重置密码
// @Summary 重置密码
// @Tags 用户模块
// @Param user body ResetPasswordRequest true "用户信息"
// @Success 200 {string} json
// @Router /controller/user/reset_password [post]
func (UserController) ResetPassword(c *gin.Context) {
	var resetPasswordRequest ResetPasswordRequest
	if err := c.ShouldBindJSON(&resetPasswordRequest); err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ParamsError"), c)
		return
	}

	phone := resetPasswordRequest.Phone
	captcha := resetPasswordRequest.Captcha
	password := resetPasswordRequest.Password
	repassword := resetPasswordRequest.Repassword
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

	user := userService.QueryUserByPhoneService(phone)
	if user.ID == 0 {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "PhoneNotRegister"), c)
		return
	}

	encrypt, err := utils.Encrypt(password)
	if err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ResetPasswordError")+": "+err.Error(), c)
		return
	}

	if err = userService.UpdateUserService(phone, encrypt); err != nil {
		tx.Rollback()
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ResetPasswordError"), c)
		return
	}

	tx.Commit()
	result.OkWithMessage(ginI18n.MustGetMessage(c, "ResetPasswordSuccess"), c)
}

// Logout 退出登录
// @Summary 退出登录
// @Tags 用户模块
// @Success 200 {string} json
// @Router /controller/auth/user/logout [post]
func (UserController) Logout(c *gin.Context) {
	userId := c.Query("user_id")
	if userId == "" {
		global.LOG.Errorln("userId is empty")
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ParamsError"), c)
		return
	}

	tokenKey := constant.UserLoginTokenRedisKey + userId
	del := redis.Del(tokenKey)

	if del.Err() != nil {
		global.LOG.Errorln(del.Err())
		result.FailWithMessage(ginI18n.MustGetMessage(c, "LogoutFailed"), c)
		return
	}
	ip := utils.GetClientIP(c)
	key := constant.UserLoginClientRedisKey + ip
	del = redis.Del(key)
	if del.Err() != nil {
		global.LOG.Errorln(del.Err())
		result.FailWithMessage(ginI18n.MustGetMessage(c, "LogoutFailed"), c)
		return
	}
	result.OkWithMessage(ginI18n.MustGetMessage(c, "LogoutSuccess"), c)
}

// GetUserLoginDevice 获取用户登录设备
func (UserController) GetUserLoginDevice(c *gin.Context) {
	userId := c.Query("user_id")
	if userId == "" {
		return
	}
	userAgent := c.GetHeader("User-Agent")
	if userAgent == "" {
		global.LOG.Errorln("user-agent is empty")
		return
	}
	ua := useragent.New(userAgent)

	ip := utils.GetClientIP(c)
	location, err := global.IP2Location.SearchByStr(ip)
	location = utils.RemoveZeroAndAdjust(location)
	if err != nil {
		global.LOG.Errorln(err)
		return
	}
	isBot := ua.Bot()
	browser, browserVersion := ua.Browser()
	os := ua.OS()
	mobile := ua.Mobile()
	mozilla := ua.Mozilla()
	platform := ua.Platform()
	engine, engineVersion := ua.Engine()
	device := model.ScaAuthUserDevice{
		UserID:          userId,
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
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	userDevice, err := userDeviceService.GetUserDeviceByUIDIPAgentService(userId, ip, userAgent)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = userDeviceService.AddUserDeviceService(&device)
		if err != nil {
			tx.Rollback()
			global.LOG.Errorln(err)
			return
		}
		tx.Commit()
		return
	} else {
		err = userDeviceService.UpdateUserDeviceService(userDevice.ID, &device)
		if err != nil {
			tx.Rollback()
			global.LOG.Errorln(err)
			return
		}
		tx.Commit()
		return
	}
}
