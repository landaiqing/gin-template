package modules

import (
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/api"
)

var oauth = api.Api.OAuthApi

func OauthRouter(router *gin.RouterGroup) {
	group := router.Group("/oauth")
	group.GET("/generate_client_id", oauth.GenerateClientId)
	group.GET("/get_temp_qrcode", oauth.GetTempQrCode)
	//group.GET("/callback", oauth.CallbackVerify)
	group.POST("/callback", oauth.CallbackNotify)
}
