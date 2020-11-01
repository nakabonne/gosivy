package process

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPad(t *testing.T) {
	tests := []struct {
		name  string
		s     string
		total int
		want  string
	}{
		{
			name:  "larger than total",
			s:     "abc",
			total: 2,
			want:  "abc",
		},
		{
			name:  "smaller than total",
			s:     "abc",
			total: 5,
			want:  "abc  ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pad(tt.s, tt.total)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		name string
		i    int
		j    int
		want int
	}{
		{
			name: "j is larger",
			i:    1,
			j:    2,
			want: 2,
		},
		{
			name: "i is larger",
			i:    2,
			j:    1,
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := max(tt.i, tt.j)
			assert.Equal(t, tt.want, got)
		})
	}
}
