package router

import (
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/router/modules"
)

func InitRouter() *gin.Engine {
	gin.SetMode(global.CONFIG.System.Env)
	Router := gin.Default()
	PublicGroup := Router.Group("api")
	modules.AuthRouter(PublicGroup)
	return Router
}
