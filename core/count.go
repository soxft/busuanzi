package core

import (
	"context"
	"github.com/soxft/busuanzi/library/tool"
	"github.com/soxft/busuanzi/process/redisutil"
	"github.com/spf13/viper"
	"strings"
	"time"
)

//index		数据类型	        key
//sitePv	Str	            bsz:site_pv:md5(host)
//siteUv	HyperLogLog		bsz:site_uv:md5(host)
//pagePv	zset	        bsz:page_pv:md5(host) / md5(path)
//pageUv	HyperLogLog		bsz:site_uv:md5(host):md5(host&path)

type Counts struct {
	SitePv int64
	SiteUv int64
	PagePv int64
	PageUv int64
}

type RKeys struct {
	SitePvKey  string
	SiteUvKey  string
	PagePvKey  string
	PageUvKey  string
	SiteUnique string
	PathUnique string
}

// Count
// @description return and count the number of users in the redis
func Count(ctx context.Context, host string, path string, userIdentity string) Counts {
	_redis := redisutil.RDB

	rk := getKeys(host, path)

	// sitePV and pagePV 使用 Str / Zset 存储
	sitePv, _ := _redis.Incr(ctx, rk.SitePvKey).Result()
	pagePv, _ := _redis.ZIncrBy(ctx, rk.PagePvKey, 1, rk.PathUnique).Result()

	// siteUv 和 pageUv 使用 HyperLogLog 存储
	_redis.PFAdd(ctx, rk.SiteUvKey, userIdentity)
	_redis.PFAdd(ctx, rk.PageUvKey, userIdentity)

	// count siteUv and pageUv
	siteUv, _ := _redis.PFCount(ctx, rk.SiteUvKey).Result()
	pageUv, _ := _redis.PFCount(ctx, rk.PageUvKey).Result()

	// setExpire
	go setExpire(ctx, rk.SiteUvKey, rk.PageUvKey, rk.SitePvKey, rk.PagePvKey)

	return Counts{
		SitePv: sitePv,
		SiteUv: siteUv,
		PagePv: int64(pagePv),
		PageUv: pageUv,
	}
}

// Put
// @description put data only
func Put(ctx context.Context, host string, path string, userIdentity string) {
	_redis := redisutil.RDB

	rk := getKeys(host, path)

	// sitePV and pagePV 使用 Str / Zset 存储
	_redis.Incr(ctx, rk.SitePvKey)
	_redis.ZIncrBy(ctx, rk.PagePvKey, 1, rk.PathUnique)

	// siteUv 和 pageUv 使用 HyperLogLog 存储
	_redis.PFAdd(ctx, rk.SiteUvKey, userIdentity)
	_redis.PFAdd(ctx, rk.PageUvKey, userIdentity)

	// setExpire
	go setExpire(ctx, rk.SiteUvKey, rk.PageUvKey, rk.SitePvKey, rk.PagePvKey)
	return
}

// Get bsz counts
func Get(ctx context.Context, host string, path string, userIdentity string) Counts {
	_redis := redisutil.RDB

	rk := getKeys(host, path)

	// sitePV and pagePV 使用 Str / Zset 存储
	sitePv, _ := _redis.Get(ctx, rk.SitePvKey).Int64()
	pagePv, _ := _redis.ZScore(ctx, rk.PagePvKey, rk.PathUnique).Result()

	// count siteUv and pageUv
	siteUv, _ := _redis.PFCount(ctx, rk.SiteUvKey).Result()
	pageUv, _ := _redis.PFCount(ctx, rk.PageUvKey).Result()

	// setExpire
	go setExpire(ctx, rk.SiteUvKey, rk.PageUvKey, rk.SitePvKey, rk.PagePvKey)

	return Counts{
		SitePv: sitePv,
		SiteUv: siteUv,
		PagePv: int64(pagePv),
		PageUv: pageUv,
	}
}

// setExpire
// @description set the expiration time of the key
// TODO 使用 MQ 进行防抖消重, 避免频繁调用
func setExpire(ctx context.Context, key ...string) {
	if viper.GetInt("bsz.expire") == 0 {
		return
	}

	_redis := redisutil.RDB

	for _, k := range key {
		_redis.Expire(ctx, k, viper.GetDuration("bsz.expire")*time.Second)
	}
}

func getKeys(host string, path string) RKeys {
	var siteUnique = host
	var pathUnique = host + path

	// 兼容旧版本
	if viper.GetBool("bsz.pathStyle") == false {
		pathUnique = host + "&" + path
	}

	// encrypt
	switch viper.GetString("bsz.Encrypt") {
	case "MD516":
		siteUnique = tool.Md5(siteUnique)[8:24]
		pathUnique = tool.Md5(pathUnique)[8:24]
	case "MD532":
		siteUnique = tool.Md5(siteUnique)
		pathUnique = tool.Md5(pathUnique)
	default:
		siteUnique = tool.Md5(siteUnique)
		pathUnique = tool.Md5(pathUnique)
	}

	redisPrefix := viper.GetString("redis.prefix")

	siteUvKey := strings.Join([]string{redisPrefix, "site_uv", siteUnique}, ":")
	pageUvKey := strings.Join([]string{redisPrefix, "page_uv", siteUnique, pathUnique}, ":")

	sitePvKey := strings.Join([]string{redisPrefix, "site_pv", siteUnique}, ":")
	pagePvKey := strings.Join([]string{redisPrefix, "page_pv", siteUnique}, ":")

	return RKeys{
		SitePvKey:  sitePvKey,
		SiteUvKey:  siteUvKey,
		PagePvKey:  pagePvKey,
		PageUvKey:  pageUvKey,
		SiteUnique: siteUnique,
		PathUnique: pathUnique,
	}
}
