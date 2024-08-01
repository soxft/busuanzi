package core

import (
	"context"
	"fmt"
	"github.com/soxft/busuanzi/library/tool"
	"github.com/soxft/busuanzi/process/redisutil"
	"github.com/spf13/viper"
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

// Count
// @description return and count the number of users in the redis
func Count(ctx context.Context, host string, path string, userIdentity string) Counts {
	_redis := redisutil.RDB

	// encode
	var siteUnique = host
	var pathUnique = path

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

	// redisPrefix
	redisPrefix := viper.GetString("redis.prefix")

	// user view keys 用户数
	siteUvKey := fmt.Sprintf("%s:site_uv:%s", redisPrefix, siteUnique)
	pageUvKey := fmt.Sprintf("%s:page_uv:%s:%s", redisPrefix, siteUnique, pathUnique)

	// page view keys 页面访问数
	sitePvKey := fmt.Sprintf("%s:site_pv:%s", redisPrefix, siteUnique)
	pagePvKey := fmt.Sprintf("%s:page_pv:%s", redisPrefix, siteUnique)

	// sitePV and pagePV 使用 Str / Zset 存储
	sitePv, _ := _redis.Incr(ctx, sitePvKey).Result()
	pagePv, _ := _redis.ZIncrBy(ctx, pagePvKey, 1, pathUnique).Result()

	// siteUv 和 pageUv 使用 HyperLogLog 存储
	_redis.PFAdd(ctx, siteUvKey, userIdentity)
	_redis.PFAdd(ctx, pageUvKey, userIdentity)

	// count siteUv and pageUv
	siteUv, _ := _redis.PFCount(ctx, siteUvKey).Result()
	pageUv, _ := _redis.PFCount(ctx, pageUvKey).Result()

	// setExpire
	go setExpire(ctx, siteUvKey, pageUvKey, sitePvKey, pagePvKey)

	return Counts{
		SitePv: sitePv,
		SiteUv: siteUv,
		PagePv: int64(pagePv),
		PageUv: pageUv,
	}
}

// setExpire
// @description set the expiration time of the key
func setExpire(ctx context.Context, key ...string) {
	if viper.GetInt("bsz.expire") == 0 {
		return
	}

	_redis := redisutil.RDB

	for _, k := range key {
		_redis.Expire(ctx, k, viper.GetDuration("bsz.expire")*time.Second)
	}
}
