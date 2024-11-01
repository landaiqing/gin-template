package middleware

import (
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"

	"schisandra-cloud-album/common/constant"
	"schisandra-cloud-album/common/result"
	"schisandra-cloud-album/utils"
)

// SessionCheckMiddleware session检查中间件
func SessionCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := utils.GetSession(c, constant.SessionKey)
		if session == nil {
			result.FailWithCodeAndMessage(403, ginI18n.MustGetMessage(c, "AuthVerifyExpired"), c)
			c.Abort()
			return
		}

		userIdAny, exists := c.Get("userId")
		if !exists {
			result.FailWithCodeAndMessage(403, ginI18n.MustGetMessage(c, "AuthVerifyExpired"), c)
			c.Abort()
			return
		}
		userId, ok := userIdAny.(*string)
		if !ok {
			result.FailWithCodeAndMessage(403, ginI18n.MustGetMessage(c, "AuthVerifyExpired"), c)
			c.Abort()
			return
		}
		if *userId != *session.UID {
			result.FailWithCodeAndMessage(403, ginI18n.MustGetMessage(c, "AuthVerifyExpired"), c)
			c.Abort()
			return
		}
		c.Next()
	}
}
