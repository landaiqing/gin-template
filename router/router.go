package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/api"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/middleware"
	"schisandra-cloud-album/router/modules"
	"time"
)

var oauth = api.Api.OAuthApi

func InitRouter() *gin.Engine {
	gin.SetMode(global.CONFIG.System.Env)
	router := gin.Default()
	err := router.SetTrustedProxies([]string{global.CONFIG.System.Ip})
	if err != nil {
		global.LOG.Error(err)
		return nil
	}
	router.Use(middleware.RateLimitMiddleware(time.Millisecond*100, 20)) // 限流中间件
	// 跨域设置
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{global.CONFIG.System.Web},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-CSRF-Token", "Accept-Language", "X-Request-Id"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// 国际化设置
	router.Use(middleware.I18n())

	publicGroup := router.Group("api") // 不需要鉴权的路由组
	authGroup := router.Group("api")   // 需要鉴权的路由组
	authGroup.Use(
		middleware.JWTAuthMiddleware(),
		middleware.CasbinMiddleware(),
		middleware.CheckClientMiddleware())

	checkClientGroup := router.Group("api") // 需要检查客户端的路由组

	checkClientGroup.Use(middleware.CheckClientMiddleware())

	modules.ClientRouter(publicGroup)    // 注册客户端路由
	modules.SwaggerRouter(publicGroup)   // 注册swagger路由
	modules.WebsocketRouter(publicGroup) // 注册websocket路由

	modules.CaptchaRouter(checkClientGroup) // 注册验证码路由
	modules.SmsRouter(checkClientGroup)     // 注册短信验证码路由
	modules.OauthRouter(checkClientGroup)   // 注册oauth路由
	modules.UserRouter(checkClientGroup)    // 注册鉴权路由

	modules.UserRouterAuth(authGroup)   // 注册鉴权路由
	modules.RoleRouter(authGroup)       // 注册角色路由
	modules.PermissionRouter(authGroup) // 注册权限路由

	return router
}
