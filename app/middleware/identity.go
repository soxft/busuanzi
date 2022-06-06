package middleware

import (
	"busuanzi/tool"
	"github.com/gin-gonic/gin"
	"strings"
)

func Identity() gin.HandlerFunc {
	return func(c *gin.Context) {
		// token
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			token = tool.Md5(c.ClientIP()) + "." + tool.Md5(c.Request.UserAgent())
			c.Writer.Header().Set("Access-Control-Expose-Headers", "Set-Bsz-Identity")
			c.Writer.Header().Set("Set-Bsz-Identity", token)
		} else {
			token = strings.Replace(token, "Bearer ", "", -1)
		}
		c.Set("user_identity", tool.Md5(token))
		c.Next()
	}
}
