package cpuload

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
	tik := time.Tick(time.Second)
	for range tik {
		select {
		case <-c.ctx.Done():
			return
		default:
			logrus.Infof("sleep time : %v", c.sleepTime)
			cpus, err := cpu.Percent(time.Second, false)
			if err != nil {
				logrus.Errorf("Err:%v", err)
			}
			percent := cpus[0]
			if percent > float64(c.targetPercent) {
				c.sleepTime = c.sleepTime * 2
			} else if percent < float64(c.targetPercent) {
				c.sleepTime = c.sleepTime / 2
			}
		}

	}
}

func (c *cl) Run() {
	go c.monitorSleepTime()

	for i := 0; i < runtime.NumCPU()/2; i++ {
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
