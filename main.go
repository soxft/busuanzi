package main

import (
	"github.com/soxft/busuanzi/process/redisutil"
	"github.com/soxft/busuanzi/process/webutil"
)

func main() {
	redisutil.Init()
	webutil.Init()
}
