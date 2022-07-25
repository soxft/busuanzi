package webutil

import (
	"github.com/gin-gonic/gin"
	"github.com/soxft/busuanzi/app/controller"
	"github.com/soxft/busuanzi/app/middleware"
)

func initRoute(r *gin.Engine) {
	{
		r.POST("/api", middleware.Identity(), controller.ApiHandler)
		r.GET("/ping", controller.PingHandler)
	}
}
