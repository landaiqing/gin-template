package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/middleware"
	"schisandra-cloud-album/router/modules"
	"time"
)

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
		AllowOrigins:     []string{global.CONFIG.System.WebURL()},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Accept-Language", "X-Sign", "X-Timestamp", "X-Nonce", "X-UID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// 国际化设置
	router.Use(middleware.I18n())

	noMiddlewareRouter := router.Group("api") // 不需要中间件的路由组
	{
		modules.WebsocketRouter(noMiddlewareRouter) // 注册websocket路由
		modules.OauthRouter(noMiddlewareRouter)     // 注册oauth路由
	}
	publicGroup := router.Group("api") // 不需要鉴权的路由组
	publicGroup.Use(middleware.SecurityHeaders())
	{
		modules.ClientRouter(publicGroup)  // 注册客户端路由
		modules.SwaggerRouter(publicGroup) // 注册swagger路由
		modules.OauthRouterAuth(publicGroup)
		modules.CaptchaRouter(publicGroup) // 注册验证码路由
		modules.SmsRouter(publicGroup)     // 注册短信验证码路由
		modules.UserRouter(publicGroup)    // 注册鉴权路由
	}

	authGroup := router.Group("api") // 需要鉴权的路由组
	authGroup.Use(
		middleware.SecurityHeaders(),
		middleware.JWTAuthMiddleware(),
		middleware.CasbinMiddleware(),
	)
	{
		modules.UserRouterAuth(authGroup)    // 注册鉴权路由
		modules.RoleRouter(authGroup)        // 注册角色路由
		modules.PermissionRouter(authGroup)  // 注册权限路由
		modules.CommentRouter(authGroup)     // 注册评论路由
		modules.CaptchaRouterAuth(authGroup) // 注册验证码路由
	}

	return router
}
