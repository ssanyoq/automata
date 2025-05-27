package main

import (
	"fmt"

	"github.com/ssanyoq/automata-uni/lab2/util"
)

type State struct {
	isAccepting bool
	transitions []Transition
}

type Transition interface {
	CheckAndGetState(r rune) *State
	GetState() *State
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

func (t *RangeTransition) GetState() *State {
	return t.s
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

func (t *AlphaTransition) GetState() *State {
	return t.s
}

type EpsilonTransition struct {
	Transition
	s *State
}

func (t *EpsilonTransition) CheckAndGetState(r rune) *State {
	return t.s
}

func (t *EpsilonTransition) GetState() *State {
	return t.s
}

type Automata struct {
	head *State
	tail *State
}

// Returns copy of this automata
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

// Recursively calculates ε-closure of given state
func epsilonClosureRecursive(s *State, beenTo map[*State]bool) []*State {
	out := []*State{s}
	for _, tr := range s.transitions {
		etr, ok := tr.(*EpsilonTransition)
		if !ok {
			continue
		}
		_, ok = beenTo[etr.s]
		if ok {
			continue
		}
		beenTo[etr.s] = true
		out = append(out, epsilonClosureRecursive(etr.s, beenTo)...)
	}
	return out
}

// Calculates ε-closure of given state set
func EpsilonClosure(s []*State) []*State {
	beenTo := make(map[*State]bool)
	out := make([]*State, 0)
	for _, state := range s {
		beenTo[state] = true
		out = append(out, epsilonClosureRecursive(state, beenTo)...)
	}
	return out
}

// prints the NFA structure for debugging
func (a *Automata) PrintAutomata() {
	fmt.Println("Automata States:")
	states := make(map[*State]int)
	last := 1
	q := util.NewQueue[*State]()
	q.Push(a.head)
	for !q.IsEmpty() {
		s, _ := q.Pop()
		_, ok := states[s]
		if ok {
			continue
		}
		states[s] = last
		last++
		for _, v := range s.transitions {
			q.Push(v.GetState())
		}
	}
	q = util.NewQueue[*State]()
	been := make(map[*State]bool)
	q.Push(a.head)
	for !q.IsEmpty() {
		s, _ := q.Pop()
		if been[s] {
			continue
		}
		been[s] = true
		fmt.Printf("State %d", states[s])
		if s.isAccepting {
			fmt.Print("[accepting]")
		}
		fmt.Print(":")
		for _, tr := range s.transitions {
			switch t := tr.(type) {
			case *EpsilonTransition:
				fmt.Printf("ε -> %d; ", states[t.s])
			case *AlphaTransition:
				fmt.Printf("%c -> %d; ", t.char, states[t.s])
			case *RangeTransition:
				fmt.Printf("%c-%c -> %d; ", t.from, t.to, states[t.s])
			}
			q.Push(tr.GetState())
		}
		fmt.Println()
	}
}
