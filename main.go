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

	"github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"

	"github.com/nakabonne/gosivy/diagnoser"
	"github.com/nakabonne/gosivy/process"
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
	diagnoser      diagnoser.Diagnoser
}

func (c *cli) usage() {
	format := `Usage:
  gosivy [flags] <pid|host:port>

Flags:
%s
Examples:
  gosivy 15788
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
	code := c.run(flagSet.Args())
	if code != 0 {
		c.usage()
		os.Exit(code)
	}
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
		ps, err := process.FindAll()
		if err != nil {
			fmt.Fprintf(c.stderr, "failed to list processes: %v\n", err)
			return 1
		}
		fmt.Fprintf(c.stderr, "%v", ps)
		return 0
	}

	var pid string
	if len(args) == 0 {
		// Automatically finds the process where the agent runs on if no args given.
		p, err := process.FindOne()
		if err != nil {
			fmt.Fprintln(c.stderr, err)
			return 1
		}
		pid = strconv.Itoa(p.PID)
	} else {
		pid = args[0]
	}
	addr, err := targetToAddr(pid)
	if err != nil {
		fmt.Fprintf(c.stderr, "failed to convert args into addresses: %v\n", err)
		return 1
	}
	if c.diagnoser == nil {
		c.diagnoser = diagnoser.NewDiagnoser(addr, c.scrapeInterval, nil)
	}
	if err := c.diagnoser.Run(); err != nil {
		fmt.Fprintf(c.stderr, "failed to start diagnoser: %s\n", err.Error())
		return 1
	}
	return 0
}

func (c *cli) validate() error {
	if c.scrapeInterval < time.Second {
		return fmt.Errorf(`"--scrape-interval" must be >= 1s`)
	}
	return nil
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
	port, err := process.GetPort(pid)
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
		return nil
	}
	if w == nil {
		var err error
		cfgDir, err := process.ConfigDir()
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
	return nil
}
