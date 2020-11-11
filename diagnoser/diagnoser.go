// Package diagnoser mainly provides two components, scraper and GUI
// for the process diagnosis.
package diagnoser

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
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
		return nil, fmt.Errorf("failed to dial TCP: %w", err)
	}

	// First up, fetch meta data of process,
	if _, err := conn.Write([]byte{stats.SignalMeta}); err != nil {
		return nil, err
	}
	reader := bufio.NewReader(conn)
	res, err := reader.ReadBytes(stats.Delimiter)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata: %w", err)
	}
	var meta stats.Meta
	if err := json.Unmarshal(res, &meta); err != nil {
		return nil, fmt.Errorf("failed to decode metadata: %w", err)
	}

	go func(ctx context.Context, ch chan<- *stats.Stats) {
		defer conn.Close()
		tick := time.NewTicker(d.scrapeInterval)
		defer tick.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				if conn == nil {
					conn, err = net.DialTCP("tcp", nil, d.addr)
					if err != nil {
						logrus.Errorf("failed to dial: %v", err)
						continue
					}
				}

				if _, err := conn.Write([]byte{stats.SignalStats}); err != nil {
					logrus.Errorf("failed to write into connection: %v", err)
					conn = nil
					continue
				}
				reader.Reset(conn)
				res, err := reader.ReadBytes(stats.Delimiter)
				if err != nil {
					logrus.Errorf("failed to read the response: %v", err)
					conn = nil
					continue
				}

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
