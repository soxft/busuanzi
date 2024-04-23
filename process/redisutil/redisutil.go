package redisutil

import (
	"context"
	"crypto/tls"
	"github.com/redis/go-redis/v9"
	"github.com/soxft/busuanzi/config"

	"log"
	"time"
)

var RDB *redis.Client

func Init() {
	log.Printf("[INFO] Redis trying connect to tcp://%s/%d", config.Redis.Address, config.Redis.Database)

	r := config.Redis

	var tlsConfig *tls.Config

	if r.TLS {
		tlsConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:            r.Address,
		Password:        r.Password,
		DB:              r.Database,
		TLSConfig:       tlsConfig,
		MinIdleConns:    r.MinIdle,
		MaxIdleConns:    r.MaxIdle,
		MaxRetries:      r.MaxRetries,
		ConnMaxLifetime: 5 * time.Minute,
		MaxActiveConns:  r.MaxActive,
	})

	RDB = rdb

	// test redis
	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("[ERROR] Redis ping failed: %v", err)
	}

	log.Printf("[INFO] Redis init success, pong: %s ", pong)
}
