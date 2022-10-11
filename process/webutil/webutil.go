package webutil

import (
	"github.com/gin-gonic/gin"
	"github.com/soxft/busuanzi/app/middleware"
	"github.com/soxft/busuanzi/config"
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
	log.Println("server listening on port:", config.Web.Address)
	err := r.Run(config.Web.Address)
	if err != nil {
		log.Fatalf("we b服务启动失败: %s", err)
	}
}
