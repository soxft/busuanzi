package main

import (
	"busuanzi/config"
	"busuanzi/redisHelper"
	"busuanzi/web"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {

	// init config
	config.Init()
	// init redis pool
	redisHelper.Pool = redisHelper.NewPool()

	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(web.AccessControl())

	r.LoadHTMLFiles("dist/index.html")
	r.StaticFile("/js", "dist/busuanzi.min.js")
	// router
	r.GET("/", web.Index)
	r.GET("/api", web.ApiHandler)
	r.GET("/ping", web.PingHandler)
	r.NoRoute(web.NoRouteHandler)

	// debug
	if !config.C.Web.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	// start server
	err := r.Run(config.C.Web.Address)
	if err != nil {
		fmt.Println("web服务启动失败 > {}", err)
	}
}
