package main

import (
	"context"

	"github.com/xiekeyi98/cpuload"
)

func main() {
	l := cpuload.NewPayloadPercent(context.Background(), 50)
	l.Run()
	select {}
}
