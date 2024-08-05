package modules

import (
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/api"
)

var authApi = api.Api.AuthApi

func AuthRouter(router *gin.RouterGroup) {
	group := router.Group("auth")
	group.GET("/user/List", authApi.GetUserList)
	group.GET("/user/query_by_username", authApi.QueryUserByUsername)
	group.GET("/user/query_by_uuid", authApi.QueryUserByUuid)
	group.DELETE("/user/delete", authApi.DeleteUser)
	group.GET("/user/query_by_phone", authApi.QueryUserByPhone)
}
