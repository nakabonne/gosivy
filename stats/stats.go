package stats

import "runtime"

// Stats represents the statistical data of the process at the time of measurement.
type Stats struct {
	// The number of goroutines that currently exist.
	Goroutines int
	// How many percent of the CPU time this process uses
	CPUUsage float64
	MemStats runtime.MemStats
}
