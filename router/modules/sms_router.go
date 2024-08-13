package modules

import (
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/api"
)

var smsApi = api.Api.SmsApi

func SmsRouter(router *gin.RouterGroup) {
	group := router.Group("/sms")
	group.GET("/ali/send", smsApi.SendMessageByAli)
	group.GET("/smsbao/send", smsApi.SendMessageBySmsBao)
	group.GET("/test/send", smsApi.SendMessageTest)
}
