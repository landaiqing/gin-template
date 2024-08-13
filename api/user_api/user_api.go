package user_api

import (
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/yitter/idgenerator-go/idgen"
	"reflect"
	"schisandra-cloud-album/api/user_api/dto"
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
var userRoleService = service.Service.UserRoleService

// GetUserList
// @Summary 获取所有用户列表
// @Tags 鉴权模块
// @Success 200 {string} json
// @Router /api/auth/user/List [get]
func (UserAPI) GetUserList(c *gin.Context) {
	userList := userService.GetUserList()
	result.OkWithData(userList, c)
}

// QueryUserByUsername
// @Summary 根据用户名查询用户
// @Tags 鉴权模块
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
// @Tags 鉴权模块
// @Param uuid query string true "用户uuid"
// @Success 200 {string} json
// @Router /api/auth/user/query_by_uuid [get]
func (UserAPI) QueryUserByUuid(c *gin.Context) {
	uuid := c.Query("uuid")
	user := userService.QueryUserByUuid(uuid)
	if reflect.DeepEqual(user, model.ScaAuthUser{}) {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "NotFoundUser"), c)
		return
	}
	result.OkWithData(user, c)
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Tags 鉴权模块
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
// @Tags 鉴权模块
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
// @Tags 鉴权模块
// @Param account query string true "账号"
// @Param password query string true "密码"
// @Success 200 {string} json
// @Router /api/user/login [post]
func (UserAPI) AccountLogin(c *gin.Context) {
	account := c.PostForm("account")
	password := c.PostForm("password")
	isPhone := utils.IsPhone(account)
	if isPhone {
		user := userService.QueryUserByPhone(account)
		if reflect.DeepEqual(user, model.ScaAuthUser{}) {
			result.FailWithMessage(ginI18n.MustGetMessage(c, "PhoneNotRegister"), c)
			return
		} else {
			verify := utils.Verify(password, *user.Password)
			if verify {
				result.OkWithData(user, c)
				return
			} else {
				result.FailWithMessage(ginI18n.MustGetMessage(c, "PasswordError"), c)
				return
			}
		}
	}
	isEmail := utils.IsEmail(account)
	if isEmail {
		user := userService.QueryUserByEmail(account)
		if reflect.DeepEqual(user, model.ScaAuthUser{}) {
			result.FailWithMessage(ginI18n.MustGetMessage(c, "EmailNotRegister"), c)
			return
		} else {
			verify := utils.Verify(password, *user.Password)
			if verify {
				result.OkWithData(user, c)
				return
			} else {
				result.FailWithMessage(ginI18n.MustGetMessage(c, "PasswordError"), c)
				return
			}
		}
	}
	isUsername := utils.IsUsername(account)
	if isUsername {
		user := userService.QueryUserByUsername(account)
		if reflect.DeepEqual(user, model.ScaAuthUser{}) {
			result.FailWithMessage(ginI18n.MustGetMessage(c, "UsernameNotRegister"), c)
			return
		} else {
			verify := utils.Verify(password, *user.Password)
			if verify {
				result.OkWithData(user, c)
				return
			} else {
				result.FailWithMessage(ginI18n.MustGetMessage(c, "PasswordError"), c)
				return
			}
		}

	}
}

// PhoneLogin 手机号登录/注册
// @Summary 手机号登录/注册
// @Tags 鉴权模块
// @Param phone query string true "手机号"
// @Param captcha query string true "验证码"
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

	user := userService.QueryUserByPhone(phone)
	if reflect.DeepEqual(user, model.ScaAuthUser{}) {
		// 未注册
		code := redis.Get("user:login:sms:" + phone)
		if code == nil {
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
			userRole := model.ScaAuthUserRole{
				UserID: addUser.ID,
				RoleID: enum.User,
			}
			e := userRoleService.AddUserRole(userRole)
			if e != nil {
				result.FailWithMessage(ginI18n.MustGetMessage(c, "LoginFailed"), c)
				return
			}
			ids, err := userRoleService.GetUserRoleIdsByUserId(addUser.ID)
			if err != nil {
				result.FailWithMessage(ginI18n.MustGetMessage(c, "LoginFailed"), c)
				return
			}
			accessToken, refreshToken, expiresAt := utils.GenerateAccessTokenAndRefreshToken(utils.JWTPayload{UserID: addUser.UID, RoleID: ids})

			data := dto.ResponseData{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
				ExpiresAt:    expiresAt,
				UID:          addUser.UID,
			}
			fail := redis.Set("user:login:token:"+*addUser.UID, data, time.Hour*24*30).Err()
			if fail != nil {
				result.FailWithMessage(ginI18n.MustGetMessage(c, "LoginFailed"), c)
				return
			}
			result.OkWithData(data, c)
			return
		}
	} else {
		code := redis.Get("user:login:sms:" + phone)
		if code == nil {
			result.FailWithMessage(ginI18n.MustGetMessage(c, "CaptchaExpired"), c)
			return
		} else {
			if captcha != code.Val() {
				result.FailWithMessage(ginI18n.MustGetMessage(c, "CaptchaError"), c)
				return
			} else {
				ids, err := userRoleService.GetUserRoleIdsByUserId(user.ID)
				if err != nil {
					result.FailWithMessage(ginI18n.MustGetMessage(c, "LoginFailed"), c)
					return
				}
				accessToken, refreshToken, expiresAt := utils.GenerateAccessTokenAndRefreshToken(utils.JWTPayload{UserID: user.UID, RoleID: ids})

				data := dto.ResponseData{
					AccessToken:  accessToken,
					RefreshToken: refreshToken,
					ExpiresAt:    expiresAt,
					UID:          user.UID,
				}
				fail := redis.Set("user:login:token:"+*user.UID, data, time.Hour*24*30).Err()
				if fail != nil {
					result.FailWithMessage(ginI18n.MustGetMessage(c, "LoginFailed"), c)
					return
				}
				result.OkWithData(data, c)
				return
			}
		}

	}

}

// RefreshHandler 刷新token
// @Summary 刷新token
// @Tags 鉴权模块
// @Param refresh_token query string true "刷新token"
// @Success 200 {string} json
// @Router /api/auth/token/refresh [post]
func (UserAPI) RefreshHandler(c *gin.Context) {
	refreshToken := c.Query("refresh_token")
	if refreshToken == "" {
		result.FailWithMessage("refresh_token不能为空！", c)
		return
	}
	parseRefreshToken, isUpd, err := utils.ParseToken(refreshToken)
	if err != nil {
		global.LOG.Errorln(err)
		result.FailWithMessage(ginI18n.MustGetMessage(c, "LoginExpired"), c)
		return
	}
	if isUpd {
		accessTokenString, refreshTokenString, expiresAt := utils.GenerateAccessTokenAndRefreshToken(utils.JWTPayload{UserID: parseRefreshToken.UserID, RoleID: parseRefreshToken.RoleID})
		data := dto.ResponseData{
			AccessToken:  accessTokenString,
			RefreshToken: refreshTokenString,
			ExpiresAt:    expiresAt,
			UID:          parseRefreshToken.UserID,
		}
		fail := redis.Set("user:login:token:"+*parseRefreshToken.UserID, data, time.Hour*24*30).Err()
		if fail != nil {
			result.FailWithMessage(ginI18n.MustGetMessage(c, "LoginExpired"), c)
			return
		}
		result.OkWithData(data, c)
		return
	}
}
