package ScheduledTask

import (
	"github.com/robfig/cron/v3"
	"golang/service"
	"golang/service/app"
	"runtime"
)

var (
	c *cron.Cron
)

// Cron在线表达式生成器
// http://cron.ciding.cc/

func InitCorn() {
	c = cron.New(cron.WithSeconds())
	c.Start()
}
func InitCornTask() {
	if c == nil {
		InitCorn()
	}

	c.AddFunc("0 */5 * * * ? ", func() {
		runtime.GC()
	})
	c.AddFunc("*/20 * * * * ? ", func() {
		// 更新网卡流量 - 守护线程
		if service.NetUpdateThreadIsRun == false {
			go service.NetUpdateThread()
		}
	})

	//看看TestSpeed队列有没有排队的
	c.AddFunc("*/1 * * * * ? ", func() {
		defer app.TestListLock.RUnlock()
		app.TestListLock.RLock()

		if len(app.TestList) > 0 {
			if app.TestList[0].IsRun == false {
				app.TestList[0].IsRun = true
				go app.TestList[0].RunNow()
			}
		}
	})
}
