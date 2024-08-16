package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/soxft/busuanzi/library/jwtutil"
	"github.com/soxft/busuanzi/library/tool"
	"net/http"
)

func Identity() gin.HandlerFunc {
	return func(c *gin.Context) {
		//c.Writer.Header().Set("Access-Control-Expose-Headers", "Set-Bsz-Identity")

		// token
		var userIdentity string
		tokenTmp, err := c.Request.Cookie("bsz_id")

		if errors.Is(err, http.ErrNoCookie) {
			// generate jwt token
			userIdentity = getUserIdentity(c)
			setBszIdentity(c, userIdentity)
		} else {
			// check if token is illegal
			if userIdentity = jwtutil.Check(tokenTmp.Value); userIdentity == "" {
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
	c.SetCookie("bsz_id", uid, 86400, "/", "", false, true)
}

func getUserIdentity(c *gin.Context) string {
	return tool.Md5(c.ClientIP() + c.Request.UserAgent())
}
