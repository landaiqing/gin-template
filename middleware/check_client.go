package middleware

import (
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/common/constant"
	"schisandra-cloud-album/common/redis"
	"schisandra-cloud-album/common/result"
	"schisandra-cloud-album/utils"
)

// CheckClientMiddleware 检查客户端请求是否合法
func CheckClientMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Request-Id")
		if id == "" {
			result.FailWithMessage(ginI18n.MustGetMessage(c, "AuthVerifyFailed"), c)
			c.Abort()
			return
		}

		ip := utils.GetClientIP(c)
		clientId := redis.Get(constant.UserLoginClientRedisKey + ip).Val()
		if clientId == "" || clientId != id {
			result.FailWithMessage(ginI18n.MustGetMessage(c, "AuthVerifyFailed"), c)
			c.Abort()
			return
		}
		c.Next()
	}
}
