package webutil

import (
	"busuanzi/app/controller"
	"busuanzi/app/middleware"
	"github.com/gin-gonic/gin"
)

func initRoute(r *gin.Engine) {
	{
		static := r.Group("/")
		{
			static.Use(middleware.Cache())
			static.GET("/", controller.Index)
			static.StaticFile("/js", "dist/busuanzi.js")
		}

		r.POST("/api", middleware.Identity(), controller.ApiHandler)
		r.GET("/ping", controller.PingHandler)
		r.NoRoute(controller.NoRouteHandler)
	}
}
