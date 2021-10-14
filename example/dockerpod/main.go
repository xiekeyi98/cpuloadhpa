package main

import (
	"context"

	"github.com/xiekeyi98/cpuloadhpa"
)

func main() {
	cl := cpuloadhpa.NewPayloadPercent(context.Background(), 50)
	cl.AsyncRun()
	select {}
}
