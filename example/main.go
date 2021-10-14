package main

import (
	"context"

	"github.com/xiekeyi98/cpuloadhpa"
)

func main() {
	l := cpuloadhpa.NewPayloadPercent(context.Background(), 50)
	l.Run()
	select {}
}
