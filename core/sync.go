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
	type Statistic struct {
		SiteUnique string
		PathUnique string
		SitePv     int64
		SiteUv     int64
		PagePv     int64
		PageUv     int64
	}

	var statistics []Statistic
	err := database.DB.WithContext(ctx).Table("statistics").Find(&statistics).Error
	if err != nil {
		return err
	}

	for _, stat := range statistics {
		siteUnique, pathUnique := stat.SiteUnique, stat.PathUnique
		sitePv := stat.SitePv
		pagePv := stat.PagePv

		rk := RKeys{
			SitePvKey: "bsz:site_pv:" + siteUnique,
			SiteUvKey: "bsz:site_uv:" + siteUnique,
			PagePvKey: "bsz:page_pv:" + siteUnique,
			PageUvKey: "bsz:page_uv:" + siteUnique,
		}

		// 更新Redis数据
		_redis := redisutil.RDB
		_redis.Set(ctx, rk.SitePvKey, sitePv, 0)
		_redis.PFAdd(ctx, rk.SiteUvKey, "1")
		_redis.ZAdd(ctx, rk.PagePvKey, redis.Z{Score: float64(pagePv), Member: pathUnique})
		_redis.PFAdd(ctx, rk.PageUvKey, "1")
		// UV数据需要重新收集，这里只是占位
	}
	return nil
}

// SyncToSQLite 将Redis数据同步到SQLite
func SyncToSQLite(ctx context.Context) error {
	// 获取所有需要同步的键
	_redis := redisutil.RDB
	keys, err := _redis.Keys(ctx, "bsz:site_pv:*").Result()
	if err != nil {
		return err
	}

	for _, key := range keys {
		// 解析site_unique
		siteUnique := key[len("bsz:site_pv:"):]

		// 获取相关的所有统计数据
		sitePv, _ := _redis.Get(ctx, "bsz:site_pv:"+siteUnique).Int64()
		siteUv, _ := _redis.PFCount(ctx, "bsz:site_uv:"+siteUnique).Result()

		// 获取该站点的所有页面PV
		pageData, _ := _redis.ZRangeWithScores(ctx, "bsz:page_pv:"+siteUnique, 0, -1).Result()

		for _, data := range pageData {
			pathUnique := data.Member.(string)
			pagePv := int64(data.Score)
			pageUv, _ := _redis.PFCount(ctx, "bsz:page_uv:"+siteUnique+":"+pathUnique).Result()

			// 更新SQLite数据
			err = database.DB.Exec(`
				INSERT INTO statistics (site_unique, path_unique, site_pv, site_uv, page_pv, page_uv)
				VALUES (?, ?, ?, ?, ?, ?)
				ON CONFLICT(site_unique, path_unique) DO UPDATE SET
				site_pv = ?,
				site_uv = ?,
				page_pv = ?,
				page_uv = ?,
				updated_at = CURRENT_TIMESTAMP
			`, siteUnique, pathUnique, sitePv, siteUv, pagePv, pageUv,
				sitePv, siteUv, pagePv, pageUv).Error
			if err != nil {
				log.Printf("同步数据失败: %v", err)
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
