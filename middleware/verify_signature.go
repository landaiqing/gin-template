package middleware

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"schisandra-cloud-album/common/constant"
	"schisandra-cloud-album/common/redis"
	"schisandra-cloud-album/common/result"
	"schisandra-cloud-album/global"
	"strconv"
	"time"
)

func VerifySignature() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 仅处理 POST 请求
		if c.Request.Method != http.MethodPost {
			c.Next()
			return
		}
		// 从请求头获取签名和时间戳
		signature := c.GetHeader("X-Sign")
		timestamp := c.GetHeader("X-Timestamp")
		nonce := c.GetHeader("X-Nonce")
		secretKey := global.CONFIG.Encrypt.Key

		// 检查时间戳是否过期，这里设置为5分钟过期
		if timestamp == "" || time.Since(parseTimestamp(timestamp)) > 5*time.Minute {
			result.FailWithMessage(ginI18n.MustGetMessage(c, "RequestVerifyError"), c)
			c.Abort()
			return
		}

		// 检查 nonce 是否已经被使用
		if data := redis.Get(constant.SystemApiNonceRedisKey + nonce).Val(); data != "" {
			result.FailWithMessage(ginI18n.MustGetMessage(c, "RequestVerifyError"), c)
			c.Abort()
			return
		}

		// 记录 nonce 到 Redis 中，并设置过期时间为 5 分钟
		if err := redis.Set(constant.SystemApiNonceRedisKey+nonce, true, 5*time.Minute).Err(); err != nil {
			global.LOG.Error(err.Error())
			c.Abort()
			return
		}

		// 获取请求方法和请求体
		var payload string
		if c.Request.Method == http.MethodPost {
			body, err := io.ReadAll(c.Request.Body)
			if err != nil {
				result.FailWithMessage(ginI18n.MustGetMessage(c, "RequestReadError"), c)
				c.Abort()
				return
			}
			payload = string(body)
			// 重新设置请求体，以便后续处理中可以再次读取
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		// 创建待签名字符串
		baseString := c.Request.Method + ":" + payload + ":" + timestamp + ":" + nonce + ":" + secretKey

		// 生成 MD5 签名
		h := md5.New()
		h.Write([]byte(baseString))
		expectedSignature := hex.EncodeToString(h.Sum(nil))

		// 验证签名
		if signature != expectedSignature {
			result.FailWithMessage(ginI18n.MustGetMessage(c, "RequestVerifyError"), c)
			c.Abort()
			return
		}
		// 继续处理请求
		c.Next()
	}
}

// 辅助函数：解析时间戳
func parseTimestamp(ts string) time.Time {
	t, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return time.Time{} // 解析错误返回零时间
	}
	return time.Unix(t/1000, 0) // 假设时间戳是毫秒
}
