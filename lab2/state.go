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
