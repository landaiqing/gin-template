package modules

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"schisandra-cloud-album/docs"
	"schisandra-cloud-album/global"
)

func SwaggerRouter(router *gin.RouterGroup) {
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Description = global.CONFIG.Swagger.Description
	router.GET("/doc/*any", gin.BasicAuth(gin.Accounts{
		global.CONFIG.Swagger.User: global.CONFIG.Swagger.Password,
	}), ginSwagger.WrapHandler(swaggerFiles.Handler, func(config *ginSwagger.Config) {
		config.Title = global.CONFIG.Swagger.Title
	}))
}
