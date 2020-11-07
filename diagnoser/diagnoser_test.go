package diagnoser

import (
	"context"
	"encoding/json"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/nakabonne/gosivy/stats"
)

func TestStartScraping(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	addr := startServer()

	_, err := startScraping(ctx, addr, time.Microsecond, make(chan<- *stats.Stats, 100))
	time.Sleep(100 * time.Millisecond)
	assert.Nil(t, err)
}

func startServer() *net.TCPAddr {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	go func() {
		defer ln.Close()
		for {
			conn, _ := ln.Accept()
			sig := make([]byte, 1)
			_, _ = conn.Read(sig)
			switch sig[0] {
			case stats.SignalMeta:
				b, _ := json.Marshal(&stats.Meta{})
				_, _ = conn.Write(b)
			case stats.SignalStats:
				b, _ := json.Marshal(&stats.Stats{})
				_, _ = conn.Write(b)
			}
			conn.Close()
		}
	}()
	return ln.Addr().(*net.TCPAddr)
}
