package core

import (
	"busuanzi/config"
	"busuanzi/redisutil"
	"busuanzi/tool"
	"github.com/gomodule/redigo/redis"
)

// Count return and count the number of users in the redis
// @return int,int,int,int site_pv,site_uv,page_pv,page_uv
func Count(host string, path string, ip string) (int, int, int, int) {
	var _redis = redisutil.Pool.Get()
	defer func(_redis redis.Conn) {
		_ = _redis.Close()
	}(_redis)

	// encode
	var pathUnique = tool.Md5(host + path)
	var siteUnique = tool.Md5(host)

	redisPrefix := config.C.Redis.Prefix
	sitePvKey := redisPrefix + ":site_pv:" + siteUnique
	siteUvKey := redisPrefix + ":site_uv:" + siteUnique
	pagePvKey := redisPrefix + ":page_pv:" + pathUnique
	pageUvKey := redisPrefix + ":page_uv:" + pathUnique

	// count sitePv ans pagePv
	sitePv, _ := redis.Int(_redis.Do("INCR", sitePvKey))
	pagePv, _ := redis.Int(_redis.Do("INCR", pagePvKey))
	_, _ = _redis.Do("SADD", siteUvKey, tool.Md5(ip))
	_, _ = _redis.Do("SADD", pageUvKey, tool.Md5(ip))

	siteUv, _ := redis.Int(_redis.Do("SCARD", siteUvKey))
	pageUv, _ := redis.Int(_redis.Do("SCARD", pageUvKey))
	setExpire(_redis, sitePvKey, siteUvKey, pagePvKey, pageUvKey)

	return sitePv, siteUv, pagePv, pageUv
}

func setExpire(redis redis.Conn, key ...string) {
	// multi-set expire
	_, _ = redis.Do("MULTI")
	for _, k := range key {
		_, _ = redis.Do("EXPIRE", k, config.C.Bsz.Expire)
	}
	_, _ = redis.Do("EXEC")
}
