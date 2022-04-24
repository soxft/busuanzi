package redisHelper

import (
	"busuanzi/config"
	"github.com/gomodule/redigo/redis"
)

func NewPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.RedisAddr,
				redis.DialPassword(config.RedisPassword),
				redis.DialDatabase(config.RedisDB),
			)
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}
