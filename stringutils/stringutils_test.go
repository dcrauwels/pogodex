package stringutils

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello worLD",
			expected: []string{"hello", "world"},
		},
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    "test sentence1",
			expected: []string{"test", "sentence1"},
		},
		{
			input:    "AAA THIS IS A AAAANA AAAA    ",
			expected: []string{"aaa", "this", "is", "a", "aaaana", "aaaa"},
		},
	}
	for _, c := range cases {
		actual := CleanInput(c.input)
		l := len(actual)
		e := len(c.expected)
		if l != e {
			t.Errorf("CleanInput() slice length (%d) does not match expected length (%d)", l, e)
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("CleanInput() result %d (%s) does not match expected output (%s)", i, word, expectedWord)
			}
		}
	}
}
