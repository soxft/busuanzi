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
	r.GET("/api", web.ApiHandler)
	r.GET("/ping", web.PingHandler)
	r.NoRoute(web.NoRouteHandler)

	// start server
	fmt.Println("server listen on port:", config.C.Web.Address)
	err := r.Run(config.C.Web.Address)
	if err != nil {
		fmt.Println("web服务启动失败 > {}", err)
	}
}
