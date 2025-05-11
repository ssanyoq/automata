package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToDFA(t *testing.T) {
	nfa := ConcatAutomata(CharAutomata('a'), KleeneeAutomata(OrAutomata(CharAutomata('a'), CharAutomata('b'))))
	nfa.PrintAutomata()
	dfa := GenerateDFA(nfa)
	dfa.PrintDFA()
}

func getPatternTests() []struct {
	name     string
	pattern  string
	matchee  string
	expected string
} {
	return []struct {
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
			pattern:  "user%@(gmail|ya)%.com",
			matchee:  "user@gmail.com",
			expected: "user@gmail.com",
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
	}
}

func TestNonMinimized(t *testing.T) {

	for _, tt := range getPatternTests() {
		if tt.name != "rangee" {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(NewLexer(tt.pattern))
			nfa, err := p.BuildNFA()
			assert.NoError(t, err)
			dfa := GenerateDFA(nfa)
			assert.Equal(t, tt.expected, dfa.Match(tt.matchee))
		})
	}
}

func TestMinimized(t *testing.T) {
	for _, tt := range getPatternTests() {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(NewLexer(tt.pattern))
			nfa, err := p.BuildNFA()
			assert.NoError(t, err)
			dfa := GenerateDFA(nfa)
			assert.NoError(t, dfa.Minimize())
			assert.Equal(t, tt.expected, dfa.Match(tt.matchee))
		})
	}
}
