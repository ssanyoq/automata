package main

func main() {
	// p := NewParser(NewLexer("a(a|b)*"))
	// nfa, err := p.BuildNFA()
	// if err != nil {
	// 	panic(err)
	// }
	nfa := ConcatAutomata(CharAutomata('a'), KleeneeAutomata(OrAutomata(CharAutomata('a'), CharAutomata('b'))))
	nfa.PrintAutomata()
	println("DFA:")
	dfa := GenerateDFA(nfa)
	dfa.PrintDFA()
	println("Minimized:")
	dfa.Minimize()
	dfa.PrintDFA()
}
