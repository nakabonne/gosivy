package diagnoser

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/nakabonne/gosivy/diagnoser/gui"
	"github.com/nakabonne/gosivy/stats"
)

const defaultRequestInterval = time.Second

func Run(addr *net.TCPAddr) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	statsCh := make(chan *stats.Stats)
	meta, err := poll(ctx, addr, statsCh)
	if err != nil {
		return err
	}
	return gui.Run(meta, statsCh)
}

func poll(ctx context.Context, addr *net.TCPAddr, statsCh chan<- *stats.Stats) (*stats.Meta, error) {
	conn, err := net.DialTCP("tcp", nil, addr)
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
		tick := time.NewTicker(defaultRequestInterval)
		defer tick.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				// TODO: Reuse connections instead of creating each time.
				conn, err := net.DialTCP("tcp", nil, addr)
				if err != nil {
					logrus.Errorf("failed to create connection: %v", err)
					continue
				}

				buf := []byte{stats.SignalStats}
				if _, err := conn.Write(buf); err != nil {
					logrus.Errorf("failed to write: %v", err)
					continue
				}
				res, err := ioutil.ReadAll(conn)
				if err != nil {
					logrus.Errorf("failed to read: %v", err)
					continue
				}
				conn.Close()

				var stats stats.Stats
				if err := json.Unmarshal(res, &stats); err != nil {
					logrus.Errorf("failed to unmarshal stats: %v", err)
					continue
				}
				ch <- &stats
			}
		}
	}(ctx, statsCh)

	return &meta, nil
}
