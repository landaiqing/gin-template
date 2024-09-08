package middleware

import (
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"schisandra-cloud-album/common/result"
	"time"
)

func RateLimitMiddleware(fillInterval time.Duration, cap int64) func(c *gin.Context) {
	bucket := ratelimit.NewBucket(fillInterval, cap)
	return func(c *gin.Context) {
		// 如果取不到令牌就中断本次请求返回 rate limit...
		if bucket.TakeAvailable(1) < 1 {
			result.FailWithMessage(ginI18n.MustGetMessage(c, "RequestLimit"), c)
			c.Abort()
			return
		}
		c.Next()
	}
}
