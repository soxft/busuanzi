package middleware

import (
	"busuanzi/library/jwtutil"
	"busuanzi/library/tool"
	"github.com/gin-gonic/gin"
	"strings"
)

func Identity() gin.HandlerFunc {
	return func(c *gin.Context) {
		// token
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			// generate jwt token
			userIdentity := tool.Md5(c.ClientIP()) + "." + tool.Md5(c.Request.UserAgent())
			token := jwtutil.Generate(userIdentity)

			c.Writer.Header().Set("Access-Control-Expose-Headers", "Set-Bsz-Identity")
			c.Writer.Header().Set("Set-Bsz-Identity", token)
		} else {
			token = strings.Replace(token, "Bearer ", "", -1)
		}
		c.Set("user_identity", tool.Md5(token))
		c.Next()
	}
}
