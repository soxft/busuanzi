package web

import (
	"busuanzi/config"
	"github.com/gin-gonic/gin"
)

func AccessControl() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", config.C.Web.AcAo)
		c.Header("Server", "busuanzi-by-xcsoft/0.1")
	}
}
