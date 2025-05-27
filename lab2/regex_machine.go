package main

import "errors"

// User interface for regex library
type RegexMachine struct {
	pattern    string
	automata   *DFAResult
	prognostic *DFAResult
}

// Creates new instance of RegexMachine
func NewRegexMachine(pattern string) *RegexMachine {
	return &RegexMachine{
		pattern: pattern,
	}
}

// Compiles pattern into minimized DFA
func (r *RegexMachine) Compile() error {
	prognosticIndex := -1
	for i, v := range r.pattern {
		if v == '/' {
			if prognosticIndex != -1 {
				return errors.New("found more than 1 prognostic operator")
			}
			prognosticIndex = i
		}
	}
	reString := r.pattern
	if prognosticIndex != -1 {
		reString = r.pattern[:prognosticIndex]

		p := NewParser(NewLexer(r.pattern[prognosticIndex+1:]))
		nfa, err := p.BuildNFA()
		if err != nil {
			return err
		}
		r.prognostic = GenerateDFA(nfa)
		err = r.prognostic.Minimize()
		if err != nil {
			return nil
		}
	}
	p := NewParser(NewLexer(reString))
	nfa, err := p.BuildNFA()
	if err != nil {
		return err
	}
	r.automata = GenerateDFA(nfa)
	return r.automata.Minimize()
}

// Returns longest match of matchee.
// Also compiles automata by pattern if not compiled yet
func (r *RegexMachine) Match(matchee string) (string, error) {
	if r.automata == nil {
		err := r.Compile()
		if err != nil {
			return "", err
		}
	}
	if r.prognostic == nil {
		return r.automata.Match(matchee), nil
	}
	var maxMatch string
	for i := range matchee {
		reMatch := r.automata.Match(matchee[:i])
		progMatch := r.automata.Match(matchee[i:])
		if progMatch != "" && reMatch != "" {
			maxMatch = reMatch
		}
	}
	return maxMatch, nil
}

// Recovers pattern from
func (r *RegexMachine) RecoverPattern() (string, error) {
	if r.automata == nil {
		err := r.Compile()
		if err != nil {
			return "", err
		}
	}
	return "", errors.New("not implemented")
}

// Returns intersection of two given RegexMachines
func (r *RegexMachine) Intersect(other *RegexMachine) (*RegexMachine, error) {
	if r.automata == nil {
		err := r.Compile()
		if err != nil {
			return nil, err
		}
	}
	if other.automata == nil {
		err := other.Compile()
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

// Subtracts other RegexMachine from this one and returns the result
func (r *RegexMachine) Subtract(other *RegexMachine) (*RegexMachine, error) {
	if r.automata == nil {
		err := r.Compile()
		if err != nil {
			return nil, err
		}
	}
	if other.automata == nil {
		err := other.Compile()
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}
