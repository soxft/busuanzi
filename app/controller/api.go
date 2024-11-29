package controller

import (
	"encoding/json"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/soxft/busuanzi/core"
	"github.com/soxft/busuanzi/library/tool"
)

var defaultData = gin.H{
	"project": "https://github.com/soxft/busuanzi",
	"usage":   "https://github.com/soxft/busuanzi/wiki/usage",
}

func ApiHandler(c *gin.Context) {
	var referer = c.Request.Header.Get("x-bsz-referer")
	if referer == "" {
		c.JSON(200, gin.H{
			"success": false,
			"message": "invalid referer",
			"data":    defaultData,
		})
		return
	}

	u, err := url.Parse(referer)
	if err != nil {
		c.JSON(200, gin.H{
			"success": false,
			"message": "unable to parse referer",
			"data":    defaultData,
		})
		return
	} else if u.Host == "" {
		c.JSON(200, gin.H{
			"success": false,
			"message": "invalid referer",
			"data":    defaultData,
		})
		return
	}

	var host = u.Host
	var path = u.Path

	// count
	counts := core.Count(c, host, path, c.GetString("user_identity"))

	// json
	c.JSON(200, gin.H{
		"success": true,
		"message": "ok",
		"data": gin.H{
			"site_pv": counts.SitePv,
			"site_uv": counts.SiteUv,
			"page_pv": counts.PagePv,
			"page_uv": counts.PageUv,
		},
	})
}

// PutHandler 仅提交数据, 不做返回
func PutHandler(c *gin.Context) {
	var referer = c.Request.Header.Get("x-bsz-referer")
	if referer == "" {
		c.Status(400)
		return
	}

	u, err := url.Parse(referer)
	if err != nil {
		c.Status(400)
		return
	} else if u.Host == "" {
		c.Status(400)
		return
	}

	var host = u.Host
	var path = u.Path

	// count
	go core.Put(c, host, path, c.GetString("user_identity"))

	// json
	c.Status(204)
}

// GetHandler 仅获取数据, 不增加
func GetHandler(c *gin.Context) {
	var referer = c.Request.Header.Get("x-bsz-referer")
	if referer == "" {
		c.JSON(200, gin.H{
			"success": false,
			"message": "invalid referer",
			"data":    defaultData,
		})
		return
	}

	u, err := url.Parse(referer)
	if err != nil {
		c.JSON(200, gin.H{
			"success": false,
			"message": "unable to parse referer",
			"data":    defaultData,
		})
		return
	} else if u.Host == "" {
		c.JSON(200, gin.H{
			"success": false,
			"message": "invalid referer",
			"data":    defaultData,
		})
		return
	}

	var host = u.Host
	var path = u.Path

	// count
	counts := core.Get(c, host, path)

	// json
	c.JSON(200, gin.H{
		"success": true,
		"message": "ok",
		"data": gin.H{
			"site_pv": counts.SitePv,
			"site_uv": counts.SiteUv,
			"page_pv": counts.PagePv,
			"page_uv": counts.PageUv,
		},
	})
}

func JsonpHandler(c *gin.Context) {
	c.Header("Content-Type", "application/javascript")

	callback := c.Query("callback")
	if callback == "" {
		c.String(200, "try{console.error(%s)}catch{}", "missing callback parameter")
		return
	}

	u, err := url.Parse(c.GetHeader("Referer"))
	if err != nil {
		c.String(200, "try{console.error(%s)}catch{}", "unable to parse referer")
		return
	} else if u.Host == "" {
		c.String(200, "try{console.error(%s)}catch{}", "invalid referer")
		return
	}

	var host = u.Host
	var path = u.Path
	counts := core.Count(c, host, path, tool.Md5(c.ClientIP()+c.Request.UserAgent()))

	data := gin.H{
		"site_pv": counts.SitePv,
		"site_uv": counts.SiteUv,
		"page_pv": counts.PagePv,
		"page_uv": counts.PageUv,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		c.String(200, "try{console.error(%s)}catch{}", "gen json failed")
		return
	}

	c.String(200, "try{%s(%s);}catch(e){}", callback, string(jsonData))
}
