package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/k0kubun/pp"
	"github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"

	"github.com/nakabonne/gosivy/diagnoser"
	"github.com/nakabonne/gosivy/pidfile"
)

const defaultScrapeInterval = time.Second

var (
	flagSet = flag.NewFlagSet("gosivy", flag.ContinueOnError)

	// Automatically populated by goreleaser during build
	version = "unversioned"
	commit  = "?"
	date    = "?"
)

type cli struct {
	debug          bool
	version        bool
	list           bool
	scrapeInterval time.Duration
	stdout         io.Writer
	stderr         io.Writer
}

func (c *cli) usage() {
	format := `Usage:
  gosivy [flags] <pid|host:port>

Flags:
%s
Examples:
  gosivy 8000
  gosivy host.xz:8080

Author:
  Ryo Nakao <ryo@nakao.dev>
`
	fmt.Fprintf(c.stderr, format, flagSet.FlagUsages())
}

func main() {
	c, err := parseFlags(os.Stdout, os.Stderr)
	if err != nil {
		os.Exit(0)
	}
	os.Exit(c.run(flagSet.Args()))
}

func parseFlags(stdout, stderr io.Writer) (*cli, error) {
	c := &cli{
		stdout: stdout,
		stderr: stderr,
	}
	flagSet.BoolVarP(&c.version, "version", "v", false, "Print the current version.")
	flagSet.BoolVar(&c.debug, "debug", false, "Run in debug mode.")
	flagSet.BoolVarP(&c.list, "list-processes", "l", false, "Show processes where gosivy agent runs on.")
	flagSet.DurationVar(&c.scrapeInterval, "scrape-interval", defaultScrapeInterval, "Interval to scrape from the agent. It must be >= 1s")
	flagSet.Usage = c.usage
	if err := flagSet.Parse(os.Args[1:]); err != nil {
		if !errors.Is(err, flag.ErrHelp) {
			fmt.Fprintln(c.stderr, err)
		}
		return nil, err
	}

	return c, nil
}

func (c *cli) validate() error {
	if c.scrapeInterval < time.Second {
		return fmt.Errorf(`"--scrape-interval" must be >= 1s`)
	}
	return nil
}

func (c *cli) run(args []string) int {
	if c.version {
		fmt.Fprintf(c.stderr, "version=%s, commit=%s, buildDate=%s, os=%s, arch=%s\n", version, commit, date, runtime.GOOS, runtime.GOARCH)
		return 0
	}
	if err := c.validate(); err != nil {
		fmt.Fprintln(c.stderr, err)
		return 1
	}
	if err := setLogger(nil, c.debug); err != nil {
		fmt.Fprintf(c.stderr, "failed to prepare for debugging: %v\n", err)
		return 1
	}
	if c.list {
		c.listProcesses()
		return 0
	}
	var addr *net.TCPAddr
	if len(args) <= 0 {
		// TODO: Make sure to use the process where the agent runs on.
		addr = &net.TCPAddr{}
		return 0
	} else {
		var err error
		addr, err = targetToAddr(args[0])
		if err != nil {
			fmt.Fprintf(c.stderr, "failed to convert args into addresses: %v\n", err)
			return 1
		}
	}
	if err := diagnoser.Run(addr, c.scrapeInterval); err != nil {
		fmt.Fprintf(c.stderr, "failed to start diagnoser: %s\n", err.Error())
		c.usage()
		return 1
	}
	return 0
}

// targetToAddr parses the target string (pid or host:port),
// and converts it into the address of a TCP end point.
func targetToAddr(target string) (*net.TCPAddr, error) {
	// The case of "host:port"
	if strings.Contains(target, ":") {
		var err error
		addr, err := net.ResolveTCPAddr("tcp", target)
		if err != nil {
			return nil, fmt.Errorf("couldn't parse dst address: %w", err)
		}
		return addr, nil
	}

	// The case of PID.
	// Find port by pid then, connect to local
	pid, err := strconv.Atoi(target)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse PID: %w", err)
	}
	port, err := pidfile.GetPort(pid)
	if err != nil {
		return nil, fmt.Errorf("couldn't get port for PID %v: %w", pid, err)
	}
	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:"+port)
	return addr, nil
}

// Makes a new file under the config directory only when debug use.
func setLogger(w io.Writer, debug bool) error {
	if !debug {
		logrus.SetOutput(ioutil.Discard)
		pp.SetDefaultOutput(ioutil.Discard)
		return nil
	}
	if w == nil {
		var err error
		cfgDir, err := pidfile.ConfigDir()
		if err != nil {
			return err
		}
		if err := os.MkdirAll(cfgDir, os.ModePerm); err != nil {
			return err
		}
		w, err = os.OpenFile(filepath.Join(cfgDir, "debug.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
	}
	logrus.SetOutput(w)
	logrus.SetLevel(logrus.TraceLevel)
	pp.SetDefaultOutput(w)
	return nil
}
