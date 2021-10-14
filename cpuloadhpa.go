package cpuloadhpa

import (
	"context"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/sirupsen/logrus"
)

type cl struct {
	sleepTime     time.Duration
	targetPercent int
	ctx           context.Context
}

func NewPayloadPercent(ctx context.Context, percent int) *cl {
	return &cl{
		sleepTime:     time.Millisecond,
		targetPercent: percent,
		ctx:           ctx,
	}
}

func (c *cl) monitorSleepTime() {
	l, r := time.Nanosecond, time.Second*10
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			cpus, err := cpu.Percent(time.Millisecond*1215, false)
			if err != nil {
				logrus.Errorf("Err:%v", err)
			}
			percent := cpus[0]
			mid := (l + r) / 2
			c.sleepTime = mid
			logrus.Infof("sleep time : %v", c.sleepTime)
			if percent > float64(c.targetPercent) {
				l = mid
			} else if percent <= float64(c.targetPercent) {
				r = mid
			}
		}

	}
}

func (c *cl) AsyncRun() {
	go c.monitorSleepTime()

	goroutines := runtime.GOMAXPROCS(0) / 2
	logrus.Infof("goroutines set : %v", goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			for {
				select {
				case <-c.ctx.Done():
					return
				default:
					time.Sleep(c.sleepTime)
				}
			}
		}()
	}

}
