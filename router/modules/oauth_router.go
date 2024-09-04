package modules

import (
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/api"
)

var oauth = api.Api.OAuthApi

func OauthRouter(router *gin.RouterGroup) {
	group := router.Group("/oauth")
	{
		wechatRouter := group.Group("/wechat")
		{
			wechatRouter.GET("/get_temp_qrcode", oauth.GetTempQrCode)
			//wechatRouter.GET("/callback", oauth.CallbackVerify)
			wechatRouter.POST("/callback", oauth.CallbackNotify)
		}
		githubRouter := group.Group("/github")
		{
			githubRouter.GET("/get_url", oauth.GetRedirectUrl)
			githubRouter.GET("/callback", oauth.Callback)
		}
		giteeRouter := group.Group("/gitee")
		{
			giteeRouter.GET("/get_url", oauth.GetGiteeRedirectUrl)
			giteeRouter.GET("/callback", oauth.GiteeCallback)
		}
		qqRouter := group.Group("/qq")
		{
			qqRouter.GET("/get_url", oauth.GetQQRedirectUrl)
			qqRouter.GET("/callback", oauth.QQCallback)
		}
		group.GET("/get_device", oauth.GetUserLoginDevice)
	}
}
