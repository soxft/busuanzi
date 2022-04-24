package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strings"
)

func NoRouteHandler(c *gin.Context) {
	path := "./dist" + c.Request.URL.Path
	if strings.HasSuffix(path, "/") {
		path += "index.html"
	}
	resource, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(path)
		c.JSON(404, gin.H{
			"success": false,
			"message": "Page not found",
		})
		return
	}
	// web
	ext := path[strings.LastIndex(path, "."):]
	if ext == ".html" {
		c.Header("Content-Type", "text/html")
	} else if ext == ".js" {
		c.Header("Content-Type", "text/javascript")
	} else if ext == ".css" {
		c.Header("Content-Type", "text/css")
	}
	c.String(200, string(resource))
	return

}
