package main

func main() {
	p := NewParser(NewLexer("a(a|b)*"))
	nfa, _ := p.BuildNFA()
	// nfa.PrintAutomata()
	dfa := GenerateDFA(nfa)
	// dfa.PrintDFA()
	_ = dfa.Minimize()
	dfa.PrintDFA()
	println("intersection")
	dfa.Intersect(dfa).PrintDFA()
	println("subtraction")
	dfa.Subtract(dfa).PrintDFA()
	// dfa.PrintDFA()
	// re := NewRegexMachine("aboba%.")
	// fmt.Println(re.Match("aboba."))
}
