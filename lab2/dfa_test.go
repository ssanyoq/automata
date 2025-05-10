package main

import "testing"

func TestToDFA(t *testing.T) {
	nfa := ConcatAutomata(CharAutomata('a'), KleeneeAutomata(OrAutomata(CharAutomata('a'), CharAutomata('b'))))
	nfa.PrintAutomata()
	dfa := GenerateDFA(nfa)
	dfa.PrintDFA()
}
