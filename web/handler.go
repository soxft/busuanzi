package web

import (
	"busuanzi/redisHelper"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

var pool = redisHelper.NewPool()

func ApiHandler(c *gin.Context) { // test redisHelper
	var _redis = pool.Get()
	defer func(_redis redis.Conn) {
		_ = _redis.Close()
	}(_redis)

	c.JSON(200, map[string]interface{}{
		"message": "Hello,  World!",
	})
}
