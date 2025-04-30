package main

type Transition interface {
	CheckAndGetState(r rune) *State
}

type RangeTransition struct {
	Transition
	from rune
	to   rune
	s    *State
}

func (t *RangeTransition) CheckAndGetState(r rune) *State {
	if t.from < r && r < t.to {
		return t.s
	}
	return nil
}

type AlphaTransition struct {
	Transition
	char rune
	s    *State
}

func (t *AlphaTransition) CheckAndGetState(r rune) *State {
	if r == t.char {
		return t.s
	}
	return nil
}

type EpsilonTransition struct {
	Transition
	s *State
}

func (t *EpsilonTransition) CheckAndGetState(r rune) *State {
	return t.s
}

type State struct {
	isAccepting bool
	transitions []Transition
}

type Automata struct {
	head *State
	tail *State
}

func (a *Automata) Duplicate() *Automata {
	if a == nil || a.head == nil {
		return nil
	}

	stateMap := make(map[*State]*State)

	cloneState := func(original *State) *State {
		if duplicate, exists := stateMap[original]; exists {
			return duplicate
		}

		duplicate := &State{
			isAccepting: original.isAccepting,
		}
		stateMap[original] = duplicate
		return duplicate
	}
	visited := make(map[*State]bool)
	var traverse func(*State)
	traverse = func(s *State) {
		if visited[s] {
			return
		}
		visited[s] = true

		cloneState(s)

		for _, trans := range s.transitions {
			switch t := trans.(type) {
			case *RangeTransition:
				cloneState(t.s)
				traverse(t.s)
			case *AlphaTransition:
				cloneState(t.s)
				traverse(t.s)
			case *EpsilonTransition:
				cloneState(t.s)
				traverse(t.s)
			}
		}
	}
	traverse(a.head)

	for original, duplicate := range stateMap {
		for _, trans := range original.transitions {
			switch t := trans.(type) {
			case *RangeTransition:
				duplicate.transitions = append(duplicate.transitions, &RangeTransition{
					from: t.from,
					to:   t.to,
					s:    stateMap[t.s],
				})
			case *AlphaTransition:
				duplicate.transitions = append(duplicate.transitions, &AlphaTransition{
					char: t.char,
					s:    stateMap[t.s],
				})
			case *EpsilonTransition:
				duplicate.transitions = append(duplicate.transitions, &EpsilonTransition{
					s: stateMap[t.s],
				})
			}
		}
	}

	return &Automata{
		head: stateMap[a.head],
		tail: stateMap[a.tail],
	}
}
