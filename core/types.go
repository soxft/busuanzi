package core

import (
	cmap "github.com/orcaman/concurrent-map/v2"
	"time"
)

// index		数据类型	        key
// sitePv	Str	            bsz:site_pv:md5(host)
// siteUv	HyperLogLog		bsz:site_uv:md5(host)
// pagePv	zset	        bsz:page_pv:md5(host) / md5(path)
// pageUv	HyperLogLog		bsz:site_uv:md5(host):md5(path)

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

type expireQueue struct {
	Queue chan string
	Cache cmap.ConcurrentMap[string, time.Time]
}
