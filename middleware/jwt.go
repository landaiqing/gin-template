package middleware

import (
	"encoding/json"
	"strings"

	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"

	"schisandra-cloud-album/common/constant"
	"schisandra-cloud-album/common/redis"
	"schisandra-cloud-album/common/result"
	"schisandra-cloud-album/common/types"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/utils"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 默认Token放在请求头Authorization的Bearer中，并以空格隔开
		authHeader := c.GetHeader(global.CONFIG.JWT.HeaderKey)
		if authHeader == "" {
			result.FailWithCodeAndMessage(403, ginI18n.MustGetMessage(c, "AuthVerifyExpired"), c)
			c.Abort()
			return
		}
		headerPrefix := global.CONFIG.JWT.HeaderPrefix
		accessToken := strings.TrimPrefix(authHeader, headerPrefix+" ")

		if accessToken == "" {
			result.FailWithCodeAndMessage(403, ginI18n.MustGetMessage(c, "AuthVerifyExpired"), c)
			c.Abort()
			return
		}

		parseToken, isUpd, err := utils.ParseAccessToken(accessToken)
		if err != nil || !isUpd {
			result.FailWithCodeAndMessage(401, ginI18n.MustGetMessage(c, "AuthVerifyExpired"), c)
			c.Abort()
			return
		}
		token := redis.Get(constant.UserLoginTokenRedisKey + *parseToken.UserID).Val()
		if token == "" {
			result.FailWithCodeAndMessage(403, ginI18n.MustGetMessage(c, "AuthVerifyExpired"), c)
			c.Abort()
			return
		}
		tokenResult := types.RedisToken{}
		err = json.Unmarshal([]byte(token), &tokenResult)
		if err != nil {
			result.FailWithCodeAndMessage(403, ginI18n.MustGetMessage(c, "AuthVerifyExpired"), c)
			c.Abort()
			return
		}
		if tokenResult.AccessToken != accessToken {
			result.FailWithCodeAndMessage(403, ginI18n.MustGetMessage(c, "AuthVerifyExpired"), c)
			c.Abort()
			return
		}
		uid := utils.GetSession(c, constant.SessionKey).UID
		if uid != *parseToken.UserID {
			result.FailWithCodeAndMessage(403, ginI18n.MustGetMessage(c, "AuthVerifyExpired"), c)
			c.Abort()
			return
		}
		c.Set("user_id", parseToken.UserID)
		global.DB.Set("user_id", parseToken.UserID) // 全局变量中设置用户ID
		c.Next()
	}
}
