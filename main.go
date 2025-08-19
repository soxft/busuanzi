package main

import (
	"github.com/soxft/busuanzi/config"
	"github.com/soxft/busuanzi/core"
	"github.com/soxft/busuanzi/process/dbutil"
	"github.com/soxft/busuanzi/process/webutil"
)

func main() {
	config.Init()
	dbutil.Init()

	core.InitExpire()

	webutil.Init()
}
