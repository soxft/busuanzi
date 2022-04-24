package core

import (
	"busuanzi/config"
	"busuanzi/redisHelper"
	"busuanzi/tool"
	"github.com/gomodule/redigo/redis"
)

// Count return and count the number of users in the redis
// @return int,int,int,int site_pv,site_uv,page_pv,page_uv
func Count(host string, path string, ip string) (int, int, int, int) {
	var _redis = redisHelper.Pool.Get()
	defer func(_redis redis.Conn) {
		_ = _redis.Close()
	}(_redis)

	var keyExpire = config.C.Bsz.Expire

	// encode
	var pathUnique = tool.Md5(host + path)
	var siteUnique = tool.Md5(host)

	redisPrefix := config.C.Redis.Prefix
	sitePvKey := redisPrefix + ":site_pv:" + siteUnique
	siteUvKey := redisPrefix + ":site_uv:" + siteUnique
	pagePvKey := redisPrefix + ":page_pv:" + pathUnique
	pageUvKey := redisPrefix + ":page_uv:" + pathUnique

	// count sitePv ans pagePv use pipeline
	_, _ = _redis.Do("WATCH", sitePvKey, pagePvKey)
	_, _ = _redis.Do("MULTI")
	_, _ = _redis.Do("INCR", sitePvKey)
	_, _ = _redis.Do("INCR", pagePvKey)
	_, _ = _redis.Do("SADD", siteUvKey, tool.Md5(ip))
	_, _ = _redis.Do("SADD", pageUvKey, tool.Md5(ip))
	_, _ = redis.Int(_redis.Do("SCARD", siteUvKey))
	_, _ = redis.Int(_redis.Do("SCARD", pageUvKey))
	_, _ = _redis.Do("EXPIRE", sitePvKey, keyExpire)
	_, _ = _redis.Do("EXPIRE", pagePvKey, keyExpire)
	_, _ = _redis.Do("EXPIRE", siteUvKey, keyExpire)
	_, _ = _redis.Do("EXPIRE", pageUvKey, keyExpire)
	res, err := redis.Values(_redis.Do("EXEC"))
	if err != nil {
		return 0, 0, 0, 0
	}

	sitePv := int(res[0].(int64))
	pagePv := int(res[1].(int64))
	siteUv := int(res[4].(int64))
	pageUv := int(res[5].(int64))

	return sitePv, siteUv, pagePv, pageUv
}
