package modules

import (
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/api"
)

var authApi = api.Api.AuthApi

func AuthRouter(router *gin.RouterGroup) {
	group := router.Group("auth")
	group.GET("/user/List", authApi.GetUserList)
}
