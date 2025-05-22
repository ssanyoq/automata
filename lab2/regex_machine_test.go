package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		matchee  string
		expected string
	}{
		{
			name:     "basic",
			pattern:  "a(a|b)*",
			matchee:  "abab",
			expected: "abab",
		},
		{
			name:     "email",
			pattern:  "[a-z]+%@(gmail|ya)%.(ru|com)",
			matchee:  "user@gmail.com",
			expected: "user@gmail.com",
		},
		{
			name:     "email 2",
			pattern:  "user%@(gmail|ya)%.(ru|com)",
			matchee:  "user@ya.ru",
			expected: "user@ya.ru",
		},
		{
			name:     "ranges",
			pattern:  "(abc){,3}",
			matchee:  "abcabcabcabcabcabc", // more than 3
			expected: "abcabcabc",
		},
		{
			name:     "ranges 2",
			pattern:  "(abc){3, 5}",
			matchee:  "abcabc", // 2
			expected: "",
		},
		{
			name:     "ranges 3",
			pattern:  "(abc){3, 5}",
			matchee:  "abcabcabcabc", // 4
			expected: "abcabcabcabc",
		},
		{
			name:     "ranges 4",
			pattern:  "(abc){,5}",
			matchee:  "abc",
			expected: "abc",
		},
		{
			name:     "ranges 5",
			pattern:  "(abc){,}",
			matchee:  "abcabcabcabcabcabc",
			expected: "abcabcabcabcabcabc",
		},
		{
			name:     "symbol range",
			pattern:  "[a-z]",
			matchee:  "f",
			expected: "f",
		},
		{
			name:     "positive closure",
			pattern:  "[a-z]+",
			matchee:  "anything but",
			expected: "anything",
		},
		{
			name:     "slashy",
			pattern:  "abc/def",
			matchee:  "abc",
			expected: "",
		},
		{
			name:     "slashy 2",
			pattern:  "abc/def",
			matchee:  "abcdef",
			expected: "abc",
		},
		{
			name:     "slashy 3",
			pattern:  "a*/a",
			matchee:  "aaaaaa", // 6
			expected: "aaaaa",  // 5
		},
	}
	for _, tc := range tests {
		if tc.name != "positive closure" {
			continue
		}
		t.Run(tc.name, func(t *testing.T) {
			rm := NewRegexMachine(tc.pattern)
			matched, err := rm.Match(tc.matchee)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, matched)
		})
	}
}
