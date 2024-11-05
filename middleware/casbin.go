package middleware

import (
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"

	"schisandra-cloud-album/common/result"
	"schisandra-cloud-album/global"
)

func CasbinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIdAny, exists := c.Get("user_id")
		if !exists {
			global.LOG.Error("casbin middleware: userId not found")
			result.FailWithMessage(ginI18n.MustGetMessage(c, "PermissionDenied"), c)
			c.Abort()
			return
		}
		userId, ok := userIdAny.(*string)
		if !ok {
			result.FailWithMessage(ginI18n.MustGetMessage(c, "PermissionDenied"), c)
			global.LOG.Error("casbin middleware: userId is not string")
			c.Abort()
			return
		}
		userIdStr := *userId
		method := c.Request.Method
		path := c.Request.URL.Path
		correct, err := global.Casbin.Enforce(userIdStr, path, method)
		if err != nil {
			result.FailWithMessage(ginI18n.MustGetMessage(c, "PermissionDenied"), c)
			global.LOG.Error("casbin middleware: ", err)
			c.Abort()
			return
		}
		if !correct {
			result.FailWithMessage(ginI18n.MustGetMessage(c, "PermissionDenied"), c)
			c.Abort()
			return
		}
		c.Next()
	}
}
