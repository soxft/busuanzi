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
	r.LoadHTMLFiles(config.DistPath + "/index.html")

	// routers
	initRoute(r)

	// start server
	log.SetOutput(os.Stdout)
	log.Println("server listening on port:", config.C.Web.Address)
	err := r.Run(config.C.Web.Address)
	if err != nil {
		log.Fatalf("we b服务启动失败: %s", err)
	}
}
