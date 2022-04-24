package web

import (
	"busuanzi/core"
	"github.com/gin-gonic/gin"
	"net/url"
)

func ApiHandler(c *gin.Context) { // test redisHelper

	// get referer url
	var referer = c.Request.Referer()
	if referer == "" {
		c.JSON(200, gin.H{
			"success": false,
			"message": "empty referer",
		})
	}

	u, err := url.Parse(referer)
	if err != nil {
		c.JSON(200, gin.H{
			"success": false,
			"message": "unable to parse referer",
		})
	}

	var host = u.Host
	var path = u.Path
	var ip = c.ClientIP()

	// count
	sitePv, siteUv, pagePv, pageUv := core.Count(host, path, ip)

	c.JSON(200, gin.H{
		"success": true,
		"message": "ok",
		"data": gin.H{
			"site_pv": sitePv,
			"site_uv": siteUv,
			"page_pv": pagePv,
			"page_uv": pageUv,
		},
	})
}
