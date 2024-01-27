package cpu

import (
	"cloudview/agents/exporter/core/iox"
	"errors"
	"fmt"
	"strings"
	"sync"
)

const (
	cpuTicks  = 100
	cpuFields = 8
	cpuMax    = 1000
	statFile  = "./test"
)

var (
	preSystem uint64
	prevReads []uint64
	limit     float64
	cores     uint64
	initOnce  sync.Once
)

func initialize() {
	preSystem, _, prevReads = systemCpuUsage()
	fmt.Println(prevReads)
}

func RefreshCpu() uint64 {
	initOnce.Do(initialize)
	system, _, curReads := systemCpuUsage()
	delta := system - preSystem
	idle := curReads[4] - prevReads[4]
	cpu_used := delta - idle
	usage := 100 * cpu_used / idle
	preSystem = system
	prevReads = curReads
	return usage
}

func systemCpuUsage() (uint64, error, []uint64) {
	lines, err := iox.ReadTextLines(statFile, iox.WithoutBlank())
	if err != nil {
		return 0, err, []uint64{}
	}

	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			if len(fields) < cpuFields {
				return 0, fmt.Errorf("bad format of cpu stats"), []uint64{}
			}

			rt, _ := parseUints(strings.Join(fields[1:cpuFields], ","))
			var totalClockTicks uint64
			for _, i := range fields[1:cpuFields] {
				v, err := parseUint(i)
				if err != nil {
					return 0, err, rt
				}

				totalClockTicks += v
			}

			return totalClockTicks, nil, rt
		}
	}

	return 0, errors.New("bad stats format"), []uint64{}
}
