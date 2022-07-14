package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/soxft/busuanzi/config"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", config.Web.Cors)
		c.Header("Server", "busuanzi-by-xcsoft/2.7.4")
		if c.Request.Method == "OPTIONS" {
			c.Header("Access-Control-Allow-Methods", "GET, POST, HEAD, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "x-bsz-referer, Authorization")
			c.Header("Access-Control-Max-Age", "86400")
			c.AbortWithStatus(204)
			return
		}
	}
}
