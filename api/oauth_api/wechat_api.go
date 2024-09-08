package oauth_api

import (
	"encoding/json"
	"errors"
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
	"schisandra-cloud-album/api/user_api/dto"
	"schisandra-cloud-album/api/websocket_api"
	"schisandra-cloud-album/common/constant"
	"schisandra-cloud-album/common/enum"
	"schisandra-cloud-album/common/redis"
	"schisandra-cloud-album/common/result"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/model"
	"schisandra-cloud-album/utils"
	"strconv"
	"strings"
	"time"
)

// CallbackNotify 微信回调
// @Summary 微信回调
// @Tags 微信公众号
// @Description 微信回调
// @Produce json
// @Router /api/oauth/callback_notify [POST]
func (OAuthAPI) CallbackNotify(c *gin.Context) {
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
				return messages.NewText("ok")

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
// @Tags 微信公众号
// @Description 获取临时二维码
// @Produce json
// @Param client_id query string true "客户端ID"
// @Router /api/oauth/get_temp_qrcode [get]
func (OAuthAPI) GetTempQrCode(c *gin.Context) {
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
func wechatLoginHandler(openId string, clientId string) bool {
	if openId == "" {
		return false
	}
	authUserSocial, err := userSocialService.QueryUserSocialByOpenID(openId, enum.OAuthSourceWechat)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		tx := global.DB.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		uid := idgen.NextId()
		uidStr := strconv.FormatInt(uid, 10)
		createUser := model.ScaAuthUser{
			UID:      &uidStr,
			Username: &openId,
		}

		// 异步添加用户
		addUserChan := make(chan *model.ScaAuthUser, 1)
		errChan := make(chan error, 1)
		go func() {
			addUser, err := userService.AddUser(createUser)
			if err != nil {
				errChan <- err
				return
			}
			addUserChan <- &addUser
		}()

		var addUser *model.ScaAuthUser
		select {
		case addUser = <-addUserChan:
		case err := <-errChan:
			tx.Rollback()
			global.LOG.Error(err)
			return false
		}

		wechat := enum.OAuthSourceWechat
		userSocial := model.ScaAuthUserSocial{
			UserID: &uidStr,
			OpenID: &openId,
			Source: &wechat,
		}

		// 异步添加用户社交信息
		wrongChan := make(chan error, 1)
		go func() {
			wrong := userSocialService.AddUserSocial(userSocial)
			wrongChan <- wrong
		}()

		select {
		case wrong := <-wrongChan:
			if wrong != nil {
				tx.Rollback()
				global.LOG.Error(wrong)
				return false
			}
		}

		// 异步添加角色
		roleErrChan := make(chan error, 1)
		go func() {
			_, err := global.Casbin.AddRoleForUser(uidStr, enum.User)
			roleErrChan <- err
		}()

		select {
		case err := <-roleErrChan:
			if err != nil {
				tx.Rollback()
				global.LOG.Error(err)
				return false
			}
		}

		// 异步处理用户登录
		resChan := make(chan bool, 1)
		go func() {
			res := handelUserLogin(*addUser.UID, clientId)
			resChan <- res
		}()

		select {
		case res := <-resChan:
			if !res {
				tx.Rollback()
				return false
			}
		}
		tx.Commit()
		return true
	} else {
		res := handelUserLogin(*authUserSocial.UserID, clientId)
		if !res {
			return false
		}
		return true
	}
}

// handelUserLogin 处理用户登录
func handelUserLogin(userId string, clientId string) bool {
	resultChan := make(chan bool, 1)

	go func() {
		accessToken, err := utils.GenerateAccessToken(utils.AccessJWTPayload{UserID: &userId})
		if err != nil {
			resultChan <- false
			return
		}
		refreshToken, expiresAt := utils.GenerateRefreshToken(utils.RefreshJWTPayload{UserID: &userId}, time.Hour*24*7)
		data := dto.ResponseData{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresAt:    expiresAt,
			UID:          &userId,
		}
		fail := redis.Set(constant.UserLoginTokenRedisKey+userId, data, time.Hour*24*7).Err()
		if fail != nil {
			resultChan <- false
			return
		}
		responseData := map[string]interface{}{
			"code":    0,
			"message": "success",
			"data":    data,
			"success": true,
		}
		tokenData, err := json.Marshal(responseData)
		if err != nil {
			resultChan <- false
			return
		}
		// gws方式发送消息
		err = websocket_api.Handler.SendMessageToClient(clientId, tokenData)
		if err != nil {
			global.LOG.Error(err)
			resultChan <- false
			return
		}
		resultChan <- true
	}()

	return <-resultChan
}
