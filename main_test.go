package main

import (
	"bytes"
	"net"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/nakabonne/gosivy/diagnoser"
)

func TestRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name string
		cli  cli
		args []string
		want int
	}{
		{
			name: "emit version",
			cli: cli{
				version: true,
			},
			want: 0,
		},
		{
			name: "invalid option",
			cli: cli{
				scrapeInterval: time.Microsecond,
			},
			want: 1,
		},
		{
			name: "run with remote addr",
			cli: cli{
				scrapeInterval: time.Second,
			},
			args: []string{"localhost:8080"},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := new(bytes.Buffer)
			tt.cli.stderr = b
			tt.cli.stdout = b
			m := diagnoser.NewMockDiagnoser(ctrl)
			m.EXPECT().Run().AnyTimes()
			tt.cli.diagnoser = m
			got := tt.cli.run(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTargetToAddr(t *testing.T) {
	tests := []struct {
		name    string
		target  string
		want    *net.TCPAddr
		wantErr bool
	}{
		{
			name:   "remote mode",
			target: "localhost:8080",
			want: &net.TCPAddr{
				IP:   net.ParseIP("127.0.0.1"),
				Port: 8080,
				Zone: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := targetToAddr(tt.target)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}
