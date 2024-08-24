package utils

import "github.com/gin-gonic/gin"

// GetClientIP 工具函数，获取客户端IP
func GetClientIP(c *gin.Context) string {
	ip := c.GetHeader("X-Real-IP")
	if ip == "" {
		ip = c.GetHeader("X-Forwarded-For")
	}
	if ip == "" {
		ip = c.ClientIP()
	}
	return ip
}
