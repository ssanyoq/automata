package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newLoadedStack[T any](load []T) *Stack[T] {
	stack := NewStack[T]()
	stack.items = load
	return stack
}

func unloadStack[T any](s *Stack[T]) []T {
	return s.items
}

func TestPopTillStop(t *testing.T) {

	tests := []struct {
		name                string
		operatorsStackItems []Token
		fragmentsStackItems []Node
		expectedFragments   []Node
	}{
		{
			name:                "single binary operator",
			operatorsStackItems: []Token{Concat},
			fragmentsStackItems: []Node{&CharNode{Character: 'a'}, &CharNode{Character: 'b'}},
			expectedFragments: []Node{&BinaryOpNode{
				Left:      &CharNode{Character: 'a'},
				Right:     &CharNode{Character: 'b'},
				Operation: Concat,
			}},
		},
		{
			name:                "parentheses",
			operatorsStackItems: []Token{Concat, OpenParenthesis, PositiveClosure},
			fragmentsStackItems: []Node{&CharNode{Character: 'b'}, &CharNode{Character: 'a'}},
			expectedFragments: []Node{
				&CharNode{Character: 'b'},
				&CaptureGroupNode{Number: -1,
					Child: &UnaryOpNode{
						Operation: PositiveClosure,
						Child:     &CharNode{Character: 'a'},
					}},
			},
		},
		{
			name:                "parentheses with others",
			operatorsStackItems: []Token{Prognostic},
			fragmentsStackItems: []Node{
				&CaptureGroupNode{Number: 1,
					Child: &UnaryOpNode{
						Operation: PositiveClosure,
						Child:     &CharNode{Character: 'a'},
					}},
				&CaptureGroupNode{Number: 2,
					Child: &UnaryOpNode{
						Operation: Kleene,
						Child:     &CharNode{Character: 'b'},
					}},
			},
			expectedFragments: []Node{
				&BinaryOpNode{
					Operation: Prognostic,
					Left: &CaptureGroupNode{Number: 1,
						Child: &UnaryOpNode{
							Operation: PositiveClosure,
							Child:     &CharNode{Character: 'a'},
						}},
					Right: &CaptureGroupNode{Number: 2,
						Child: &UnaryOpNode{
							Operation: Kleene,
							Child:     &CharNode{Character: 'b'},
						}},
				},
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
