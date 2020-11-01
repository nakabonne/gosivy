package diagnoser

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/nakabonne/gosivy/stats"
)

func TestStartScraping(t *testing.T) {
	type args struct {
		ctx      context.Context
		addr     *net.TCPAddr
		interval time.Duration
		statsCh  chan<- *stats.Stats
	}
	tests := []struct {
		name    string
		args    args
		want    *stats.Meta
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}
