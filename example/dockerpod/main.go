package main

import (
	"context"
	"math/rand"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/xiekeyi98/cpuloadhpa"
	_ "go.uber.org/automaxprocs"
)

func main() {
	cl := cpuloadhpa.NewPayloadPercent(context.Background(), 75, cpuloadhpa.WithGoroutineNums(runtime.GOMAXPROCS(0)))
	cl.AsyncRun()
	tic := time.NewTicker(time.Minute * 5)
	defer tic.Stop()
	for range tic.C {
		shouldTarget := 60 + rand.Int()%30 // [60,90]
		if cl.GetTarget() != shouldTarget {
			logrus.Infof("update target to %v ", shouldTarget)
			cl.UpdateTarget(shouldTarget)
		}

	}

	select {}
}
