package main

import (
	"context"

	"github.com/soxft/busuanzi/config"
	"github.com/soxft/busuanzi/core"
	"github.com/soxft/busuanzi/process/persist"
	"github.com/soxft/busuanzi/process/redisutil"
	"github.com/soxft/busuanzi/process/webutil"
)

func main() {
	config.Init()
	redisutil.Init()

	core.InitExpire()

	persist.Init(context.Background())

	webutil.Init()
}
