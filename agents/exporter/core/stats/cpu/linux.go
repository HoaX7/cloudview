package cpu

import (
	"cloudview/agents/exporter/core/iox"
	"cloudview/agents/exporter/core/logging"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
)

const (
	cpuTicks  = 100
	cpuFields = 8
	cpuMax    = 100
	statFile  = "/proc/stat"
)

var (
	preTotal uint64
	preIdle  uint64
	limit    float64
	cores    uint64
	initOnce sync.Once
	usage    uint64
)

func initialize() {
	cpus, err := effectiveCpus()
	if err != nil {
		logging.Error(err.Error())
		return
	}

	cores = uint64(cpus)
	limit = float64(cpus)
	quota, err := cpuQuota()
	if err == nil && quota > 0 {
		if quota < limit {
			limit = quota
		}
	}
}

func RefreshCpu() uint64 {
	initOnce.Do(initialize)
	idle, total := systemCpuUsage()

	idleTicks := float64(idle - atomic.LoadUint64(&preIdle))
	totalTicks := float64(total - atomic.LoadUint64(&preTotal))
	atomic.StoreUint64(&preIdle, idle)
	atomic.StoreUint64(&preTotal, total)

	cpuUsage := uint64(100 * (totalTicks - idleTicks) / totalTicks)
	if cpuUsage > cpuMax {
		cpuUsage = cpuMax
	}
	atomic.StoreUint64(&usage, cpuUsage)
	return cpuUsage
}

func GetCpuUsage() uint64 {
	return atomic.LoadUint64(&usage)
}

func cpuQuota() (float64, error) {
	cg, err := currentCgroup()
	if err != nil {
		return 0, err
	}

	return cg.cpuQuota()
}

func cpuUsage() (uint64, error) {
	cg, err := currentCgroup()
	if err != nil {
		return 0, err
	}

	return cg.cpuUsage()
}

func effectiveCpus() (int, error) {
	cg, err := currentCgroup()
	if err != nil {
		return 0, err
	}

	return cg.effectiveCpus()
}

func systemCpuUsage() (idle, total uint64) {
	lines, err := iox.ReadTextLines(statFile, iox.WithoutBlank())
	if err != nil {
		return
	}
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := iox.ParseUint(fields[i])
				if err != nil {
					fmt.Println("Error: ", i, fields[i], err)
				}
				total += val // tally up all the numbers to get total ticks
				if i == 4 {  // idle is the 5th field in the cpu line
					idle = val
				}
			}
			return
		}
	}
	return
}
