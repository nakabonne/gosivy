// Package diagnoser mainly provides two components, scraper and TUI
// for the process diagnosis.
package diagnoser

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/nakabonne/gosivy/diagnoser/tui"
	"github.com/nakabonne/gosivy/stats"
)

type GUI interface {
	Run(context.Context) error
}

type Diagnoser interface {
	Run() error
}

type diagnoser struct {
	addr           *net.TCPAddr
	scrapeInterval time.Duration
	gui            GUI
}

func NewDiagnoser(addr *net.TCPAddr, scrapeInterval time.Duration, gui GUI) Diagnoser {
	return &diagnoser{
		addr:           addr,
		scrapeInterval: scrapeInterval,
		gui:            gui,
	}
}

// Run performs the scraper which periodically scrapes from the agent,
// and then draws charts to show the stats.
func (d *diagnoser) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	statsCh := make(chan *stats.Stats)
	meta, err := d.startScraping(ctx, statsCh)
	if err != nil {
		return err
	}
	if d.gui == nil {
		d.gui = tui.NewTUI(d.scrapeInterval, cancel, statsCh, meta)
	}
	return d.gui.Run(ctx)
}

func (d *diagnoser) startScraping(ctx context.Context, statsCh chan<- *stats.Stats) (*stats.Meta, error) {
	conn, err := net.DialTCP("tcp", nil, d.addr)
	if err != nil {
		return nil, err
	}
	// First up, fetch meta data of process,
	buf := []byte{stats.SignalMeta}
	if _, err := conn.Write(buf); err != nil {
		return nil, err
	}
	res, err := ioutil.ReadAll(conn)
	if err != nil {
		return nil, err
	}
	conn.Close()
	var meta stats.Meta
	if err := json.Unmarshal(res, &meta); err != nil {
		return nil, err
	}

	go func(ctx context.Context, ch chan<- *stats.Stats) {
		tick := time.NewTicker(d.scrapeInterval)
		defer tick.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				// TODO: Reuse connections instead of creating each time.
				conn, err := net.DialTCP("tcp", nil, d.addr)
				if err != nil {
					logrus.Errorf("failed to create connection: %v", err)
					continue
				}

				buf := []byte{stats.SignalStats}
				if _, err := conn.Write(buf); err != nil {
					logrus.Errorf("failed to write into connection: %v", err)
					continue
				}
				res, err := ioutil.ReadAll(conn)
				if err != nil {
					logrus.Errorf("failed to read the response: %v", err)
					continue
				}
				conn.Close()

				var stats stats.Stats
				if err := json.Unmarshal(res, &stats); err != nil {
					logrus.Errorf("failed to decode stats: %v", err)
					continue
				}
				ch <- &stats
			}
		}
	}(ctx, statsCh)

	return &meta, nil
}
