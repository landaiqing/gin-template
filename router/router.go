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
	publicGroup := router.Group("api")
	// 跨域设置
	publicGroup.Use(cors.New(cors.Config{
		AllowOrigins:     []string{global.CONFIG.System.Web},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-CSRF-Token", "Accept-Language"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// 国际化设置
	publicGroup.Use(middleware.I18n())

	modules.SwaggerRouter(router)        // 注册swagger路由
	modules.UserRouter(publicGroup)      // 注册鉴权路由
	modules.CaptchaRouter(publicGroup)   // 注册验证码路由
	modules.SmsRouter(publicGroup)       // 注册短信验证码路由
	modules.OauthRouter(publicGroup)     // 注册oauth路由
	modules.WebsocketRouter(publicGroup) // 注册websocket路由
	modules.RoleRouter(publicGroup)      // 注册角色路由
	return router
}
