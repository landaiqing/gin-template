package modules

import "github.com/gin-gonic/gin"

func AuthRouter(router *gin.RouterGroup) {
	group := router.Group("auth")
	group.GET("/user", func(c *gin.Context) {

	})
}
