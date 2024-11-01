package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/soxft/busuanzi/config"
	"github.com/spf13/viper"
	"net/http"
	"strings"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		cors := viper.GetString("Web.Cors")

		// 多 cors 匹配 Failed to load resource: Access-Control-Allow-Origin cannot contain more than one origin.
		var corsPass = false
		var origin = c.Request.Header.Get("Origin")

		if strings.Contains(cors, ",") {
			// 多 Cors 匹配, 判断请求多域名是否在 cors 列表中
			for _, v := range strings.Split(cors, ",") {
				allow := strings.ToLower(strings.TrimSpace(v))

				if origin == allow {
					corsPass = true
					c.Header("Access-Control-Allow-Origin", allow)
					break
				}
			}
		} else {
			// 单 cors 匹配  // * 或者单域名
			if cors == "*" || origin == cors {
				corsPass = true
			}
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Server", "busuanzi-by-xcsoft/"+config.VERSION)
		if c.Request.Method == http.MethodOptions {
			c.Header("Access-Control-Allow-Methods", "GET, POST, HEAD, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "x-bsz-referer, Authorization, Content-Type")
			c.Header("Access-Control-Max-Age", "86400")
			if corsPass {
				c.AbortWithStatus(204)
			} else {
				c.AbortWithStatus(403)
			}
			return
		}
	}
}
