package web

import (
	"busuanzi/redisHelper"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

func ApiHandler(c *gin.Context) { // test redisHelper
	var _redis = redisHelper.Pool.Get()
	defer func(_redis redis.Conn) {
		_ = _redis.Close()
	}(_redis)

	var a, _ = _redis.Do("GET", "a")

	// get referer url
	var referer = c.Request.Referer()
	var ip = c.ClientIP()

	c.JSON(200, map[string]interface{}{
		"message": "Hello,  World!",
		"a":       a,
		"referer": referer,
		"ip":      ip,
	})
}
