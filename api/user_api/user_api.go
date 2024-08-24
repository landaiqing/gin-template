package user_api

import (
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/yitter/idgenerator-go/idgen"
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

// AddUser 添加用户
// @Summary 添加用户
// @Tags 用户模块
// @Param user body dto.AddUserRequest true "用户信息"
// @Success 200 {string} json
// @Router /api/user/add [post]
func (UserAPI) AddUser(c *gin.Context) {
	addUserRequest := dto.AddUserRequest{}
	err := c.ShouldBindJSON(&addUserRequest)
	if err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ParamsError"), c)
		return
	}

	username := userService.QueryUserByUsername(addUserRequest.Username)
	if !reflect.DeepEqual(username, model.ScaAuthUser{}) {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "UsernameExists"), c)
		return
	}

	phone := userService.QueryUserByPhone(addUserRequest.Phone)
	if !reflect.DeepEqual(phone, model.ScaAuthUser{}) {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "PhoneExists"), c)
		return
	}
	encrypt, err := utils.Encrypt(addUserRequest.Password)
	if err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "AddUserError"), c)
		return
	}
	uid := idgen.NextId()
	uidStr := strconv.FormatInt(uid, 10)
	user := model.ScaAuthUser{
		UID:      &uidStr,
		Username: &addUserRequest.Username,
		Password: &encrypt,
		Phone:    &addUserRequest.Phone,
	}
	_, err = userService.AddUser(user)
	if err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "AddUserError"), c)
		return
	}
	_, err = global.Casbin.AddRoleForUser(uidStr, enum.User)
	if err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "AddUserRoleError"), c)
		return
	}
	result.OkWithMessage(ginI18n.MustGetMessage(c, "AddUserSuccess"), c)
	return
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
	isPhone := utils.IsPhone(phone)
	if !isPhone {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "PhoneErrorFormat"), c)
		return
	}

	// 异步查询用户信息
	userChan := make(chan *model.ScaAuthUser)
	go func() {
		user := userService.QueryUserByPhone(phone)
		userChan <- &user
	}()

	// 异步获取验证码
	codeChan := make(chan string)
	go func() {
		code := redis.Get(constant.UserLoginSmsRedisKey + phone)
		if code == nil {
			codeChan <- ""
		} else {
			codeChan <- code.Val()
		}
	}()

	user := <-userChan
	code := <-codeChan

	if reflect.DeepEqual(user, model.ScaAuthUser{}) {
		// 未注册
		if code == "" {
			result.FailWithMessage(ginI18n.MustGetMessage(c, "CaptchaExpired"), c)
			return
		} else {
			uid := idgen.NextId()
			uidStr := strconv.FormatInt(uid, 10)
			createUser := model.ScaAuthUser{
				UID:   &uidStr,
				Phone: &phone,
			}
			addUser, err := userService.AddUser(createUser)
			if err != nil {
				result.FailWithMessage(ginI18n.MustGetMessage(c, "RegisterUserError"), c)
				return
			}

			_, err = global.Casbin.AddRoleForUser(uidStr, enum.User)
			if err != nil {
				result.FailWithMessage(ginI18n.MustGetMessage(c, "RegisterUserError"), c)
				return
			}
			handelUserLogin(addUser, request.AutoLogin, c)
			return
		}
	} else {
		if code == "" {
			result.FailWithMessage(ginI18n.MustGetMessage(c, "CaptchaExpired"), c)
			return
		} else {
			if captcha != code {
				result.FailWithMessage(ginI18n.MustGetMessage(c, "CaptchaError"), c)
				return
			} else {
				handelUserLogin(*user, request.AutoLogin, c)
				return
			}
		}
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
	if err != nil {
		global.LOG.Errorln(err)
		result.FailWithMessage(ginI18n.MustGetMessage(c, "LoginExpired"), c)
		return
	}
	if !isUpd {
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
	if token == "" || err != nil {
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
	accessToken, err := utils.GenerateAccessToken(utils.AccessJWTPayload{UserID: user.UID})
	if err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "LoginFailed"), c)
		return
	}
	var days time.Duration
	if autoLogin {
		days = time.Hour * 24 * 7
	} else {
		days = time.Hour * 24 * 1
	}
	refreshToken, expiresAt := utils.GenerateRefreshToken(utils.RefreshJWTPayload{UserID: user.UID}, days)
	data := dto.ResponseData{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		UID:          user.UID,
	}
	fail := redis.Set(constant.UserLoginTokenRedisKey+*user.UID, data, time.Hour*24*1).Err()
	if fail != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "LoginFailed"), c)
		return
	}
	result.OkWithData(data, c)
	return
}

// ResetPassword 重置密码
// @Summary 重置密码
// @Tags 用户模块
// @Param user body dto.ResetPasswordRequest true "用户信息"
// @Success 200 {string} json
// @Router /api/user/reset_password [post]
func (UserAPI) ResetPassword(c *gin.Context) {
	resetPasswordRequest := dto.ResetPasswordRequest{}
	err := c.ShouldBindJSON(&resetPasswordRequest)
	if err != nil {
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
	isPhone := utils.IsPhone(phone)
	if !isPhone {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "PhoneErrorFormat"), c)
		return
	}
	code := redis.Get(constant.UserLoginSmsRedisKey + phone)
	if code == nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CaptchaExpired"), c)
		return
	} else {
		if captcha != code.Val() {
			result.FailWithMessage(ginI18n.MustGetMessage(c, "CaptchaError"), c)
			return
		}
	}
	user := userService.QueryUserByPhone(phone)
	if reflect.DeepEqual(user, model.ScaAuthUser{}) {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "PhoneNotRegister"), c)
		return
	}
	encrypt, err := utils.Encrypt(password)
	if err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ResetPasswordError"), c)
		return
	}
	wrong := userService.UpdateUser(phone, encrypt)
	if wrong != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ResetPasswordError"), c)
		return
	}
	result.OkWithMessage(ginI18n.MustGetMessage(c, "ResetPasswordSuccess"), c)
	return
}
