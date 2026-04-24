package hw02unpackstring

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/stretchr/testify/require" //nolint
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: "🙃0", expected: ""},
		{input: "aaф0b", expected: "aab"},
		// uncomment if task with asterisk completed
		// {input: `qwe\4\5`, expected: `qwe45`},
		// {input: `qwe\45`, expected: `qwe44444`},
		// {input: `qwe\\5`, expected: `qwe\\\\\`},
		// {input: `qwe\\\3`, expected: `qwe\3`},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b"}
	for _, tc := range invalidStrings {
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}

func TestAssertUnpackError(t *testing.T) {
	var out, err = Unpack("A45")
	assert.True(t, errors.Is(err, ErrInvalidString))
	assert.Equal(t, "", out)
}

func TestUnpackErrorF(t *testing.T) {
	const str, want = "a3bc", "aaabc"
	got, err := Unpack(str)

	if err != nil {
		t.Errorf("unpack returns unexpected error: %q", str)
	}

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
