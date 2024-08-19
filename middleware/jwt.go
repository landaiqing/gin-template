package middleware

import (
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/common/result"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/utils"
	"strings"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 默认双Token放在请求头Authorization的Bearer中，并以空格隔开
		authHeader := c.GetHeader(global.CONFIG.JWT.HeaderKey)
		if authHeader == "" {
			c.Abort()
			result.FailWithMessage(ginI18n.MustGetMessage(c, "AuthVerifyFailed"), c)
			return
		}
		headerPrefix := global.CONFIG.JWT.HeaderPrefix
		accessToken := strings.TrimPrefix(authHeader, headerPrefix+" ")

		if accessToken == "" {
			c.Abort()
			result.FailWithMessage(ginI18n.MustGetMessage(c, "AuthVerifyFailed"), c)
			return
		}
		parseToken, isUpd, err := utils.ParseAccessToken(accessToken)
		if err != nil || !isUpd {
			c.Abort()
			result.FailWithCodeAndMessage(401, ginI18n.MustGetMessage(c, "AuthVerifyExpired"), c)
			return
		}
		c.Set("userId", parseToken.UserID)
		c.Next()
	}
}
