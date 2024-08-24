package middleware

import (
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/global"
)

func CasbinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			global.LOG.Error("casbin middleware: userId not found")
			c.Abort()
			return
		}
		method := c.Request.Method
		path := c.Request.URL.Path
		ok, err := global.Casbin.Enforce(userId.(string), path, method)
		if err != nil {
			global.LOG.Error("casbin middleware: ", err)
			c.Abort()
			return
		}
		if !ok {
			c.Abort()
			return
		}
		c.Next()
	}
}
