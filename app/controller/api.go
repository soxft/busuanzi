package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/soxft/busuanzi/core"
	"net/url"
)

func ApiHandler(c *gin.Context) {
	var referer = c.Request.Header.Get("x-bsz-referer")
	if referer == "" {
		c.JSON(200, gin.H{
			"success": false,
			"message": "invalid referer",
			"data": gin.H{
				"project": "https://github.com/soxft/busuanzi",
				"usage":   "https://github.com/soxft/busuanzi/wiki/usage",
			},
		})
		return
	}

	u, err := url.Parse(referer)
	if err != nil {
		c.JSON(200, gin.H{
			"success": false,
			"message": "unable to parse referer",
			"data": gin.H{
				"project": "https://github.com/soxft/busuanzi",
				"usage":   "https://github.com/soxft/busuanzi/wiki/usage",
			},
		})
		return
	} else if u.Host == "" {
		c.JSON(200, gin.H{
			"success": false,
			"message": "invalid referer",
			"data": gin.H{
				"project": "https://github.com/soxft/busuanzi",
				"usage":   "https://github.com/soxft/busuanzi/wiki/usage",
			},
		})
		return
	}

	var host = u.Host
	var path = u.Path

	userIdentity := c.GetString("user_identity")

	// count
	counts := core.Count(c, host, path, userIdentity)

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
