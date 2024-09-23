package middleware

import (
	"encoding/json"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/common/constant"
	"schisandra-cloud-album/common/redis"
	"schisandra-cloud-album/common/result"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/utils"
	"strings"
)

type TokenData struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	ExpiresAt    int64   `json:"expires_at"`
	UID          *string `json:"uid"`
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 默认双Token放在请求头Authorization的Bearer中，并以空格隔开
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
		uid := c.GetHeader("X-UID")
		if uid == "" {
			result.FailWithCodeAndMessage(403, ginI18n.MustGetMessage(c, "AuthVerifyExpired"), c)
			c.Abort()
			return
		}
		if *parseToken.UserID != uid {
			result.FailWithCodeAndMessage(403, ginI18n.MustGetMessage(c, "AuthVerifyExpired"), c)
			c.Abort()
			return
		}
		token := redis.Get(constant.UserLoginTokenRedisKey + *parseToken.UserID).Val()
		if token == "" {
			result.FailWithCodeAndMessage(403, ginI18n.MustGetMessage(c, "AuthVerifyExpired"), c)
			c.Abort()
			return
		}
		tokenResult := TokenData{}
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
		c.Set("userId", parseToken.UserID)
		c.Next()
	}
}
