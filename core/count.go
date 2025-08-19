package core

import (
	"context"
	"github.com/soxft/busuanzi/library/tool"
	"github.com/soxft/busuanzi/process/dbutil"
	"github.com/spf13/viper"
	"strings"
)

// Count
// @description return and count the number of users in the database
func Count(ctx context.Context, host string, path string, userIdentity string) Counts {
	rk := getKeys(host, path)

	// sitePV and pagePV
	sitePv, _ := dbutil.DB.IncrSitePv(ctx, rk.SiteUnique)
	pagePv, _ := dbutil.DB.IncrPagePv(ctx, rk.SiteUnique, rk.PathUnique)

	// siteUv and pageUv
	dbutil.DB.AddSiteUv(ctx, rk.SiteUnique, userIdentity)
	dbutil.DB.AddPageUv(ctx, rk.SiteUnique, rk.PathUnique, userIdentity)

	// count siteUv and pageUv
	siteUv, _ := dbutil.DB.CountSiteUv(ctx, rk.SiteUnique)
	pageUv, _ := dbutil.DB.CountPageUv(ctx, rk.SiteUnique, rk.PathUnique)

	// setExpire (only for Redis)
	go setExpire(rk.SiteUvKey, rk.PageUvKey, rk.SitePvKey, rk.PagePvKey)

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
	rk := getKeys(host, path)

	// sitePV and pagePV
	dbutil.DB.IncrSitePv(ctx, rk.SiteUnique)
	dbutil.DB.IncrPagePv(ctx, rk.SiteUnique, rk.PathUnique)

	// siteUv and pageUv
	dbutil.DB.AddSiteUv(ctx, rk.SiteUnique, userIdentity)
	dbutil.DB.AddPageUv(ctx, rk.SiteUnique, rk.PathUnique, userIdentity)

	// setExpire (only for Redis)
	go setExpire(rk.SiteUvKey, rk.PageUvKey, rk.SitePvKey, rk.PagePvKey)
	return
}

// Get bsz counts
func Get(ctx context.Context, host string, path string) Counts {
	rk := getKeys(host, path)

	// sitePV and pagePV
	sitePv, _ := dbutil.DB.GetSitePv(ctx, rk.SiteUnique)
	pagePv, _ := dbutil.DB.GetPagePv(ctx, rk.SiteUnique, rk.PathUnique)

	// count siteUv and pageUv
	siteUv, _ := dbutil.DB.CountSiteUv(ctx, rk.SiteUnique)
	pageUv, _ := dbutil.DB.CountPageUv(ctx, rk.SiteUnique, rk.PathUnique)

	// setExpire (only for Redis)
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

	redisPrefix := viper.GetString("database.prefix")

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
