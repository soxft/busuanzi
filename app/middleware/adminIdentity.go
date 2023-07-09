package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/soxft/busuanzi/config"
)

func AdminIdentity() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查 POST 中的 password 是否正确
		password := c.PostForm("password")
		if password != config.Admin.Password {
			c.JSON(200, gin.H{
				"success": false,
				"message": "password is incorrect",
				"data":    gin.H{},
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
