package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/soxft/busuanzi/config"
	"github.com/soxft/busuanzi/core"
	"github.com/soxft/busuanzi/library/database"
	"github.com/soxft/busuanzi/process/redisutil"
	"github.com/soxft/busuanzi/process/webutil"
)

func main() {
	config.Init()
	redisutil.Init()

	core.InitExpire()

	// 初始化SQLite
	database.InitDB("./data/busuanzi.db")

	// 创建上下文用于控制后台任务
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 从SQLite加载数据到Redis
	if err := core.LoadDataFromSQLite(ctx); err != nil {
		log.Printf("从SQLite加载数据失败: %v", err)
	}

	// 启动定时同步任务（每2分钟同步一次）
	core.StartSyncTask(ctx, 5*time.Minute)

	// 设置优雅退出
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("正在关闭服务...")
		cancel()
		// 最后同步一次数据
		if err := core.SyncToSQLite(context.Background()); err != nil {
			log.Printf("最终数据同步失败: %v", err)
		}
		os.Exit(0)
	}()

	webutil.Init()
}
