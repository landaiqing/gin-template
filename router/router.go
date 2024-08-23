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
	err := router.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		global.LOG.Error(err)
		return nil
	}
	publicGroup := router.Group("api")
	// 跨域设置
	publicGroup.Use(cors.Default()).Use(middleware.RateLimitMiddleware(time.Millisecond*100, 20))
	// 国际化设置
	publicGroup.Use(middleware.I18n())

	modules.SwaggerRouter(router)        // 注册swagger路由
	modules.UserRouter(publicGroup)      // 注册鉴权路由
	modules.CaptchaRouter(publicGroup)   // 注册验证码路由
	modules.SmsRouter(publicGroup)       // 注册短信验证码路由
	modules.OauthRouter(publicGroup)     // 注册oauth路由
	modules.WebsocketRouter(publicGroup) // 注册websocket路由
	return router
}
