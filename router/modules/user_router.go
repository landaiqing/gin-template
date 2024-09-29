package modules

import (
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/controller"
)

var userApi = controller.Controller.UserController

// UserRouter 用户相关路由 有auth接口组需要token验证,没有auth接口组不需要token验证
func UserRouter(router *gin.RouterGroup) {
	userGroup := router.Group("user")
	{
		userGroup.POST("/login", userApi.AccountLogin)
		userGroup.POST("/phone_login", userApi.PhoneLogin)
		userGroup.POST("/reset_password", userApi.ResetPassword)
	}
	tokenGroup := router.Group("token")
	{
		tokenGroup.POST("/refresh", userApi.RefreshHandler)
	}
}

// UserRouterAuth 用户相关路由 有auth接口组需要token验证
func UserRouterAuth(router *gin.RouterGroup) {
	authGroup := router.Group("auth")
	{
		authGroup.GET("/user/list", userApi.GetUserList)
		authGroup.GET("/user/query_by_uuid", userApi.QueryUserByUuid)
	}

}
