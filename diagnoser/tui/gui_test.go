package tui

import (
	"context"
	"fmt"
	"testing"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/stretchr/testify/assert"

	"github.com/nakabonne/gosivy/stats"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name    string
		r       runner
		wantErr bool
	}{
		{
			name: "successful running",
			r: func(context.Context, terminalapi.Terminal, *container.Container, ...termdash.Option) error {
				return nil
			},
			wantErr: false,
		},
		{
			name: "failed running",
			r: func(context.Context, terminalapi.Terminal, *container.Container, ...termdash.Option) error {
				return fmt.Errorf("error")
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			g := NewTUI(0, cancel, nil, &stats.Meta{})
			err := g.run(ctx, &termbox.Terminal{}, tt.r)
			assert.Equal(t, tt.wantErr, err != nil)
			cancel()
		})
	}
}
