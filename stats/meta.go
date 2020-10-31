package stats

import (
	"fmt"
	"os"
	"runtime"

	"github.com/shirou/gopsutil/process"
)

// Meta represents process metadata, which will be not changed
// as long as the process continues.
type Meta struct {
	PID        int
	Username   string
	Command    string
	GoMaxProcs int
	NumCPU     int
}

func NewMeta() (*Meta, error) {
	// TODO: Make it singleton if possible.
	process, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return nil, err
	}

	var username, command string
	if u, err := process.Username(); err == nil {
		username = u
	}
	if c, err := process.Cmdline(); err == nil {
		command = c
	}
	return &Meta{
		PID:        os.Getpid(),
		Username:   username,
		Command:    command,
		GoMaxProcs: runtime.GOMAXPROCS(0),
		NumCPU:     runtime.NumCPU(),
	}, nil
}

func (m *Meta) String() string {
	return fmt.Sprintf(
		"PID: %d, CMD: %s, User: %s, Num CPU: %d, GOPAXPROS: %d",
		m.PID,
		m.Command,
		m.Username,
		m.NumCPU,
		m.GoMaxProcs,
	)
}
