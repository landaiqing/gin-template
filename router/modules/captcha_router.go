package modules

import (
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/api"
)

var captchaApi = api.Api.CaptchaApi

func CaptchaRouter(router *gin.RouterGroup) {
	group := router.Group("/captcha")
	group.GET("/rotate/get", captchaApi.GenerateRotateCaptcha)
	group.POST("/rotate/check", captchaApi.CheckRotateData)
}
