package main

import (
	"busuanzi/config"
	"busuanzi/web"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r.GET("/api", web.ApiHandler)
	err := r.Run(config.WebAddr)
	if err != nil {
		fmt.Println("web服务启动失败 > {}", err)
	}
}
