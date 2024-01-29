package memory

import (
	"cloudview/agents/exporter/core/iox"
	"cloudview/agents/exporter/core/utils"
	"fmt"
	"regexp"
	"strings"
)

type MemStats map[string]uint64

const statFile = "/proc/meminfo"

func RefreshMemoryUsage() int {
	stats, err := getMemStats()
	if err != nil {
		return 0
	}
	return utils.GetPercentage(int(stats["memtotal"]), int(stats["memused"]))
}

// getMemStats gets the memory stats of a linux system from the
// file /proc/meminfo
func getMemStats() (memStats MemStats, err error) {
	lines, err := iox.ReadTextLines(statFile, iox.WithoutBlank())
	if err != nil {
		return
	}

	memStats = MemStats{}
	re := regexp.MustCompile(`^((?:Mem|Swap)(?:Total|Free)|Buffers|Cached|` +
		`SwapCached|Active|Inactive|Dirty|Writeback|Mapped|Slab|` +
		`Commit(?:Limit|ted_AS)):\s*(\d+)`)

	for _, line := range lines {
		stat := re.FindStringSubmatch(line)
		if stat == nil {
			continue
		}
		key := stat[1]
		value, err := iox.ParseUint(stat[2])
		if err != nil {
			fmt.Println(err)
			continue
		} else {
			memStats[strings.ToLower(key)] = value
		}
	}

	memStats[`memused`] = memStats[`memtotal`] - memStats[`memfree`]
	memStats[`swapused`] = memStats[`swaptotal`] - memStats[`swapfree`]
	memStats[`realfree`] = memStats[`memfree`] + memStats[`buffers`] + memStats[`cached`]

	return memStats, nil
}
