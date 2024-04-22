package core

import (
	"context"
	"fmt"
	"github.com/soxft/busuanzi/config"
	"github.com/soxft/busuanzi/library/tool"
	"github.com/soxft/busuanzi/process/redisutil"
)

// Count
// @description return and count the number of users in the redis
func Count(ctx context.Context, host string, path string, userIdentity string) (int64, int64, int64, int64) {
	_redis := redisutil.RDB

	// encode
	var pathUnique = tool.Md5(host + "&" + path)
	var siteUnique = tool.Md5(host)

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

	// siteUv 和 pageUv 使用 Set 存储
	_, _ = _redis.SAdd(ctx, siteUvKey, userIdentity).Result()
	_, _ = _redis.SAdd(ctx, pageUvKey, userIdentity).Result()

	siteUv, _ := _redis.SCard(ctx, siteUvKey).Result()
	pageUv, _ := _redis.SCard(ctx, pageUvKey).Result()

	return sitePv, siteUv, int64(pagePv), pageUv
}
