package agent

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nakabonne/gosivy/stats"
)

func TestListenAndClose(t *testing.T) {
	err := Listen(Options{})
	assert.Nil(t, err)
	Close()
	_, err = os.Stat(pidFile)

	assert.True(t, os.IsNotExist(err))
	assert.Empty(t, pidFile)
}

func TestHandle(t *testing.T) {
	tests := []struct {
		name    string
		signal  byte
		wantErr bool
	}{
		{
			name:    "signal meta received",
			signal:  stats.SignalMeta,
			wantErr: false,
		},
		{
			name:    "signal stats received",
			signal:  stats.SignalStats,
			wantErr: false,
		},
		{
			name:    "unknown signal received",
			signal:  byte(0x9),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := new(bytes.Buffer)
			err := handle(b, []byte{tt.signal})
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
