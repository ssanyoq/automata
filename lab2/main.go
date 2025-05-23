package main

func main() {
	p := NewParser(NewLexer("a(a|b)*"))
	nfa, _ := p.BuildNFA()
	dfa := GenerateDFA(nfa)
	_ = dfa.Minimize()
	dfa.PrintDFA()
	println("\nIntersection")
	dfa.Intersect(dfa).PrintDFA()
	println("\nSubtraction")
	dfa.Subtract(dfa).PrintDFA()
	println("\nRecovering")
	println(dfa.RecoverRegex())
}
