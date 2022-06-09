package redisutil

import (
	"busuanzi/config"
	"github.com/gomodule/redigo/redis"
	"log"
)

var (
	Pool *redis.Pool
)

func init() {
	Pool = &redis.Pool{
		MaxIdle:   config.Redis.MaxIdle,
		MaxActive: config.Redis.MaxActive,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.Redis.Address,
				redis.DialPassword(config.Redis.Password),
				redis.DialDatabase(config.Redis.Database),
			)
			if err != nil {
				log.Fatalf("redis connect error: %s", err.Error())
			}
			return c, err
		},
	}
	_, _ = Pool.Get().Do("PING")
}
