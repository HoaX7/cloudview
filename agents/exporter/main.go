package main

import (
	"cloudview/agents/exporter/core/logging"
	"cloudview/agents/exporter/core/stats/cpu"
	"cloudview/agents/exporter/core/stats/memory"
	"cloudview/agents/exporter/core/stats/sysinfo"
	"fmt"
	"sync"
	"time"
)

var initOnce sync.Once

func main() {
	initOnce.Do(initialize)
	sysinfo, _ := sysinfo.GetSysInfo()
	fmt.Println("Exporter active and monitoring cpu usage: see 'usage.log'")
	fmt.Println("Logging errors to 'error.log'")
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case tm := <-ticker.C:
			runSafe(func() {
				usage := cpu.RefreshCpu()
				msg := fmt.Sprintf("Cpu Usage: %d%%", usage)
				logging.Log(msg)

				memUsage := memory.RefreshMemoryUsage()
				msg = fmt.Sprintf("Memory Usage: %d%%", memUsage)
				logging.Log(msg)

				// send usage data to cloudview backend
				if staticConfig.Reporting {
					cpuUsage := &Usage{
						Type:    "cpu",
						Percent: int(usage),
					}
					memUsage := &Usage{
						Type:    "ram",
						Percent: memUsage,
					}
					sysinfo.Timestamp = tm
					go reportMetrics(sysinfo, *cpuUsage, *memUsage)
				}
			})
		}
	}
}
