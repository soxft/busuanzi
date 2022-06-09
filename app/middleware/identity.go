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
		var userIdentity string
		tokenTmp := c.Request.Header.Get("Authorization")

		if tokenTmp == "" {
			// generate jwt token
			userIdentity = tool.Md5(c.ClientIP()) + tool.Md5(c.Request.UserAgent())
			setBszIdentity(c, userIdentity)
		} else {
			token := strings.Replace(tokenTmp, "Bearer ", "", -1)

			if userIdentity = jwtutil.Check(token); userIdentity == "" {
				// fake data, regenerate jwt token
				userIdentity = tool.Md5(c.ClientIP()) + tool.Md5(c.Request.UserAgent())
				setBszIdentity(c, userIdentity)
			}
		}
		c.Set("user_identity", tool.Md5(userIdentity))
		c.Next()
	}
}

func setBszIdentity(c *gin.Context, userIdentity string) {
	uid := jwtutil.Generate(userIdentity)
	c.Writer.Header().Set("Access-Control-Expose-Headers", "Set-Bsz-Identity")
	c.Writer.Header().Set("Set-Bsz-Identity", uid)
}
