package router

import (
	"github.com/gin-contrib/cors"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/common/result"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/middleware"
	"schisandra-cloud-album/router/modules"
	"time"
)

func InitRouter() *gin.Engine {
	gin.SetMode(global.CONFIG.System.Env)
	router := gin.Default()
	router.NoRoute(HandleNotFound)
	router.NoMethod(HandleNotFound)
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
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Accept-Language", "X-Sign", "X-Timestamp", "X-Nonce"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// 国际化设置
	router.Use(middleware.I18n())

	websocketRouter := router.Group("api")
	{
		modules.WebsocketRouter(websocketRouter) // 注册websocket路由
	}
	publicGroup := router.Group("api") // 不需要鉴权的路由组
	publicGroup.Use(middleware.SecurityHeaders())
	{
		modules.ClientRouter(publicGroup)  // 注册客户端路由
		modules.SwaggerRouter(publicGroup) // 注册swagger路由

		modules.OauthRouter(publicGroup)
		modules.CaptchaRouter(publicGroup) // 注册验证码路由
		modules.SmsRouter(publicGroup)     // 注册短信验证码路由
		modules.UserRouter(publicGroup)    // 注册鉴权路由
	}
	authGroup := router.Group("api") // 需要鉴权的路由组

	authGroup.Use(
		middleware.JWTAuthMiddleware(),
		middleware.CasbinMiddleware(),
		middleware.SecurityHeaders(),
	)
	{
		modules.UserRouterAuth(authGroup)   // 注册鉴权路由
		modules.RoleRouter(authGroup)       // 注册角色路由
		modules.PermissionRouter(authGroup) // 注册权限路由
	}

	return router
}

// HandleNotFound 404处理
func HandleNotFound(c *gin.Context) {
	result.FailWithCodeAndMessage(404, ginI18n.MustGetMessage(c, "404NotFound"), c)
	return
}
