package middleware

import (
	"github.com/gin-gonic/gin"
)

// Cache set CacheControl
func Cache() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "max-age=43200")
	}
}
