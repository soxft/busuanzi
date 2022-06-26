package webutil

import (
	"busuanzi/app/controller"
	"busuanzi/app/middleware"
	"busuanzi/config"
	"github.com/gin-gonic/gin"
)

func initRoute(r *gin.Engine) {
	{
		r.POST("/api", middleware.Identity(), controller.ApiHandler)
		r.GET("/ping", controller.PingHandler)

		static := r.Group("/")
		{
			static.Use(middleware.Cache())
			static.GET("/", controller.Index)
			static.StaticFile("/js", config.DistPath+"/busuanzi.js")
		}
		r.NoRoute(middleware.Cache(), controller.Index)
	}
}
