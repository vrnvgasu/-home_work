package main

import (
	"io"
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	src = "./testdata/input.txt"
	out = "./out.txt"
)

func TestCopyFailed(t *testing.T) {
	tests := []struct {
		name   string
		offset int64
	}{
		{
			name:   "fileSize < 0",
			offset: -1,
		},
		{
			name:   "offset > fileSize",
			offset: math.MaxInt64,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Copy(src, out, tt.offset, 0)
			require.Error(t, err)
		})
	}
}

func TestCopy(t *testing.T) {
	tests := []struct {
		name   string
		offset int64
		limit  int64
		dest   string
	}{
		{
			name:   "./go-cp -from testdata/input.txt -to out.txt",
			offset: 0,
			limit:  0,
			dest:   "testdata/out_offset0_limit0.txt",
		},
		{
			name:   "./go-cp -from testdata/input.txt -to out.txt -limit 10",
			offset: 0,
			limit:  10,
			dest:   "testdata/out_offset0_limit10.txt",
		},
		{
			name:   "./go-cp -from testdata/input.txt -to out.txt -limit 1000",
			offset: 0,
			limit:  1000,
			dest:   "testdata/out_offset0_limit1000.txt",
		},
		{
			name:   "./go-cp -from testdata/input.txt -to out.txt -limit 10000",
			offset: 0,
			limit:  10000,
			dest:   "testdata/out_offset0_limit10000.txt",
		},
		{
			name:   "./go-cp -from testdata/input.txt -to out.txt -offset 100 -limit 1000",
			offset: 100,
			limit:  1000,
			dest:   "testdata/out_offset100_limit1000.txt",
		},
		{
			name:   "./go-cp -from testdata/input.txt -to out.txt -offset 6000 -limit 1000",
			offset: 6000,
			limit:  1000,
			dest:   "testdata/out_offset6000_limit1000.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Remove(out)

			err := Copy(src, out, tt.offset, tt.limit)
			assert.NoError(t, err)

			orig, err := os.Open(tt.dest)
			assert.NoError(t, err)
			defer orig.Close()

			copied, err := os.Open(out)
			assert.NoError(t, err)
			defer copied.Close()

			bytesIn, err := io.ReadAll(orig)
			assert.NoError(t, err)

			bytesOut, err := io.ReadAll(copied)
			assert.NoError(t, err)

			assert.Equal(t, bytesIn, bytesOut)
		})
	}
}
