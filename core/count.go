package core

import (
	"context"
	"fmt"
	"github.com/soxft/busuanzi/config"
	"github.com/soxft/busuanzi/library/tool"
	"github.com/soxft/busuanzi/process/redisutil"
	"strings"
)

//index		数据类型	        key
//sitePv	string	        bsz:site_pv:md5(example.com)
//siteUv	HyperLogLog		bsz:site_uv:md5(example.com)
//pagePv	zset	        bsz:page_pv:md5(example.com)
//pageUv	HyperLogLog		bsz:site_uv:md5(example.com):md5(example.com&index.html)

// Count
// @description return and count the number of users in the redis
func Count(ctx context.Context, host string, path string, userIdentity string) (int64, int64, int64, int64) {
	_redis := redisutil.RDB

	// encode
	var pathUnique = strings.ToLower(tool.Md5(host + "&" + path))
	var siteUnique = strings.ToLower(tool.Md5(host))

	redisPrefix := config.Redis.Prefix

	// user view keys 用户数
	siteUvKey := fmt.Sprintf("%s:site_uv:%s", redisPrefix, siteUnique)
	pageUvKey := fmt.Sprintf("%s:page_uv:%s:%s", redisPrefix, siteUnique, pathUnique)

	// page view keys 页面访问数
	sitePvKey := fmt.Sprintf("%s:site_pv:%s", redisPrefix, siteUnique)
	pagePvKey := fmt.Sprintf("%s:page_pv:%s", redisPrefix, siteUnique)

	// count sitePv ans pagePv
	sitePv, _ := _redis.Incr(ctx, sitePvKey).Result()
	pagePv, _ := _redis.ZIncrBy(ctx, pagePvKey, 1, pathUnique).Result() // pagePv 使用 ZSet 存储

	// siteUv 和 pageUv 使用 HyperLogLog 存储
	_redis.PFAdd(ctx, siteUvKey, userIdentity)
	_redis.PFAdd(ctx, pageUvKey, userIdentity)

	// count siteUv and pageUv
	siteUv, _ := _redis.PFCount(ctx, siteUvKey).Result()
	pageUv, _ := _redis.PFCount(ctx, pageUvKey).Result()

	return sitePv, siteUv, int64(pagePv), pageUv
}
