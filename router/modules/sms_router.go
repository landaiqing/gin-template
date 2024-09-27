package modules

import (
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/api"
)

var smsApi = api.Api.SmsApi

func SmsRouter(router *gin.RouterGroup) {
	group := router.Group("/sms")
	group.POST("/ali/send", smsApi.SendMessageByAli)
	group.POST("/smsbao/send", smsApi.SendMessageBySmsBao)
	group.POST("/test/send", smsApi.SendMessageTest)
}
