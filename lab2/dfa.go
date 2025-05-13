package main

import (
	"errors"
	"fmt"
	"strings"

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
	newDFA := make(map[*State]bool)
	for _, k := range keys {
		if k.isAccepting {
			d.isAccepting = true
		}
		newDFA[k] = true
	}
	d.nfaStates = newDFA
}

// Transits by given character
func (d *DFAState) transit(r rune) *DFAState {
	s, ok := d.transitions[r]
	if !ok {
		return nil
	}
	return s
}

// DFAResult holds the resulting DFA
type DFAResult struct {
	startState *DFAState
	states     []*DFAState
}

func (d *DFAResult) Match(matchee string) string {
	var maxOut, curString strings.Builder
	currentState := d.startState
	for _, symbol := range matchee {
		newState := currentState.transit(symbol)
		if newState == nil {
			break
		}
		curString.WriteRune(symbol)
		if newState.isAccepting {
			maxOut = curString
		}
		currentState = newState
	}
	return maxOut.String()
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

// Checks equality of 2 given sets
func setEquals(left map[*State]bool, right map[*State]bool) bool {
	if len(left) != len(right) {
		return false
	}
	for lKey, lVal := range left {
		rVal, exists := right[lKey]
		if !exists || rVal != lVal {
			return false
		}
	}
	return true
}

// Checks if DFA with the same set of nfa states is already present.
// If present, returns pointer to said DFA
func getDFABySet(dfas []*DFAState, newDFA *DFAState) *DFAState {
	for _, dfa := range dfas {
		if setEquals(dfa.nfaStates, newDFA.nfaStates) {
			return dfa
		}
	}
	return nil
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
	startState := NewDFAState()
	startState.nfaStates[nfa.head] = true
	startState.epsilonClosureDFA()

	stateQueue := util.NewQueue[*DFAState]()
	stateQueue.Push(startState)

	states := make([]*DFAState, 0)

	for !stateQueue.IsEmpty() {
		state, _ := stateQueue.Pop()

		states = append(states, state)
		transitions := getTransitions(state)
		state.transitions = make(map[rune]*DFAState)
		for r, s := range transitions {
			s.epsilonClosureDFA()
			existing := getDFABySet(states, s)
			if existing != nil {
				state.transitions[r] = existing
				continue
			}
			state.transitions[r] = s
			stateQueue.Push(s)
		}
	}

	return cleanupDFAResult(&DFAResult{
		startState: startState,
		states:     states,
	})
}

// Returns 'local' alphabet of the given DFA - all symbols that
// act as a transition in at least one state
func getAlphabet(dfa *DFAResult) map[rune]bool {
	alphabet := make(map[rune]bool)
	for _, state := range dfa.states {
		for r := range state.transitions {
			alphabet[r] = true
		}
	}
	return alphabet
}

// Alias for better readability
type DFASet map[*DFAState]bool

// Replaces DFAs in allSets with those specified in replaceMap
func replaceSets(allSets []*DFASet, replaceMap map[*DFASet][]*DFASet) []*DFASet {
	result := make([]*DFASet, 0)
	for _, elem := range allSets {
		if replacements, ok := replaceMap[elem]; ok {
			result = append(result, replacements...)
		} else {
			result = append(result, elem)
		}
	}
	return result
}

// Minimizes given DFA
func (dfa *DFAResult) Minimize() error {
	acceptingGroup := make(DFASet)
	otherGroup := make(DFASet)
	for _, s := range dfa.states {
		if s.isAccepting {
			acceptingGroup[s] = true
		} else {
			otherGroup[s] = true
		}
	}
	groupsList := make([]*DFASet, 0)
	groupsQueue := util.NewQueue[*DFASet]()

	if len(otherGroup) != 0 {
		groupsList = append(groupsList, &otherGroup)
		groupsQueue.Push(&otherGroup)
	}
	if len(acceptingGroup) != 0 {
		groupsList = append(groupsList, &acceptingGroup)
		groupsQueue.Push(&acceptingGroup)
	}

	alphabet := getAlphabet(dfa)
	// Forming groups
	for !groupsQueue.IsEmpty() {
		targetGroup, _ := groupsQueue.Pop()
		for symbol := range alphabet {

			// Map that stores which state should be replaced with which states
			replaceMap := make(map[*DFASet][]*DFASet)

			for _, group := range groupsList {
				inTarget := make(DFASet)
				other := make(DFASet)

				for state := range *group {
					tr := state.transit(symbol)
					if tr == nil {
						// continue
						other[state] = true
					}

					if (*targetGroup)[tr] {
						inTarget[state] = true
					} else {
						other[state] = true
					}
				}
				if len(inTarget) != 0 && len(other) != 0 {
					// Group can be divided, so we will do this
					replaceMap[group] = []*DFASet{&inTarget, &other}
					groupsQueue.Push(&inTarget)
					groupsQueue.Push(&other)
				}
			}

			groupsList = replaceSets(groupsList, replaceMap)
		}
	}

	// Creating new states
	newStates := make([]*DFAState, 0)
	// Which old state now belongs to which new state
	oldToNew := make(map[*DFAState]*DFAState)

	for _, group := range groupsList {
		newState := NewDFAState()
		for state := range *group {
			oldToNew[state] = newState
			if state.isAccepting {
				newState.isAccepting = true
			}
			if state == dfa.startState {
				dfa.startState = newState
			}
		}
		newStates = append(newStates, newState)
	}

	// Making connections

	// Since newState[i] was made from groupsList[i], we can do this
	for i := 0; i < len(groupsList); i++ {
		curGroup := groupsList[i]
		curNewState := newStates[i]
		for state := range *curGroup {
			for r, tr := range state.transitions {
				newState, ok := oldToNew[tr]
				if !ok {
					return errors.New("found old state that doesn't correspond to the new one")
				}
				curNewState.transitions[r] = newState
			}
		}
	}
	dfa.states = newStates
	return nil
}

// PrintDFA prints the DFA structure for debugging
func (dfa *DFAResult) PrintDFA() {
	fmt.Println("DFA States:")
	for i, state := range dfa.states {
		fmt.Printf("State %d (Accepting: %t, Start: %t):\n", i, state.isAccepting, dfa.startState == state)
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
