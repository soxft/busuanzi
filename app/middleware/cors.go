package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/soxft/busuanzi/config"
	"github.com/spf13/viper"
	"strings"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		cors := viper.GetString("Web.Cors")

		// 多 cors 匹配 Failed to load resource: Access-Control-Allow-Origin cannot contain more than one origin.
		if strings.Contains(cors, ",") {
			for _, v := range strings.Split(cors, ",") {
				allow := strings.ToLower(strings.TrimSpace(v))

				if c.Request.Header.Get("Origin") == allow {
					c.Header("Access-Control-Allow-Origin", allow)
					break
				}
			}
		} else {
			c.Header("Access-Control-Allow-Origin", viper.GetString("Web.Cors"))
		}

		c.Header("Server", "busuanzi-by-xcsoft/"+config.VERSION)
		if c.Request.Method == "OPTIONS" {
			c.Header("Access-Control-Allow-Methods", "GET, POST, HEAD, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "x-bsz-referer, Authorization")
			c.Header("Access-Control-Max-Age", "86400")
			c.AbortWithStatus(204)
			return
		}
	}
}
