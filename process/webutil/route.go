package webutil

import (
	"github.com/gin-gonic/gin"
	"github.com/soxft/busuanzi/app/controller"
	"github.com/soxft/busuanzi/app/middleware"
	"github.com/soxft/busuanzi/config"
)

func initRoute(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.Use(middleware.Identity())
		api.POST("", controller.ApiHandler)
		api.GET("", controller.GetHandler)
		api.PUT("", controller.PutHandler)
	}

	// 仿原版 Jsonp 请求模式
	r.GET("/jsonp", controller.JsonpHandler)
	r.GET("/ping", controller.PingHandler)

	static := r.Group("/")
	{
		static.Use(middleware.Cache())
		static.GET("/", controller.Index)
	}

	js := r.Group("js")
	{
		static.Use(middleware.Cache())
		js.StaticFile("", config.DistPath+"/busuanzi.js")
		js.StaticFile("/jsonp", config.DistPath+"/busuanzi.jsonp.js")
		js.StaticFile("/lite", config.DistPath+"/busuanzi.lite.js")
		js.StaticFile("/lite_pjax", config.DistPath+"/busuanzi.pjax.lite.js")
	}

	r.NoRoute(middleware.Cache(), controller.Index)
}
