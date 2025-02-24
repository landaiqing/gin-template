package oauth_controller

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/ArtisanCloud/PowerLibs/v3/http/helper"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/basicService/qrCode/response"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/contract"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/messages"
	models2 "github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/models"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount/server/handlers/models"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/yitter/idgenerator-go/idgen"
	"gorm.io/gorm"

	"schisandra-cloud-album/common/constant"
	"schisandra-cloud-album/common/enum"
	"schisandra-cloud-album/common/randomname"
	"schisandra-cloud-album/common/redis"
	"schisandra-cloud-album/common/result"
	"schisandra-cloud-album/common/types"
	"schisandra-cloud-album/controller/websocket_controller/qr_ws_controller"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/model"
	"schisandra-cloud-album/utils"
)

// CallbackNotify 微信回调
// @Summary 微信回调
// @Tags 微信公众号
// @Description 微信回调
// @Produce json
// @Router /controller/oauth/callback_notify [POST]
func (OAuthController) CallbackNotify(c *gin.Context) {
	rs, err := global.Wechat.Server.Notify(c.Request, func(event contract.EventInterface) interface{} {
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
				res := wechatLoginHandler(msg.FromUserName, key, c)
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
				return messages.NewText("ok")

			case models.CALLBACK_EVENT_SCAN:
				msg := models.EventScan{}
				err := event.ReadMessage(&msg)
				if err != nil {
					println(err.Error())
					return "error"
				}
				res := wechatLoginHandler(msg.FromUserName, msg.EventKey, c)
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
// @Tags 微信公众号
// @Description 微信回调验证
// @Produce json
// @Router /controller/oauth/callback_verify [get]
func (OAuthController) CallbackVerify(c *gin.Context) {
	rs, err := global.Wechat.Server.VerifyURL(c.Request)
	if err != nil {
		panic(err)
	}
	err = helper.HttpResponseSend(rs, c.Writer)
}

// GetTempQrCode 获取临时二维码
// @Summary 获取临时二维码
// @Tags 微信公众号
// @Description 获取临时二维码
// @Produce json
// @Param client_id query string true "客户端ID"
// @Router /controller/oauth/get_temp_qrcode [get]
func (OAuthController) GetTempQrCode(c *gin.Context) {
	clientId := c.Query("client_id")
	if clientId == "" {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ParamsError"), c)
		return
	}
	ip := utils.GetClientIP(c) // 使用工具函数获取客户端IP
	key := constant.UserLoginQrcodeRedisKey + ip

	// 从Redis获取二维码数据
	qrcode := redis.Get(key).Val()
	if qrcode != "" {
		data := new(response.ResponseQRCodeCreate)
		if err := json.Unmarshal([]byte(qrcode), data); err != nil {
			global.LOG.Error(err)
			result.FailWithMessage(ginI18n.MustGetMessage(c, "QRCodeGetFailed"), c)
			return
		}
		result.OK(ginI18n.MustGetMessage(c, "QRCodeGetSuccess"), data.Url, c)
		return
	}

	// 生成临时二维码
	data, err := global.Wechat.QRCode.Temporary(c.Request.Context(), clientId, 7*24*3600)
	if err != nil {
		global.LOG.Error(err)
		result.FailWithMessage(ginI18n.MustGetMessage(c, "QRCodeGetFailed"), c)
		return
	}

	// 序列化数据并存储到Redis
	serializedData, err := json.Marshal(data)
	if err != nil {
		global.LOG.Error(err)
		result.FailWithMessage(ginI18n.MustGetMessage(c, "QRCodeGetFailed"), c)
		return
	}
	if err := redis.Set(key, serializedData, time.Hour*24*7).Err(); err != nil {
		global.LOG.Error(err)
		result.FailWithMessage(ginI18n.MustGetMessage(c, "QRCodeGetFailed"), c)
		return
	}

	result.OK(ginI18n.MustGetMessage(c, "QRCodeGetSuccess"), data.Url, c)
}

// wechatLoginHandler 微信登录处理
func wechatLoginHandler(openId string, clientId string, c *gin.Context) bool {
	if openId == "" {
		return false
	}

	authUserSocial, err := userSocialService.QueryUserSocialByOpenIDService(openId, enum.OAuthSourceWechat)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		tx := global.DB.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		uid := idgen.NextId()
		uidStr := strconv.FormatInt(uid, 10)
		avatar := utils.GenerateAvatar(uidStr)
		name := randomname.GenerateName()
		createUser := model.ScaAuthUser{
			UID:      uidStr,
			Username: openId,
			Avatar:   avatar,
			Nickname: name,
			Gender:   enum.Male,
		}

		// 添加用户
		addUser, worn := userService.AddUserService(createUser)
		if worn != nil {
			tx.Rollback()
			global.LOG.Error(worn)
			return false
		}

		wechat := enum.OAuthSourceWechat
		userSocial := model.ScaAuthUserSocial{
			UserID: uidStr,
			OpenID: openId,
			Source: wechat,
		}

		// 添加用户社交信息
		if wrong := userSocialService.AddUserSocialService(userSocial); wrong != nil {
			tx.Rollback()
			global.LOG.Error(wrong)
			return false
		}

		// 添加角色
		if _, err = global.Casbin.AddRoleForUser(uidStr, enum.User); err != nil {
			tx.Rollback()
			global.LOG.Error(err)
			return false
		}

		// 处理用户登录
		if res := handelUserLogin(addUser.UID, clientId, c); !res {
			tx.Rollback()
			return false
		}

		tx.Commit()
		return true
	} else {
		res := handelUserLogin(authUserSocial.UserID, clientId, c)
		return res
	}
}

// handelUserLogin 处理用户登录
func handelUserLogin(userId string, clientId string, c *gin.Context) bool {

	user := userService.QueryUserByUuidService(&userId)
	resultChan := make(chan bool, 1)

	go func() {
		accessToken, err := utils.GenerateAccessToken(utils.AccessJWTPayload{UserID: &userId})
		if err != nil {
			resultChan <- false
			return
		}
		refreshToken := utils.GenerateRefreshToken(utils.RefreshJWTPayload{UserID: &userId}, time.Hour*24*7)
		data := types.ResponseData{
			AccessToken: accessToken,
			UID:         &userId,
			Username:    user.Username,
			Nickname:    user.Nickname,
			Avatar:      user.Avatar,
			Status:      user.Status,
		}
		redisTokenData := types.RedisToken{
			AccessToken: accessToken,
			UID:         userId,
		}
		fail := redis.Set(constant.UserLoginTokenRedisKey+userId, redisTokenData, time.Hour*24*7).Err()
		if fail != nil {
			resultChan <- false
			return
		}
		responseData := result.Response{
			Data:    data,
			Message: "success",
			Code:    200,
			Success: true,
		}
		tokenData, err := json.Marshal(responseData)
		if err != nil {
			resultChan <- false
			return
		}
		sessionData := utils.SessionData{
			RefreshToken: refreshToken,
			UID:          userId,
		}
		wrong := utils.SetSession(c, constant.SessionKey, sessionData)
		if wrong != nil {
			resultChan <- false
			return
		}
		// gws方式发送消息
		err = qr_ws_controller.Handler.SendMessageToClient(clientId, tokenData)
		if err != nil {
			global.LOG.Error(err)
			resultChan <- false
			return
		}
		resultChan <- true
	}()

	return <-resultChan
}
