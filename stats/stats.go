package stats

import (
	"os"
	"runtime"

	"github.com/shirou/gopsutil/process"
)

// Stats represents the statistical data of the process at the time of measurement.
type Stats struct {
	// The number of goroutines that currently exist.
	Goroutines int
	// How many percent of the CPU time this process uses
	CPUUsage float64
	MemStats
}

// MemStats records statistics about the memory allocator.
type MemStats struct {
	HeapAlloc uint64
	HeapIdle  uint64
	HeapInuse uint64
}

func NewStats() (*Stats, error) {
	// TODO: Make it singleton if possible.
	process, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return nil, err
	}
	var cpuUsage float64
	if c, err := process.CPUPercent(); err == nil {
		cpuUsage = c
	}
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return &Stats{
		Goroutines: runtime.NumGoroutine(),
		CPUUsage:   cpuUsage,
		MemStats: MemStats{
			HeapAlloc: m.HeapAlloc,
			HeapIdle:  m.HeapIdle,
			HeapInuse: m.HeapInuse,
		},
	}, nil
}
