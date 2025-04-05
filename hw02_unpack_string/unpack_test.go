package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

// test ok --------------------------------------------------------

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
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
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

func TestUnpack_Additional(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// gitlab
		{input: "abcd", expected: "abcd"},
		{input: "d\n5abc", expected: "d\n\n\n\n\nabc"},
		// self
		{input: "a", expected: "a"},
		{input: "务", expected: "务"},
		{input: "\a", expected: "\a"},
		{input: "a🙃0", expected: "a"},
		{input: "aa🙃1", expected: "aa🙃"},
		{input: "aa1🙃", expected: "aa🙃"},
		{input: "\aq", expected: "\aq"},
		{input: "q\a", expected: "q\a"},
		{input: "\aqw", expected: "\aqw"},
		{input: "q\aw", expected: "q\aw"},
		{input: "qw\a", expected: "qw\a"},
		// uncomment if task with asterisk completed
		{input: `\5ab`, expected: `5ab`},
		{input: `a\5b`, expected: `a5b`},
		{input: `ab\5`, expected: `ab5`},
		{input: `\\`, expected: `\`},
		{input: `\\a`, expected: `\a`},
		{input: `a\\`, expected: `a\`},
		{input: `\\ab`, expected: `\ab`},
		{input: `a\\b`, expected: `a\b`},
		{input: `ab\\`, expected: `ab\`},
		{input: `q5we2\\5a`, expected: `qqqqqwee\\\\\a`},
		{input: `务\\许2可\\\\证0\1a🙃4\00\24`, expected: `务\许许可\\1a🙃🙃🙃🙃2222`},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString_Additional(t *testing.T) {
	invalidStrings := []string{
		"3",
		"\\",
		"d\\n5abc",
		"d\\\n5abc",
		"qwen\\",
		`qw\ne`,
		`\qwne`,
		`qwe\\55`,
		`\`,
		`\a`,
		`a\`,
		`\ab`,
		`a\b`,
		`ab\`,
	}
	for _, tc := range invalidStrings {
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}
