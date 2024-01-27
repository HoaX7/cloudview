package main

import (
	"cloudview/agents/exporter/core/stats/cpu"
	"fmt"
	"time"
)

var cpuUsage int64

const beta = 0.95

func main() {
	cpuTicker := time.NewTicker(time.Millisecond * 250)
	defer cpuTicker.Stop()

	for {
		select {
		case tm := <-cpuTicker.C:
			usage := cpu.RefreshCpu()
			fmt.Println("Current Time:", tm, usage)
		}
	}
}
