package core

import (
	"context"
	"github.com/soxft/busuanzi/library/tool"
	"github.com/soxft/busuanzi/process/redisutil"
	"github.com/spf13/viper"
	"strings"
)

// Count
// @description return and count the number of users in the redis
func Count(ctx context.Context, host string, path string, userIdentity string) Counts {
	return CountWithScope(ctx, host, path, userIdentity, ScopeAll)
}

// CountWithScope
// @description return and count the number of users in the redis with specified scope
func CountWithScope(ctx context.Context, host string, path string, userIdentity string, scope Scope) Counts {
	_redis := redisutil.RDB

	rk := getKeys(host, path)

	var sitePv int64
	var siteUv int64
	var pagePvFloat float64
	var pageUv int64

	// 根据 scope 更新不同的统计
	switch scope {
	case ScopePage:
		// 仅更新页面统计
		pagePvFloat, _ = _redis.ZIncrBy(ctx, rk.PagePvKey, 1, rk.PathUnique).Result()
		_redis.PFAdd(ctx, rk.PageUvKey, userIdentity)
		pageUv, _ = _redis.PFCount(ctx, rk.PageUvKey).Result()
		// 获取站点统计（不更新）
		sitePv, _ = _redis.Get(ctx, rk.SitePvKey).Int64()
		siteUv, _ = _redis.PFCount(ctx, rk.SiteUvKey).Result()
		go setExpire(rk.PageUvKey, rk.PagePvKey)
	case ScopeSite:
		// 仅更新站点统计
		sitePv, _ = _redis.Incr(ctx, rk.SitePvKey).Result()
		_redis.PFAdd(ctx, rk.SiteUvKey, userIdentity)
		siteUv, _ = _redis.PFCount(ctx, rk.SiteUvKey).Result()
		// 获取页面统计（不更新）
		pagePvFloat, _ = _redis.ZScore(ctx, rk.PagePvKey, rk.PathUnique).Result()
		pageUv, _ = _redis.PFCount(ctx, rk.PageUvKey).Result()
		go setExpire(rk.SiteUvKey, rk.SitePvKey)
	default: // ScopeAll
		// 更新所有统计（原有行为）
		sitePv, _ = _redis.Incr(ctx, rk.SitePvKey).Result()
		pagePvFloat, _ = _redis.ZIncrBy(ctx, rk.PagePvKey, 1, rk.PathUnique).Result()
		_redis.PFAdd(ctx, rk.SiteUvKey, userIdentity)
		_redis.PFAdd(ctx, rk.PageUvKey, userIdentity)
		siteUv, _ = _redis.PFCount(ctx, rk.SiteUvKey).Result()
		pageUv, _ = _redis.PFCount(ctx, rk.PageUvKey).Result()
		go setExpire(rk.SiteUvKey, rk.PageUvKey, rk.SitePvKey, rk.PagePvKey)
	}

	return Counts{
		SitePv: sitePv,
		SiteUv: siteUv,
		PagePv: int64(pagePvFloat),
		PageUv: pageUv,
	}
}

// Put
// @description put data only
func Put(ctx context.Context, host string, path string, userIdentity string) {
	PutWithScope(ctx, host, path, userIdentity, ScopeAll)
}

// PutWithScope
// @description put data only with specified scope
func PutWithScope(ctx context.Context, host string, path string, userIdentity string, scope Scope) {
	_redis := redisutil.RDB

	rk := getKeys(host, path)

	switch scope {
	case ScopePage:
		// 仅更新页面统计
		_redis.ZIncrBy(ctx, rk.PagePvKey, 1, rk.PathUnique)
		_redis.PFAdd(ctx, rk.PageUvKey, userIdentity)
		go setExpire(rk.PageUvKey, rk.PagePvKey)
	case ScopeSite:
		// 仅更新站点统计
		_redis.Incr(ctx, rk.SitePvKey)
		_redis.PFAdd(ctx, rk.SiteUvKey, userIdentity)
		go setExpire(rk.SiteUvKey, rk.SitePvKey)
	default: // ScopeAll
		// 更新所有统计
		_redis.Incr(ctx, rk.SitePvKey)
		_redis.ZIncrBy(ctx, rk.PagePvKey, 1, rk.PathUnique)
		_redis.PFAdd(ctx, rk.SiteUvKey, userIdentity)
		_redis.PFAdd(ctx, rk.PageUvKey, userIdentity)
		go setExpire(rk.SiteUvKey, rk.PageUvKey, rk.SitePvKey, rk.PagePvKey)
	}
	return
}

// Get bsz counts
func Get(ctx context.Context, host string, path string) Counts {
	_redis := redisutil.RDB

	rk := getKeys(host, path)

	// sitePV and pagePV 使用 Str / Zset 存储
	sitePv, _ := _redis.Get(ctx, rk.SitePvKey).Int64()
	pagePv, _ := _redis.ZScore(ctx, rk.PagePvKey, rk.PathUnique).Result()

	// count siteUv and pageUv
	siteUv, _ := _redis.PFCount(ctx, rk.SiteUvKey).Result()
	pageUv, _ := _redis.PFCount(ctx, rk.PageUvKey).Result()

	// setExpire
	go setExpire(rk.SiteUvKey, rk.PageUvKey, rk.SitePvKey, rk.PagePvKey)

	return Counts{
		SitePv: sitePv,
		SiteUv: siteUv,
		PagePv: int64(pagePv),
		PageUv: pageUv,
	}
}

func getKeys(host string, path string) RKeys {
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
