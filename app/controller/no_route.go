package controller

import (
	"github.com/gin-gonic/gin"
)

func NoRouteHandler(c *gin.Context) {
	c.JSON(404, gin.H{
		"success": false,
		"message": "route not found",
		"data":    gin.H{},
	})
}
