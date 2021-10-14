package cpuloadhpa

import (
	"container/list"
	"context"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/sirupsen/logrus"
)

type cl struct {
	targetPercent int
	ctx           context.Context
	goroutineNums int

	latestCpuPercent float64
	workCancels      *list.List
}
type Options func(*cl)

func WithGoroutineNums(num int) Options {
	return func(c *cl) {
		if num <= 0 {
			num = 1
		}
		c.goroutineNums = num
	}
}
func NewPayloadPercent(ctx context.Context, percent int, options ...Options) *cl {
	res := &cl{
		targetPercent: percent,
		ctx:           ctx,
		goroutineNums: runtime.GOMAXPROCS(0) / 2,
		workCancels:   list.New().Init(),
	}
	for _, o := range options {
		o(res)
	}
	logrus.Infof("set goroutineNums:%v ,and target cpu: %v", res.goroutineNums, res.targetPercent)
	return res
}

func (c *cl) UpdateTarget(target int) {
	logrus.Infof("update target cpu: %v -> %v", c.targetPercent, target)
	c.targetPercent = target
}
func (c *cl) GetTarget() int {
	return c.targetPercent
}
func (c *cl) monitorCPU() {
	for {
		cpus, err := cpu.Percent(time.Second, false)
		if err != nil {
			logrus.Errorf("Err:%v", err)
		}
		percent := cpus[0]
		c.latestCpuPercent = percent

	}
}

func (c *cl) GetLatestCPUPercent() float64 {
	return c.latestCpuPercent
}
func (c *cl) GetWokers() int {
	return c.workCancels.Len()

}
func (c *cl) AsyncRun() {
	go c.monitorCPU()
	for i := 0; i < c.goroutineNums; i++ {
		c.runNewWorker()
	}
}

func (c *cl) killOneWorker() {

}
func (c *cl) runNewWorker() {
	logrus.Infof("new worker added.")

	go func() {
		runtime.LockOSThread()
		for {
			select {
			case <-c.ctx.Done():
				return
			default:
				for {
					unitHundresOfMicrosecond := 1000
					runMicrosecond := unitHundresOfMicrosecond * c.targetPercent
					sleepMicrosecond := unitHundresOfMicrosecond*100 - runMicrosecond
					begin := time.Now()
					for {
						if time.Now().Sub(begin) > time.Duration(runMicrosecond)*time.Microsecond {
							break
						}
					}
					// sleep
					time.Sleep(time.Duration(sleepMicrosecond) * time.Microsecond)
				}
			}
		}
	}()

}
