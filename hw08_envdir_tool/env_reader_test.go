package main

import (
	"testing"

	"github.com/stretchr/testify/require" //nolint
)

func TestReadDir(t *testing.T) {
	// Place your code here
	t.Run("The directory does not exist", func(t *testing.T) {
		_, err := ReadDir("testdata/env1/")
		require.Error(t, err)
	})

	t.Run("Read directory", func(t *testing.T) {
		env, _ := ReadDir("testdata/env/")
		require.Len(t, env, 5)
	})

	t.Run("Read directory empty / no empty file", func(t *testing.T) {
		env, _ := ReadDir("testdata/env/")
		require.True(t, env["BAR"].NeedRemove)
		require.False(t, env["EMPTY"].NeedRemove)
	})
}
