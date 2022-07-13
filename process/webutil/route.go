package webutil

import (
	"github.com/gin-gonic/gin"
	"github.com/soxft/busuanzi/app/controller"
	"github.com/soxft/busuanzi/app/middleware"
	"github.com/soxft/busuanzi/config"
)

func initRoute(r *gin.Engine) {
	r.LoadHTMLFiles(config.AssetsPath + "/index.html")
	{
		r.POST("/api", middleware.Identity(), controller.ApiHandler)
		r.GET("/ping", controller.PingHandler)

		static := r.Group("/")
		{
			static.Use(middleware.Cache())
			static.GET("/", controller.Index)
			static.StaticFile("/js", config.AssetsPath+"/busuanzi.js")
		}
		r.NoRoute(middleware.Cache(), controller.Index)
	}
}
