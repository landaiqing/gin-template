package router

import (
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/router/modules"
)

func InitRouter() *gin.Engine {
	gin.SetMode(global.CONFIG.System.Env)
	router := gin.Default()
	err := router.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		global.LOG.Error(err)
		return nil
	}
	PublicGroup := router.Group("api")
	modules.SwaggerRouter(router)   // 注册swagger路由
	modules.AuthRouter(PublicGroup) // 注册鉴权路由
	return router
}
