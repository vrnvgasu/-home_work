package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const dir = "./testdata/env"
		expected := Environment{
			"HELLO": EnvValue{
				Value:      "\"hello\"",
				NeedRemove: false,
			},
			"BAR": EnvValue{
				Value:      "bar",
				NeedRemove: false,
			},
			"FOO": EnvValue{
				Value:      "   foo\nwith new line",
				NeedRemove: false,
			},
			"UNSET": EnvValue{
				Value:      "",
				NeedRemove: true,
			},
			"EMPTY": EnvValue{
				Value:      "",
				NeedRemove: false,
			},
		}
		evv, err := ReadDir(dir)
		require.NoError(t, err)
		require.Equal(t, expected, evv)
	})
	t.Run("failure", func(t *testing.T) {
		evv, err := ReadDir("dummy")
		require.Error(t, err)
		require.Nil(t, evv)
	})
}

func TestEnvironmentHandle(t *testing.T) {
	err := os.Setenv("dummy_remove", "")
	require.NoError(t, err)
	err = os.Setenv("dummy_replace", "")
	require.NoError(t, err)

	env := Environment{
		"dummy_remove": {
			Value:      "",
			NeedRemove: true,
		},
		"dummy_replace": {
			Value:      "xxx",
			NeedRemove: false,
		},
	}
	err = env.Handle()
	require.NoError(t, err)

	_, ok := os.LookupEnv("dummy_remove")
	require.False(t, ok)
	v, ok := os.LookupEnv("dummy_replace")
	require.True(t, ok)
	require.Equal(t, env["dummy_replace"].Value, v)
}
