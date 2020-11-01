package process

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/keybase/go-ps"
)

// Process represents an OS process.
type Process struct {
	PID        int
	Executable string
	// Full path to the executable.
	Path string
}

type Processes []Process

// String formats as:
//
// PID   Exec Path
// 15788 foo  /path/to/src/foo
// 14054 main /private/var/folders/sy/5rwqjr1j3kl5r3kwxfgntkfr0000gn/T/go-build227076651/b001/exe/main
func (ps Processes) String() string {
	var (
		b          strings.Builder
		pidTitle   = "PID"
		execTitle  = "Exec"
		pathTitle  = "Path"
		maxPIDLen  = len(pidTitle)
		maxExecLen = len(execTitle)
	)
	// Take the maximum length to align the width of each column.
	for _, p := range ps {
		maxPIDLen = max(maxPIDLen, len(strconv.Itoa(p.PID)))
		maxExecLen = max(maxExecLen, len(p.Executable))
	}

	b.WriteString(fmt.Sprintf("%s %s %s\n",
		pad(pidTitle, maxPIDLen),
		pad(execTitle, maxExecLen),
		pathTitle,
	))

	for _, p := range ps {
		b.WriteString(fmt.Sprintf("%s %s %s\n",
			pad(strconv.Itoa(p.PID), maxPIDLen),
			pad(p.Executable, maxExecLen),
			p.Path,
		))
	}
	return b.String()
}

// FindAll gives back all processes where the agent runs on.
func FindAll() (Processes, error) {
	processes := make(Processes, 0)
	ps, err := ps.Processes()
	if err != nil {
		return nil, err
	}
	for _, p := range ps {
		process, err := newProcess(p)
		if err != nil {
			continue
		}
		processes = append(processes, *process)
	}
	return processes, nil
}

// FindOne finds processes where the agent runs on and gives back the first one it found.
func FindOne() (*Process, error) {
	ps, err := ps.Processes()
	if err != nil {
		return nil, err
	}
	for _, p := range ps {
		process, err := newProcess(p)
		if err == nil {
			return process, nil
		}
	}
	return nil, fmt.Errorf("no process where the agent runs found")
}

func newProcess(p ps.Process) (*Process, error) {
	pid := p.Pid()
	if pid == 0 {
		return nil, fmt.Errorf("system process given")
	}
	pidfile, err := PIDFile(pid)
	if err != nil {
		return nil, fmt.Errorf("failed to find pid file: %w", err)
	}
	if _, err := os.Stat(pidfile); err != nil {
		return nil, fmt.Errorf("given process isn't the process where the agent doesn't run on")
	}
	path, err := p.Path()
	if err != nil {
		return nil, fmt.Errorf("failed to detect full path to the executable: %w", err)
	}
	return &Process{
		PID:        pid,
		Executable: p.Executable(),
		Path:       path,
	}, nil

}

func pad(s string, total int) string {
	if len(s) >= total {
		return s
	}
	return s + strings.Repeat(" ", total-len(s))
}

func max(i, j int) int {
	if i > j {
		return i
	}
	return j
}
