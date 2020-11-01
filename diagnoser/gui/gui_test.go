package gui

import (
	"context"
	"testing"
	"time"

	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/stretchr/testify/assert"

	"github.com/nakabonne/gosivy/stats"
)

func TestRun(t *testing.T) {
	type fields struct {
		RedrawInterval time.Duration
		Cancel         context.CancelFunc
		StatsCh        <-chan *stats.Stats
		Metadata       stats.Meta
		widgets        *widgets
	}
	type args struct {
		ctx context.Context
		t   terminalapi.Terminal
		r   runner
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GUI{
				RedrawInterval: tt.fields.RedrawInterval,
				Cancel:         tt.fields.Cancel,
				StatsCh:        tt.fields.StatsCh,
				Metadata:       tt.fields.Metadata,
				widgets:        tt.fields.widgets,
			}
			err := g.run(tt.args.ctx, tt.args.t, tt.args.r)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
