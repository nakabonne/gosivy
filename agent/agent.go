// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package agent

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/shirou/gopsutil/process"

	"github.com/nakabonne/gosivy/pidfile"
	"github.com/nakabonne/gosivy/stats"
)

const defaultAddr = "127.0.0.1:0"

var (
	mu       sync.Mutex
	portfile string
	listener net.Listener
	errLog   io.Writer
)

// Options is optional settings for the started agent.
type Options struct {
	// The address the agent will be listening at.
	// It must be in the form of "host:port".
	Addr string

	// The directory to store the configuration file,
	ConfigDir string

	// Where to emit the error log to. By default io.Stderr is used.
	ErrorLog io.Writer
}

// Listen starts the gosivy agent on a host process. It automatically
// cleans up resources if the running process receives an interrupt.
//
// Note that the agent exposes an endpoint via a TCP connection that
// can be used by any program on the system.
func Listen(opts Options) error {
	mu.Lock()
	defer mu.Unlock()
	errLog = opts.ErrorLog
	if errLog == nil {
		errLog = os.Stderr
	}

	if portfile != "" {
		return fmt.Errorf("gosivy agent already listening at: %v", listener.Addr())
	}

	gosivyDir := opts.ConfigDir
	if gosivyDir == "" {
		cfgDir, err := pidfile.ConfigDir()
		if err != nil {
			return err
		}
		gosivyDir = cfgDir
	}

	err := os.MkdirAll(gosivyDir, os.ModePerm)
	if err != nil {
		return err
	}
	gracefulShutdown()

	addr := opts.Addr
	if addr == "" {
		addr = defaultAddr
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	listener = ln
	port := listener.Addr().(*net.TCPAddr).Port
	portfile = fmt.Sprintf("%s/%d", gosivyDir, os.Getpid())
	err = ioutil.WriteFile(portfile, []byte(strconv.Itoa(port)), os.ModePerm)
	if err != nil {
		return err
	}

	go listen()
	return nil
}

// Close closes the agent, removing temporary files and closing the TCP listener.
// If no agent is listening, Close does nothing.
func Close() {
	mu.Lock()
	defer mu.Unlock()

	if portfile != "" {
		os.Remove(portfile)
		portfile = ""
	}
	if listener != nil {
		listener.Close()
	}
}

// gracefulShutdown enables to automatically clean up resources if the
// running process receives an interrupt.
func gracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		// cleanup the socket on shutdown.
		sig := <-c
		Close()
		ret := 1
		if sig == syscall.SIGTERM {
			ret = 0
		}
		os.Exit(ret)
	}()
}

func listen() {
	buf := make([]byte, 1)
	for {
		fd, err := listener.Accept()
		if err != nil {
			// TODO: Find better way to check for closed connection, see: https://golang.org/issues/4373.
			if !strings.Contains(err.Error(), "use of closed network connection") {
				fmt.Fprintf(errLog, "gosivy: %v\n", err)
			}
			if netErr, ok := err.(net.Error); ok && !netErr.Temporary() {
				break
			}
			continue
		}
		if _, err := fd.Read(buf); err != nil {
			fmt.Fprintf(errLog, "gosivy: %v\n", err)
			continue
		}
		if err := handle(fd, buf); err != nil {
			fmt.Fprintf(errLog, "gosivy: %v\n", err)
			continue
		}
		fd.Close()
	}
}

func handle(conn io.ReadWriter, msg []byte) error {
	// TODO: Make it singleton if possible.
	process, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return err
	}
	switch msg[0] {
	case stats.SignalMeta:
		meta := stats.Meta{
			GoMaxProcs: runtime.GOMAXPROCS(0),
			NumCPU:     runtime.NumCPU(),
		}
		if u, err := process.Username(); err == nil {
			meta.Username = u
		}
		if c, err := process.Cmdline(); err == nil {
			meta.Cmmand = c
		}
		b, err := json.Marshal(meta)
		if err != nil {
			return err
		}
		_, err = conn.Write(b)
		return err
	case stats.SignalStats:
		s := stats.Stats{
			Goroutines: runtime.NumGoroutine(),
		}
		runtime.ReadMemStats(&s.MemStats)
		if c, err := process.CPUPercent(); err == nil {
			s.CPUUsage = c
		}
		b, err := json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = conn.Write(b)
		return err
	default:
		return fmt.Errorf("unknown signal received: %b", msg[0])
	}
}
