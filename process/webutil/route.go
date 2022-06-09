package webutil

import (
	"busuanzi/app/controller"
	"busuanzi/app/middleware"
	"github.com/gin-gonic/gin"
)

func initRoute(r *gin.Engine) {
	{
		r.POST("/api", controller.ApiHandler)
		r.GET("/ping", controller.PingHandler)

		static := r.Group("/")
		{
			static.Use(middleware.Cache())
			static.GET("/", controller.Index)
			static.StaticFile("/js", "dist/busuanzi.js")
		}
		r.NoRoute(middleware.Cache(), controller.Index)
	}
}
