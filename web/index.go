package web

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

func Index(c *gin.Context) {
	resource, err := ioutil.ReadFile("./build/index.html")
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.Header("Content-Type", "text/html")
	c.String(200, string(resource))
}
