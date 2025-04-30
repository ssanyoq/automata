package main

// shortcuts for making NFA, for parsing

func CharAutomata(r rune) *Automata {
	st2 := State{
		isAccepting: true,
		transitions: []Transition{},
	}
	st1 := State{
		isAccepting: false,
		transitions: []Transition{
			&AlphaTransition{
				char: r,
				s:    &st2,
			},
		},
	}
	return &Automata{
		head: &st1,
		tail: &st2,
	}
}

func RangeCharAutomata(from rune, to rune) *Automata {
	st2 := State{
		isAccepting: true,
		transitions: []Transition{},
	}
	st1 := State{
		isAccepting: false,
		transitions: []Transition{
			&RangeTransition{
				from: from,
				to:   to,
				s:    &st2,
			},
		},
	}
	return &Automata{
		head: &st1,
		tail: &st2,
	}
}

func OrAutomata(left *Automata, right *Automata) *Automata {
	left.tail.isAccepting = false
	right.tail.isAccepting = false
	st1 := State{
		isAccepting: false,
		transitions: []Transition{
			&EpsilonTransition{
				s: left.head,
			},
			&EpsilonTransition{
				s: right.head,
			},
		},
	}
	st2 := State{
		isAccepting: true,
		transitions: []Transition{},
	}
	left.tail.transitions = append(left.tail.transitions, &EpsilonTransition{s: &st2})
	right.tail.transitions = append(right.tail.transitions, &EpsilonTransition{s: &st2})
	return &Automata{
		head: &st1,
		tail: &st2,
	}
}

func ConcatAutomata(left *Automata, right *Automata) *Automata {
	left.tail.isAccepting = false
	left.tail.transitions = append(left.tail.transitions, right.head.transitions...) // squash 2 states

	return &Automata{
		head: left.head,
		tail: right.tail,
	}
}

func KleeneeAutomata(a *Automata) *Automata {
	a.tail.isAccepting = false
	a.tail.transitions = append(a.tail.transitions, &EpsilonTransition{s: a.head})
	st2 := &State{
		isAccepting: true,
		transitions: []Transition{},
	}
	st1 := &State{
		isAccepting: false,
		transitions: []Transition{
			&EpsilonTransition{
				s: a.head,
			},
			&EpsilonTransition{
				s: st2,
			},
		},
	}
	return &Automata{
		head: st1,
		tail: st2,
	}
}

func PositiveClosureAutomata(a *Automata) *Automata {
	// TODO:
	return nil
}

func RepeatAutomata(a *Automata, from int, to int) *Automata {
	// TODO:
	return nil
}
