package oauth_api

import (
	"encoding/json"
	"errors"
	"github.com/ArtisanCloud/PowerLibs/v3/fmt"
	"github.com/ArtisanCloud/PowerLibs/v3/http/helper"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/basicService/qrCode/response"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/contract"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/messages"
	models2 "github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/models"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount/server/handlers/models"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/yitter/idgenerator-go/idgen"
	"gorm.io/gorm"
	"schisandra-cloud-album/api/user_api/dto"
	"schisandra-cloud-album/api/websocket_api"
	"schisandra-cloud-album/common/constant"
	"schisandra-cloud-album/common/enum"
	"schisandra-cloud-album/common/redis"
	"schisandra-cloud-album/common/result"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/model"
	"schisandra-cloud-album/service"
	"schisandra-cloud-album/utils"
	"strconv"
	"strings"
	"time"
)

var userService = service.Service.UserService
var userRoleService = service.Service.UserRoleService
var userSocialService = service.Service.UserSocialService
var rolePermissionService = service.Service.RolePermissionService
var permissionServiceService = service.Service.PermissionService
var roleService = service.Service.RoleService

// GenerateClientId 生成客户端ID
// @Summary 生成客户端ID
// @Description 生成客户端ID
// @Produce json
// @Success 200 {object} result.Result{data=string} "客户端ID"
// @Router /api/oauth/generate_client_id [get]
func (OAuthAPI) GenerateClientId(c *gin.Context) {
	ip := c.ClientIP()
	clientId := redis.Get(constant.UserLoginClientRedisKey + ip).Val()
	if clientId != "" {
		result.OkWithData(clientId, c)
		return
	}
	v1 := uuid.NewV1()
	redis.Set(constant.UserLoginClientRedisKey+ip, v1.String(), 0)
	result.OkWithData(v1.String(), c)
	return
}

// CallbackNotify 微信回调验证
// @Summary 微信回调验证
// @Description 微信回调验证
// @Produce json
// @Success 200 {object} result.Result{data=string} "验证结果"
// @Router /api/oauth/callback_notify [POST]
func (OAuthAPI) CallbackNotify(c *gin.Context) {
	rs, err := global.Wechat.Server.Notify(c.Request, func(event contract.EventInterface) interface{} {
		fmt.Dump("event", event)

		switch event.GetMsgType() {

		case models2.CALLBACK_MSG_TYPE_EVENT:
			switch event.GetEvent() {
			case models.CALLBACK_EVENT_SUBSCRIBE:
				msg := models.EventSubscribe{}
				err := event.ReadMessage(&msg)
				if err != nil {
					println(err.Error())
					return "error"
				}
				key := strings.TrimPrefix(msg.EventKey, "qrscene_")
				res := wechatLoginHandler(msg.FromUserName, key)
				if !res {
					return messages.NewText(ginI18n.MustGetMessage(c, "LoginFailed"))
				}
				return messages.NewText(ginI18n.MustGetMessage(c, "LoginSuccess"))

			case models.CALLBACK_EVENT_UNSUBSCRIBE:
				msg := models.EventUnSubscribe{}
				err := event.ReadMessage(&msg)
				if err != nil {
					println(err.Error())
					return "error"
				}
				fmt.Dump(msg)
				return messages.NewText("再见，我的宝！")

			case models.CALLBACK_EVENT_SCAN:
				msg := models.EventScan{}
				err := event.ReadMessage(&msg)
				if err != nil {
					println(err.Error())
					return "error"
				}
				res := wechatLoginHandler(msg.FromUserName, msg.EventKey)
				if !res {
					return messages.NewText(ginI18n.MustGetMessage(c, "LoginFailed"))
				}
				return messages.NewText(ginI18n.MustGetMessage(c, "LoginSuccess"))

			}

		case models2.CALLBACK_MSG_TYPE_TEXT:
			msg := models.MessageText{}
			err := event.ReadMessage(&msg)
			if err != nil {
				println(err.Error())
				return "error"
			}
			fmt.Dump(msg)
		}
		return messages.NewText("ok")

	})
	if err != nil {
		panic(err)
	}
	err = helper.HttpResponseSend(rs, c.Writer)
	if err != nil {
		panic(err)
	}
}

// CallbackVerify 微信回调验证
// @Summary 微信回调验证
// @Description 微信回调验证
// @Produce json
// @Success 200 {object} result.Result{data=string} "验证结果"
// @Router /api/oauth/callback_verify [get]
func (OAuthAPI) CallbackVerify(c *gin.Context) {
	rs, err := global.Wechat.Server.VerifyURL(c.Request)
	if err != nil {
		panic(err)
	}
	err = helper.HttpResponseSend(rs, c.Writer)
}

// GetTempQrCode 获取临时二维码
// @Summary 获取临时二维码
// @Description 获取临时二维码
// @Produce json
// @Param client_id query string true "客户端ID"
// @Success 200 {object} result.Result{data=string} "临时二维码"
// @Router /api/oauth/get_temp_qrcode [get]
func (OAuthAPI) GetTempQrCode(c *gin.Context) {
	clientId := c.Query("client_id")
	ip := c.ClientIP()
	if clientId == "" {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ParamsError"), c)
		return
	}
	qrcode := redis.Get(constant.UserLoginQrcodeRedisKey + ip + ":" + clientId).Val()

	if qrcode != "" {
		data := response.ResponseQRCodeCreate{}
		err := json.Unmarshal([]byte(qrcode), &data)
		if err != nil {
			return
		}
		result.OK(ginI18n.MustGetMessage(c, "QRCodeGetSuccess"), data.Url, c)
		return
	}
	data, err := global.Wechat.QRCode.Temporary(c.Request.Context(), clientId, 30*24*3600)
	if err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "QRCodeGetFailed"), c)
		return
	}
	serializedData, err := json.Marshal(data)
	if err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "QRCodeGetFailed"), c)
		return
	}
	wrong := redis.Set(constant.UserLoginQrcodeRedisKey+ip+":"+clientId, serializedData, time.Hour*24*30).Err()

	if wrong != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "QRCodeGetFailed"), c)
		return
	}
	result.OK(ginI18n.MustGetMessage(c, "QRCodeGetSuccess"), data.Url, c)
}

// wechatLoginHandler 微信登录处理
func wechatLoginHandler(openId string, clientId string) bool {
	if openId == "" {
		return false
	}
	authUserSocial, err := userSocialService.QueryUserSocialByOpenID(openId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		uid := idgen.NextId()
		uidStr := strconv.FormatInt(uid, 10)
		createUser := model.ScaAuthUser{
			UID:      &uidStr,
			Username: &openId,
		}
		addUser, err := userService.AddUser(createUser)
		if err != nil {
			return false
		}
		wechat := enum.OAuthSourceWechat
		userSocial := model.ScaAuthUserSocial{
			UserID: &addUser.ID,
			OpenID: &openId,
			Source: &wechat,
		}
		wrong := userSocialService.AddUserSocial(userSocial)
		if wrong != nil {
			return false
		}
		userRole := model.ScaAuthUserRole{
			UserID: addUser.ID,
			RoleID: enum.User,
		}
		e := userRoleService.AddUserRole(userRole)
		if e != nil {
			return false
		}
		res := handelUserLogin(addUser, clientId)
		if !res {
			return false
		}
		return true
	} else {
		user, err := userService.QueryUserById(authUserSocial.UserID)
		if err != nil {
			return false
		}
		res := handelUserLogin(user, clientId)
		if !res {
			return false
		}
		return true
	}
}

// handelUserLogin 处理用户登录
func handelUserLogin(user model.ScaAuthUser, clientId string) bool {
	ids, err := userRoleService.GetUserRoleIdsByUserId(user.ID)
	if err != nil {
		return false
	}
	permissionIds := rolePermissionService.QueryPermissionIdsByRoleId(ids)
	permissions, err := permissionServiceService.GetPermissionsByIds(permissionIds)
	if err != nil {
		return false
	}
	serializedPermissions, err := json.Marshal(permissions)
	if err != nil {
		return false
	}
	wrong := redis.Set(constant.UserAuthPermissionRedisKey+*user.UID, serializedPermissions, 0).Err()
	if wrong != nil {
		return false
	}
	roleList, err := roleService.GetRoleListByIds(ids)
	if err != nil {
		return false
	}
	serializedRoleList, err := json.Marshal(roleList)
	if err != nil {
		return false
	}
	er := redis.Set(constant.UserAuthRoleRedisKey+*user.UID, serializedRoleList, 0).Err()
	if er != nil {
		return false
	}
	accessToken, refreshToken, expiresAt := utils.GenerateAccessTokenAndRefreshToken(utils.JWTPayload{UserID: user.UID, RoleID: ids})

	data := dto.ResponseData{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		UID:          user.UID,
	}
	fail := redis.Set(constant.UserLoginTokenRedisKey+*user.UID, data, time.Hour*24*7).Err()
	w := redis.Set(constant.UserLoginWechatRedisKey+clientId, data, time.Minute*5).Err()
	if fail != nil || w != nil {
		return false
	}
	res := websocket_api.SendMessageData(clientId, data)
	if !res {
		return false
	}
	return true
}
