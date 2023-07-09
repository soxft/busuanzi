package core

import (
	"github.com/gomodule/redigo/redis"
	"github.com/soxft/busuanzi/config"
	"github.com/soxft/busuanzi/library/tool"
	"github.com/soxft/busuanzi/process/redisutil"
)

// Count
// @description return and count the number of users in the redis
func Count(host string, path string, userIdentity string) (int, int, int, int) {
	_redis := redisutil.Pool.Get()
	defer func(_redis redis.Conn) {
		_ = _redis.Close()
	}(_redis)

	sitePvKey, siteUvKey, pagePvKey, pageUvKey, pathUnique := getRedisKey(host, path)

	// count sitePv ans pagePv
	sitePv, _ := redis.Int(_redis.Do("INCR", sitePvKey))                   // 站点总访问量 +1
	pagePv, _ := redis.Int(_redis.Do("ZINCRBY", pagePvKey, 1, pathUnique)) // 页面总访问量 对应路径 +1
	_, _ = _redis.Do("SADD", siteUvKey, userIdentity)                      // 站点总访问用户 对应用户 +1
	_, _ = _redis.Do("SADD", pageUvKey, userIdentity)                      // 页面总访问用户 对应用户 +1

	siteUv, _ := redis.Int(_redis.Do("SCARD", siteUvKey)) // 获取站点总访问用户
	pageUv, _ := redis.Int(_redis.Do("SCARD", pageUvKey)) // 获取页面总访问用户

	if config.Bsz.Expire > 0 {
		go setExpire(sitePvKey, siteUvKey, pagePvKey, pageUvKey)
	}

	return sitePv, siteUv, pagePv, pageUv
}

func setExpire(key ...string) {
	var _redis = redisutil.Pool.Get()
	defer func(_redis redis.Conn) {
		_ = _redis.Close()
	}(_redis)
	// multi-set expire
	_, _ = _redis.Do("MULTI")
	for _, k := range key {
		_, _ = _redis.Do("EXPIRE", k, config.Bsz.Expire)
	}
	_, _ = _redis.Do("EXEC")
}

func getRedisKey(host string, path string) (string, string, string, string, string) {
	// encode
	var pathUnique = tool.Md5(host + "&" + path)
	var siteUnique = tool.Md5(host)

	redisPrefix := config.Redis.Prefix
	siteUvKey := redisPrefix + ":site_uv:" + siteUnique
	pageUvKey := redisPrefix + ":page_uv:" + siteUnique + ":" + pathUnique

	sitePvKey := redisPrefix + ":site_pv:" + siteUnique
	pagePvKey := redisPrefix + ":page_pv:" + siteUnique

	return sitePvKey, siteUvKey, pagePvKey, pageUvKey, pathUnique
}

// Get
// @description return the number of users in the redis
func Get(host string, path string) (int, int, int, int) {
	_redis := redisutil.Pool.Get()
	defer func(_redis redis.Conn) {
		_ = _redis.Close()
	}(_redis)

	sitePvKey, siteUvKey, pagePvKey, pageUvKey, pathUnique := getRedisKey(host, path)

	// count sitePv ans pagePv
	sitePv, _ := redis.Int(_redis.Do("GET", sitePvKey))
	pagePv, _ := redis.Int(_redis.Do("ZSCORE", pagePvKey, pathUnique))
	siteUv, _ := redis.Int(_redis.Do("SCARD", siteUvKey))
	pageUv, _ := redis.Int(_redis.Do("SCARD", pageUvKey))

	if config.Bsz.Expire > 0 {
		go setExpire(sitePvKey, siteUvKey, pagePvKey, pageUvKey)
	}

	return sitePv, siteUv, pagePv, pageUv
}

// Update
// @description update the number of users in the redis
func Update(host string, path string, sitePv int, siteUv int, pagePv int, pageUv int) int {
	_redis := redisutil.Pool.Get()
	defer func(_redis redis.Conn) {
		_ = _redis.Close()
	}(_redis)

	sitePvKey, siteUvKey, pagePvKey, pageUvKey, pathUnique := getRedisKey(host, path)

	sitePvOld, siteUvOld, pagePvOld, pageUvOld := Get(host, path)

	var diff int
	var count int

	if sitePvOld > sitePv {
		// 删除站点总访问量
		diff = sitePvOld - sitePv
		_, _ = _redis.Do("DECRBY", sitePvKey, diff)
	} else if sitePvOld < sitePv {
		// 增加站点总访问量
		diff = sitePv - sitePvOld
		_, _ = _redis.Do("INCRBY", sitePvKey, diff)
	}
	count = count + diff

	if pagePvOld > pagePv {
		// 删除页面总访问量
		diff = pagePvOld - pagePv
		_, _ = _redis.Do("ZINCRBY", pagePvKey, -diff, pathUnique)
	} else if pagePvOld < pagePv {
		// 增加页面总访问量
		diff = pagePv - pagePvOld
		_, _ = _redis.Do("ZINCRBY", pagePvKey, diff, pathUnique)
	}
	count = count + diff

	if siteUvOld > siteUv {
		// 删除站点总访问用户
		diff = siteUvOld - siteUv
		_, _ = _redis.Do("SPOP", siteUvKey, diff)
	} else if siteUvOld < siteUv {
		// 增加站点总访问用户
		diff = siteUv - siteUvOld
		for i := 0; i < diff; i++ {
			_, _ = _redis.Do("SADD", siteUvKey, tool.Md5(tool.RandString(32)))
		}
	}
	count = count + diff

	if pageUvOld > pageUv {
		// 删除页面总访问用户
		diff = pageUvOld - pageUv
		_, _ = _redis.Do("SPOP", pageUvKey, diff)
	} else if pageUvOld < pageUv {
		// 增加页面总访问用户
		diff = pageUv - pageUvOld
		for i := 0; i < diff; i++ {
			_, _ = _redis.Do("SADD", pageUvKey, tool.Md5(tool.RandString(32)))
		}
	}
	count = count + diff

	return count
}
