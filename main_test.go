package main

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
