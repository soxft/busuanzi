package middleware

import (
	"github.com/gin-gonic/gin"
)

// Cache set CacheControl
func Cache() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "public, max-age=86400")
		c.Next()
	}
}
