package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	tests := []struct {
		name     string
		cmd      []string
		expected int
	}{
		{
			name:     "is null",
			cmd:      nil,
			expected: 0,
		},
		{
			name:     "empty",
			cmd:      nil,
			expected: 0,
		},
		{
			name:     "echo 1",
			cmd:      []string{"echo", "1"},
			expected: 0,
		},
		{
			name:     "error",
			cmd:      []string{"dummy", "1"},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, RunCmd(tt.cmd, make(Environment)))
		})
	}
}
