package main

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	defer os.Remove("out.txt")
	err := exec.Command("./test.sh").Run()
	require.NoError(t, err)
}
