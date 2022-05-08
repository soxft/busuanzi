package main

import (
	"busuanzi/config"
	"busuanzi/web"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// debug
	if !config.C.Web.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// middleware
	if config.C.Web.Log {
		r.Use(gin.Logger())
	}
	r.Use(gin.Recovery())
	r.Use(web.AccessControl())

	// web
	r.LoadHTMLFiles("dist/index.html")
	r.StaticFile("/js", "dist/busuanzi.js")

	// router
	r.GET("/", web.Index)
	r.GET("/ping", web.PingHandler)
	r.GET("/api", web.ApiHandler)
	r.OPTIONS("/api", web.ApiHandler)
	r.NoRoute(web.NoRouteHandler)

	// start server
	log.Println("server listen on port:", config.C.Web.Address)
	err := r.Run(config.C.Web.Address)
	if err != nil {
		log.Fatalf("web服务启动失败: %s", err)
	}
}
