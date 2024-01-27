package main

import (
	"cloudview/agents/exporter/core/stats/cpu"
	"fmt"
	"sync/atomic"
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
			curUsage := cpu.RefreshCpu()
			prevUsage := atomic.LoadInt64(&cpuUsage)
			// cpu = cpuᵗ⁻¹ * beta + cpuᵗ * (1 - beta)
			usage := int64(float64(prevUsage)*beta + float64(curUsage)*(1-beta))
			atomic.StoreInt64(&cpuUsage, usage)
			fmt.Println("Current Time:", tm, usage)
		}
	}
}
