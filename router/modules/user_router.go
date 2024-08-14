package modules

import (
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/api"
	"schisandra-cloud-album/middleware"
)

var userApi = api.Api.UserApi

// UserRouter 用户相关路由 有auth接口组需要token验证,没有auth接口组不需要token验证
func UserRouter(router *gin.RouterGroup) {
	userGroup := router.Group("user")
	{
		userGroup.POST("/login", userApi.AccountLogin)
		userGroup.POST("/phone_login", userApi.PhoneLogin)
	}
	authGroup := router.Group("auth").Use(middleware.JWTAuthMiddleware())
	{
		authGroup.GET("/user/list", userApi.GetUserList)
		authGroup.GET("/user/query_by_uuid", userApi.QueryUserByUuid)

	}
	tokenGroup := router.Group("token")
	{
		tokenGroup.POST("/refresh", userApi.RefreshHandler)
	}

}
