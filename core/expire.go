package core

import (
	"context"
	camp "github.com/orcaman/concurrent-map/v2"
	"github.com/soxft/busuanzi/config"
	"github.com/soxft/busuanzi/process/redisutil"
	"github.com/spf13/viper"
	"log"
	"time"
)

var expQ expireQueue
var expTime = time.Second * 20

// InitExpire init expireQueue
func InitExpire() {
	if viper.GetInt("bsz.expire") == 0 {
		return
	}

	expQ.Queue = make(chan string, 100)
	expQ.Cache = camp.New[time.Time]()

	defer func() {
		if err := recover(); err != nil {
			close(expQ.Queue)
			expQ.Cache.Clear()

			log.Printf("[ERROR] core.Init %s \n", err)

			InitExpire()
		}
	}()

	// 监听 Expire chan 事件 消费
	go func() {
		for t := range expQ.Queue {
			if config.DEBUG {
				log.Printf("[expire] Set %s expired with %d s \n", t, viper.GetDuration("bsz.expire"))
			}
			redisutil.RDB.Expire(context.Background(), t, viper.GetDuration("bsz.expire")*time.Second)
		}
	}()

	go func() {
		// 定期清理过期数据
		timer := time.NewTicker(time.Second * 1)

		for {
			select {
			case <-timer.C:
				for k, v := range expQ.Cache.Items() {
					if time.Now().Sub(v) > expTime {
						expQ.Cache.Remove(k)

						// log.Printf("[%s] expire left %.2f delete \n", k, time.Now().Sub(v).Seconds())
					}
				}
			}
		}
	}()
}

func isInList(str string) bool {
	if t, ok := expQ.Cache.Get(str); ok {

		var tNow = time.Now()
		var tSub = tNow.Sub(t)

		if tSub < expTime {
			if config.DEBUG {
				log.Printf("[expire] %s %.2f (%d -> %d) ignored \n", str, tSub.Seconds(), tNow.Unix(), t.Unix())
			}

			return true
		} else {
			if config.DEBUG {
				log.Printf("[expire] %s %.2f (%d -> %d) expired \n", str, tSub.Seconds(), tNow.Unix(), t.Unix())
			}
		}
	}

	return false
}

// setExpire
// @description set the expiration time of the key
func setExpire(key ...string) {
	if viper.GetInt("bsz.expire") == 0 {
		return
	}

	for _, k := range key {
		if isInList(k) {
			continue
		}

		expQ.Queue <- k

		expQ.Cache.Set(k, time.Now())
	}
}
