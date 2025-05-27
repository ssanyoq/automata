package main

// shortcuts for making NFA, for parsing

// Makes NFA automata for single character transition
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

// Makes NFA automata for [from-to]
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

// Makes NFA automata for left|right
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

// Makes NFA automata for left.right
func ConcatAutomata(left *Automata, right *Automata) *Automata {
	if left == nil {
		return right
	}
	if right == nil {
		return left
	}
	left.tail.isAccepting = false
	left.tail.transitions = append(left.tail.transitions, right.head.transitions...) // squash 2 states

	return &Automata{
		head: left.head,
		tail: right.tail,
	}
}

// Makes NFA automata for a*
func KleeneeAutomata(a *Automata) *Automata {
	a.tail.isAccepting = false
	a.tail.transitions = append(a.tail.transitions, &EpsilonTransition{s: a.head})
	st2 := &State{
		isAccepting: true,
		transitions: []Transition{},
	}
	a.tail.transitions = append(a.tail.transitions, &EpsilonTransition{s: st2})
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

// Makes NFA automata for a{from,to}
func RepeatAutomata(a *Automata, from int, to int) *Automata {
	var res *Automata

	for i := 0; i < from; i++ {
		res = ConcatAutomata(res, a.Duplicate())
	}
	if to == -1 {
		duplicate := a.Duplicate()
		return ConcatAutomata(res, KleeneeAutomata(duplicate))
	}
	end := &State{
		isAccepting: true,
	}

	middle := &State{
		isAccepting: false,
		transitions: []Transition{
			&EpsilonTransition{s: end},
		},
	}
	if res == nil {
		res = &Automata{
			head: middle,
			tail: middle,
		}
	} else {
		res.tail.transitions = append(res.tail.transitions, &EpsilonTransition{s: middle})
		res.tail = middle
	}

	for i := from; i < to; i++ {
		res = ConcatAutomata(res, a.Duplicate())
		res.tail.isAccepting = false
		skip := &State{
			isAccepting: false,
			transitions: []Transition{
				&EpsilonTransition{s: end},
			},
		}
		res.tail.transitions = append(res.tail.transitions, &EpsilonTransition{s: skip})
	}
	res.tail.isAccepting = false
	res.tail = end
	return res
}

// Makes NFA automata for a+
func PositiveClosureAutomata(a *Automata) *Automata {
	return RepeatAutomata(a, 1, -1)
}
