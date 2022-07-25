package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/soxft/busuanzi/library/jwtutil"
	"github.com/soxft/busuanzi/library/tool"
	"strings"
)

func Identity() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Set-Bsz-Identity")

		// token
		var userIdentity string
		tokenTmp := c.Request.Header.Get("Authorization")

		if tokenTmp == "" {
			// generate jwt token
			userIdentity = getUserIdentity(c)
			setBszIdentity(c, userIdentity)
		} else {
			token := strings.Replace(tokenTmp, "Bearer ", "", -1)
			// check if token is illegal
			if userIdentity = jwtutil.Check(token); userIdentity == "" {
				// fake data, regenerate jwt token
				userIdentity = getUserIdentity(c)
				setBszIdentity(c, userIdentity)
			}
		}
		c.Set("user_identity", userIdentity)
		c.Next()
	}
}

func setBszIdentity(c *gin.Context, userIdentity string) {
	uid := jwtutil.Generate(userIdentity)
	c.Writer.Header().Set("Set-Bsz-Identity", uid)
}

func getUserIdentity(c *gin.Context) string {
	return tool.Md5(c.ClientIP() + c.Request.UserAgent())
}
