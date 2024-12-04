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

	r.GET("/", NoRouteHandler)

	js := r.Group("js")
	{
		js.Use(middleware.Cache())
		js.StaticFile("", config.DistPath+"/busuanzi.js")
		js.StaticFile("/jsonp", config.DistPath+"/busuanzi.jsonp.js")
		js.StaticFile("/lite", config.DistPath+"/busuanzi.lite.js")
		js.StaticFile("/lite_pjax", config.DistPath+"/busuanzi.pjax.lite.js")
	}

	r.NoRoute(NoRouteHandler)
}

func NoRouteHandler(c *gin.Context) {
	c.Redirect(302, "https://busuanzi.9420.ltd")
}
