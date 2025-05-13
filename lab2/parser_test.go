package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ssanyoq/automata-uni/lab2/util"
)

func newLoadedStack[T any](load []T) *util.Stack[T] {
	stack := util.NewStack[T]()
	stack.Items = load
	return stack
}

func unloadStack[T any](s *util.Stack[T]) []T {
	return s.Items
}

func TestPopTillStop(t *testing.T) {

	tests := []struct {
		name                string
		operatorsStackItems []Token
		fragmentsStackItems []*Automata
		expectedFragments   []*Automata
	}{
		{
			name:                "single binary operator",
			operatorsStackItems: []Token{Concat},
			fragmentsStackItems: []*Automata{CharAutomata('a'), CharAutomata('b')},
			expectedFragments: []*Automata{
				ConcatAutomata(CharAutomata('a'), CharAutomata('b')),
			},
		},
		{
			name:                "positive closure",
			operatorsStackItems: []Token{PositiveClosure},
			fragmentsStackItems: []*Automata{CharAutomata('a')},
			expectedFragments: []*Automata{
				ConcatAutomata(CharAutomata('a'), KleeneeAutomata(CharAutomata('a'))),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := &Parser{
				fragmentsStack: newLoadedStack(test.fragmentsStackItems),
				operatorsStack: newLoadedStack(test.operatorsStackItems),
				Errors:         []error{},
			}

			// Run the function we are testing
			parser.popTillStop()

			assert.Empty(t, parser.Errors)
			assert.Equal(t, test.expectedFragments, unloadStack(parser.fragmentsStack))
		})
	}
}
