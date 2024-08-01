package webutil

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/soxft/busuanzi/app/middleware"
	"github.com/soxft/busuanzi/config"
	"github.com/spf13/viper"
	"log"
	"net/url"
	"os"
)

func Init() {
	if viper.GetBool("web.debug") == false {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// middleware
	if viper.GetBool("web.log") {
		r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
			refererRaw := param.Request.Referer()
			refererUrl, _ := url.Parse(refererRaw)
			referer := refererUrl.Host
			if referer == "" {
				referer = "N/A"
			}

			// ban ping
			if param.Request.URL.Path == "/ping" || param.Method == "OPTIONS" {
				return ""
			}

			return fmt.Sprintf("[GIN] %v | %d | %13v | %20s | %40s | %-6s \"%s\"\n",
				param.TimeStamp.Format("2006/01/02 - 15:04:05"),
				param.StatusCode,
				param.Latency,
				referer,
				param.ClientIP,
				param.Method,
				param.Path,
			)
		}))
	}
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())
	r.LoadHTMLFiles(config.DistPath + "/index.html")

	// routers
	initRoute(r)

	// start server
	log.SetOutput(os.Stdout)
	log.Println("[INFO] server running on", viper.GetString("web.address"))
	err := r.Run(viper.GetString("web.address"))
	if err != nil {
		log.Fatalf("[ERROR] web 服务启动失败: %s", err)
	}
}
