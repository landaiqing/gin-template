package middleware

import (
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/messages"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/common/result"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/utils"
	"time"
)

// ExceptionNotification 异常通知中间件
func ExceptionNotification() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				openID := global.CONFIG.Wechat.OpenID
				content := `
				系统异常通知：
				请求时间：` + time.Now().Format("2006-01-02 15:04:05") + `
				请求IP：` + utils.GetClientIP(c) + `
				请求地址：` + c.Request.URL.String() + `
				请求方法：` + c.Request.Method + `
				请求参数：` + c.Request.Form.Encode() + `
				错误信息：` + err.(error).Error() + `
`
				messages.NewRaw(`
{
    "touser":"` + openID + `",
    "msgtype":"text",
    "text":{"content":"` + content + `"}"}
}
`)
				result.FailWithMessage(ginI18n.MustGetMessage(c, "SystemError"), c)
			}
		}()
		c.Next()
	}
}
