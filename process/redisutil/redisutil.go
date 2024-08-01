package redisutil

import (
	"context"
	"crypto/tls"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"time"

	"log"
)

var RDB *redis.Client

func Init() {
	log.Printf("[INFO] Redis trying connect to tcp://%s/%d", viper.GetString("redis.address"), viper.GetInt("redis.database"))

	option := &redis.Options{
		Addr:            viper.GetString("redis.address"),
		Password:        viper.GetString("redis.password"),
		DB:              viper.GetInt("redis.database"),
		MinIdleConns:    viper.GetInt("redis.MinIdle"),
		MaxIdleConns:    viper.GetInt("redis.MaxIdle"),
		MaxRetries:      viper.GetInt("redis.MaxRetries"),
		MaxActiveConns:  viper.GetInt("redis.MaxActive"),
		ConnMaxLifetime: 5 * time.Minute,
	}
	if viper.GetBool("redis.tls") {
		option.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	rdb := redis.NewClient(option)

	RDB = rdb

	// test redis
	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("[ERROR] Redis ping failed: %v", err)
	}

	log.Printf("[INFO] Redis init success, pong: %s ", pong)
}
