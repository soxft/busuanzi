package web

import "github.com/gin-gonic/gin"

func NoRouteHandler(c *gin.Context) {
	c.JSON(404, gin.H{
		"success": false,
		"message": "Page not found",
	})
}
