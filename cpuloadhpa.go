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
	goroutineNums int
}
type Options func(*cl)

func WithGoroutineNums(num int) Options {
	return func(c *cl) {
		c.goroutineNums = num
	}
}
func NewPayloadPercent(ctx context.Context, percent int, options ...Options) *cl {
	res := &cl{
		sleepTime:     time.Millisecond,
		targetPercent: percent,
		ctx:           ctx,
		goroutineNums: runtime.GOMAXPROCS(0) / 2,
	}
	for _, o := range options {
		o(res)
	}
	logrus.Infof("set goroutineNums:%v ,and target cpu: %v", res.goroutineNums, res.targetPercent)
	return res
}

func (c *cl) UpdateTarget(target int) {
	logrus.Infof("update target cpu: %v", c.targetPercent)
	c.targetPercent = target
}
func (c *cl) GetTarget() int {
	return c.targetPercent
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
			//logrus.Infof("sleep time : %v", c.sleepTime)
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

	for i := 0; i < c.goroutineNums; i++ {
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
