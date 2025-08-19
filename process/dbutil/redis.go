package dbutil

import (
	"context"
	"github.com/soxft/busuanzi/process/redisutil"
	"github.com/spf13/viper"
	"strings"
)

type RedisDB struct{}

func (r *RedisDB) IncrSitePv(ctx context.Context, siteKey string) (int64, error) {
	redisPrefix := viper.GetString("database.prefix")
	key := strings.Join([]string{redisPrefix, "site_pv", siteKey}, ":")
	return redisutil.RDB.Incr(ctx, key).Result()
}

func (r *RedisDB) IncrPagePv(ctx context.Context, siteKey, pathKey string) (float64, error) {
	redisPrefix := viper.GetString("database.prefix")
	key := strings.Join([]string{redisPrefix, "page_pv", siteKey}, ":")
	return redisutil.RDB.ZIncrBy(ctx, key, 1, pathKey).Result()
}

func (r *RedisDB) AddSiteUv(ctx context.Context, siteKey, userIdentity string) error {
	redisPrefix := viper.GetString("database.prefix")
	key := strings.Join([]string{redisPrefix, "site_uv", siteKey}, ":")
	return redisutil.RDB.PFAdd(ctx, key, userIdentity).Err()
}

func (r *RedisDB) AddPageUv(ctx context.Context, siteKey, pathKey, userIdentity string) error {
	redisPrefix := viper.GetString("database.prefix")
	key := strings.Join([]string{redisPrefix, "page_uv", siteKey, pathKey}, ":")
	return redisutil.RDB.PFAdd(ctx, key, userIdentity).Err()
}

func (r *RedisDB) GetSitePv(ctx context.Context, siteKey string) (int64, error) {
	redisPrefix := viper.GetString("database.prefix")
	key := strings.Join([]string{redisPrefix, "site_pv", siteKey}, ":")
	return redisutil.RDB.Get(ctx, key).Int64()
}

func (r *RedisDB) GetPagePv(ctx context.Context, siteKey, pathKey string) (float64, error) {
	redisPrefix := viper.GetString("database.prefix")
	key := strings.Join([]string{redisPrefix, "page_pv", siteKey}, ":")
	return redisutil.RDB.ZScore(ctx, key, pathKey).Result()
}

func (r *RedisDB) CountSiteUv(ctx context.Context, siteKey string) (int64, error) {
	redisPrefix := viper.GetString("database.prefix")
	key := strings.Join([]string{redisPrefix, "site_uv", siteKey}, ":")
	return redisutil.RDB.PFCount(ctx, key).Result()
}

func (r *RedisDB) CountPageUv(ctx context.Context, siteKey, pathKey string) (int64, error) {
	redisPrefix := viper.GetString("database.prefix")
	key := strings.Join([]string{redisPrefix, "page_uv", siteKey, pathKey}, ":")
	return redisutil.RDB.PFCount(ctx, key).Result()
}

func (r *RedisDB) SetExpire(keys ...string) {
	// This will be handled by the existing setExpire function in core
}