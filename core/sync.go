package core

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/soxft/busuanzi/library/database"
	"github.com/soxft/busuanzi/process/redisutil"
)

// LoadDataFromSQLite 从SQLite加载数据到Redis
func LoadDataFromSQLite(ctx context.Context) error {
	// 加载站点统计数据
	var siteStats []struct {
		SiteUnique string
		SitePv     int64
		SiteUv     int64
	}
	err := database.DB.WithContext(ctx).Table("site_statistics").Find(&siteStats).Error
	if err != nil {
		return err
	}

	// 加载页面统计数据
	var pageStats []struct {
		SiteUnique string
		PathUnique string
		PagePv     int64
		PageUv     int64
	}
	err = database.DB.WithContext(ctx).Table("page_statistics").Find(&pageStats).Error
	if err != nil {
		return err
	}

	_redis := redisutil.RDB
	// 更新站点统计数据
	for _, stat := range siteStats {
		rk := RKeys{
			SitePvKey: "bsz:site_pv:" + stat.SiteUnique,
			SiteUvKey: "bsz:site_uv:" + stat.SiteUnique,
		}
		_redis.Set(ctx, rk.SitePvKey, stat.SitePv, 0)
		_redis.PFAdd(ctx, rk.SiteUvKey, "1")
	}

	// 更新页面统计数据
	for _, stat := range pageStats {
		rk := RKeys{
			PagePvKey: "bsz:page_pv:" + stat.SiteUnique,
			PageUvKey: "bsz:page_uv:" + stat.SiteUnique,
		}
		_redis.ZAdd(ctx, rk.PagePvKey, redis.Z{Score: float64(stat.PagePv), Member: stat.PathUnique})
		_redis.PFAdd(ctx, rk.PageUvKey, "1")
	}
	return nil
}

// SyncToSQLite 将Redis数据同步到SQLite
func SyncToSQLite(ctx context.Context) error {
	_redis := redisutil.RDB
	// 获取所有需要同步的站点
	keys, err := _redis.Keys(ctx, "bsz:site_pv:*").Result()
	if err != nil {
		return err
	}

	for _, key := range keys {
		siteUnique := key[len("bsz:site_pv:"):]

		// 获取站点统计数据
		sitePv, _ := _redis.Get(ctx, "bsz:site_pv:"+siteUnique).Int64()
		siteUv, _ := _redis.PFCount(ctx, "bsz:site_uv:"+siteUnique).Result()

		// 更新站点统计数据
		err = database.DB.Exec(`
			INSERT INTO site_statistics (site_unique, site_pv, site_uv)
			VALUES (?, ?, ?)
			ON CONFLICT(site_unique) DO UPDATE SET
			site_pv = ?,
			site_uv = ?,
			updated_at = CURRENT_TIMESTAMP
		`, siteUnique, sitePv, siteUv, sitePv, siteUv).Error
		if err != nil {
			log.Printf("同步站点数据失败: %v", err)
			continue
		}

		// 获取该站点的所有页面PV
		pageData, _ := _redis.ZRangeWithScores(ctx, "bsz:page_pv:"+siteUnique, 0, -1).Result()
		for _, data := range pageData {
			pathUnique := data.Member.(string)
			pagePv := int64(data.Score)

			// 更新页面统计数据
			pageUv, _ := _redis.PFCount(ctx, "bsz:page_uv:"+siteUnique+":"+pathUnique).Result()
			err = database.DB.Exec(`
				INSERT INTO page_statistics (site_unique, path_unique, page_pv, page_uv)
				VALUES (?, ?, ?, ?)
				ON CONFLICT(site_unique, path_unique) DO UPDATE SET
				page_pv = ?,
				page_uv = ?,
				updated_at = CURRENT_TIMESTAMP
			`, siteUnique, pathUnique, pagePv, pageUv, pagePv, pageUv).Error
			if err != nil {
				log.Printf("同步页面数据失败: %v", err)
			}
		}
	}
	return nil
}

// StartSyncTask 启动定时同步任务
func StartSyncTask(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := SyncToSQLite(ctx); err != nil {
					log.Printf("同步到SQLite失败: %v", err)
				}
			}
		}
	}()
}
