package controller

import (
	"github.com/gin-gonic/gin"
)

func PingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"success": true,
		"message": "pong",
		"data":    gin.H{},
	})
}
