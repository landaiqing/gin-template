package middleware

import (
	"github.com/gin-gonic/gin"
)

func ValidateSignMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Next()
	}
}
