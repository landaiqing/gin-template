package modules

import (
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/common/result"
)

func AuthRouter(router *gin.RouterGroup) {
	group := router.Group("auth")
	group.GET("/user", func(c *gin.Context) {
		result.FailWithCode(result.SystemError, c)
	})
}
