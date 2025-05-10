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

func TestMatch(t *testing.T) {

	nfa := ConcatAutomata(CharAutomata('a'), KleeneeAutomata(OrAutomata(CharAutomata('a'), CharAutomata('b'))))
	dfa := GenerateDFA(nfa)
	assert.NoError(t, dfa.Minimize())

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
			pattern:  "user%@(gmail|ya)",
			matchee:  "user@gmail.com",
			expected: "user@gmail",
		},
		{
			name:     "ranges",
			pattern:  "(abc){,1}",
			matchee:  "abcabcabcabcabcabc",
			expected: "abc",
		},
	}

	for _, tt := range tests {
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
