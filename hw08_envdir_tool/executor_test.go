package main

import (
	"testing"

	"github.com/stretchr/testify/require" //nolint
)

func TestRunCmd(t *testing.T) {
	// Place your code here
	t.Run("Empty cmd", func(t *testing.T) {
		env := make(Environment)
		code := RunCmd([]string{}, env)
		require.Equal(t, 127, code)
	})

	t.Run("Command is exist", func(t *testing.T) {
		env := make(Environment)
		code := RunCmd([]string{"cd", "..", "-l"}, env)
		require.Equal(t, 0, code)
	})
}
