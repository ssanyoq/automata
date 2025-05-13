package main

import "fmt"

func main() {
	p := NewParser(NewLexer("a(a|b)*"))
	nfa, _ := p.BuildNFA()
	nfa.PrintAutomata()
	dfa := GenerateDFA(nfa)
	dfa.PrintDFA()
	_ = dfa.Minimize()
	dfa.PrintDFA()
	re := NewRegexMachine("aboba%.")
	fmt.Println(re.Match("aboba."))
}
