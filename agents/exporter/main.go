package main

import (
	"cloudview/agents/exporter/core/stats/cpu"
)

func main() {
	cpu.RefreshCpu()
}
