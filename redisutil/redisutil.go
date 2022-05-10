package redisutil

import (
	"busuanzi/config"
	"github.com/gomodule/redigo/redis"
)

var (
	Pool *redis.Pool
)

func init() {
	Pool = &redis.Pool{
		MaxIdle:   config.C.Redis.MaxIdle,
		MaxActive: config.C.Redis.MaxActive,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.C.Redis.Address,
				redis.DialPassword(config.C.Redis.Password),
				redis.DialDatabase(config.C.Redis.Database),
			)
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}
