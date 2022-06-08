package webutil

import (
	"busuanzi/app/middleware"
	"busuanzi/config"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func Init() {
	if !config.Web.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// middleware
	if config.Web.Log {
		r.Use(gin.Logger())
	}
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())
	r.LoadHTMLFiles("dist/index.html")

	// routers
	initRoute(r)

	// start server
	log.SetOutput(os.Stdout)
	log.Println("server listen on port:", config.Web.Address)
	err := r.Run(config.Web.Address)
	if err != nil {
		log.Fatalf("web服务启动失败: %s", err)
	}
}
