package main

import (
	"fmt"

	"github.com/ssanyoq/automata-uni/lab2/util"
)

// DFAState represents a state in the DFA, which is a set of NFA states
type DFAState struct {
	nfaStates   map[*State]bool // Set of NFA states this DFA state represents
	isAccepting bool
	transitions map[rune]*DFAState // DFA transitions
}

func NewDFAState() *DFAState {
	return &DFAState{
		nfaStates:   make(map[*State]bool),
		transitions: make(map[rune]*DFAState),
		isAccepting: false,
	}
}

// Adds epsilon closure to nfa states of DFAState
func (d *DFAState) epsilonClosureDFA() {
	keys := make([]*State, 0)
	for k := range d.nfaStates {
		keys = append(keys, k)
	}
	keys = EpsilonClosure(keys)
	newSet := make(map[*State]bool)
	for _, k := range keys {
		if k.isAccepting {
			d.isAccepting = true
		}
		newSet[k] = true
	}
	d.nfaStates = newSet
}

// DFAResult holds the resulting DFA
type DFAResult struct {
	startState *DFAState
	states     []*DFAState
}

// Adds s to states of DFA state if not present already.
// Also turns DFA state into accepting if s is accepting
func checkAndAdd(s *State, dfaState *DFAState) *DFAState {
	if dfaState == nil {
		dfaState = NewDFAState()
	}
	if !dfaState.nfaStates[s] {
		dfaState.nfaStates[s] = true
		if s.isAccepting {
			dfaState.isAccepting = true
		}
	}
	return dfaState
}

// Get all non-epsilon transitions from the given set of states
func getTransitions(dfaState *DFAState) map[rune]*DFAState {
	out := make(map[rune]*DFAState)
	for state := range dfaState.nfaStates {
		for _, transition := range state.transitions {
			switch tr := transition.(type) {
			case *AlphaTransition:
				out[tr.char] = checkAndAdd(tr.s, out[tr.char])
			case *RangeTransition:
				for i := tr.from; i <= tr.to; i++ {
					out[i] = checkAndAdd(tr.s, out[i])
				}
			}
		}
	}
	return out
}

// Checks if newSet is already present in sets
func containsSet(sets []map[*State]bool, newSet map[*State]bool) bool {
	for _, existingSet := range sets {
		if len(existingSet) != len(newSet) {
			continue
		}

		allMatch := true
		for key, newValue := range newSet {
			existingValue, exists := existingSet[key]
			if !exists || existingValue != newValue {
				allMatch = false
				break
			}
		}
		if allMatch {
			return true
		}
	}
	return false
}

// Removes unknown states
func cleanupDFAResult(r *DFAResult) *DFAResult {
	existingStates := make(map[*DFAState]bool)
	for _, s := range r.states {
		existingStates[s] = true
	}
	for i, s := range r.states {
		newtr := make(map[rune]*DFAState)
		for r, tr := range s.transitions {
			if existingStates[tr] {
				newtr[r] = tr
			}
		}
		r.states[i].transitions = newtr
	}
	return r
}

func GenerateDFA(nfa *Automata) *DFAResult {
	stateQueue := util.NewQueue[*DFAState]()

	seenSets := make([]map[*State]bool, 0)
	start := NewDFAState()
	states := make([]*DFAState, 0)

	start.nfaStates[nfa.head] = true
	start.epsilonClosureDFA()
	stateQueue.Push(start)
	for !stateQueue.IsEmpty() {
		state, _ := stateQueue.Pop()
		if containsSet(seenSets, state.nfaStates) {
			continue
		}

		seenSets = append(seenSets, state.nfaStates)
		states = append(states, state)
		transitions := getTransitions(state)
		state.transitions = make(map[rune]*DFAState)
		for r, s := range transitions {
			s.epsilonClosureDFA()
			state.transitions[r] = s
			stateQueue.Push(s)
		}
	}

	return cleanupDFAResult(&DFAResult{
		startState: start,
		states:     states,
	})
}

// PrintDFA prints the DFA structure for debugging
func (dfa *DFAResult) PrintDFA() {
	fmt.Println("DFA States:")
	for i, state := range dfa.states {
		fmt.Printf("State %d (Accepting: %t):\n", i, state.isAccepting)
		fmt.Print("  NFA States: ")
		for s := range state.nfaStates {
			fmt.Printf("[%p] ", s)
		}
		fmt.Println()

		fmt.Println("  Transitions:")
		for char, target := range state.transitions {
			targetIndex := -1
			for i, s := range dfa.states {
				if s == target {
					targetIndex = i
					break
				}
			}
			fmt.Printf("    '%c' -> State %d\n", char, targetIndex)
		}
	}
}
