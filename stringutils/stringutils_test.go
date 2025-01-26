package stringutils

import (
    "testing",
    "strings"
)

func TestCleanInput(t *testing.T) {
    cases := []struct {
        input string
        expected []string
    }{
        {
            input:      "  hello worLD",
            expected:   []string{"hello", "world"},
        },
        {   
            input:      "",
            expected:   []string{},
        },
        {   
            input:      "test sentence1",
            expected:   []string{"test", "sentence1"},
        },
        {   
            input:      "AAA THIS IS A AAAANA AAAA    ",
            expected:   []string{"AAA", "THIS", "IS", "A", "AAAANA", "AAAA"},
        },
    }
    for _, c := range cases {
        actual := stringutils.cleanInput(c.input)
        l := len(actual)
        if l != len(strings.Fields(c)) {
            
        }
    }
}
