package controller

import (
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/soxft/busuanzi/core"
)

func Admin(c *gin.Context) {
	c.HTML(200, "admin.html", gin.H{})
}

func AdminGet(c *gin.Context) {
	pageUrl := c.PostForm("url")
	u, err := url.Parse(pageUrl)

	if err != nil {
		c.JSON(200, gin.H{
			"success": false,
			"message": "unable to parse url",
			"data":    gin.H{},
		})
		return
	} else if u.Host == "" {
		c.JSON(200, gin.H{
			"success": false,
			"message": "invalid url",
			"data":    gin.H{},
		})
		return
	}

	var host = u.Host
	var path = u.Path

	sitePv, siteUv, pagePv, pageUv := core.Get(host, path)

	c.JSON(200, gin.H{
		"success": true,
		"message": "ok",
		"data": gin.H{
			"site_pv":  sitePv,
			"site_uv":  siteUv,
			"page_pv":  pagePv,
			"page_uv":  pageUv,
			"page_url": pageUrl,
			"host":     host,
			"path":     path,
		},
	})
}

func AdminUpdate(c *gin.Context) {
	pageUrl := c.PostForm("url")
	u, err := url.Parse(pageUrl)
	if err != nil {
		c.JSON(200, gin.H{
			"success": false,
			"message": "unable to parse url",
			"data":    gin.H{},
		})
		return
	} else if u.Host == "" {
		c.JSON(200, gin.H{
			"success": false,
			"message": "invalid url",
			"data":    gin.H{},
		})
		return
	}

	var host = u.Host
	var path = u.Path

	sitePv, _ := strconv.Atoi(c.PostForm("site_pv"))
	siteUv, _ := strconv.Atoi(c.PostForm("site_uv"))
	pagePv, _ := strconv.Atoi(c.PostForm("page_pv"))
	pageUv, _ := strconv.Atoi(c.PostForm("page_uv"))

	count := core.Update(host, path, sitePv, siteUv, pagePv, pageUv)

	c.JSON(200, gin.H{
		"success": true,
		"message": "ok",
		"data": gin.H{
			"count": count,
		},
	})
}
