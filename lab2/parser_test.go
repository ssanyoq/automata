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

// func TestParseBrackets(t *testing.T) {
// 	tests := []struct {
// 		name              string
// 		inputString       string
// 		expectedError     bool
// 		expectedFragments []Automata
// 	}{
// 		{
// 			name:          "basic example",
// 			inputString:   "[a-z]",
// 			expectedError: false,
// 			expectedFragments: []Automata{&Automata{
// 				From: 'a',
// 				To:   'z',
// 			}},
// 		},
// 		{
// 			name:              "bad range",
// 			inputString:       "[z-a]",
// 			expectedError:     true,
// 			expectedFragments: []Automata{},
// 		},
// 		{
// 			name:              "not enough arguments",
// 			inputString:       "[a-]",
// 			expectedError:     true,
// 			expectedFragments: []Automata{},
// 		},
// 		{
// 			name:              "not enough arguments2",
// 			inputString:       "[-z]",
// 			expectedError:     true,
// 			expectedFragments: []Automata{},
// 		},
// 		{
// 			name:              "not enough arguments3",
// 			inputString:       "[az]",
// 			expectedError:     true,
// 			expectedFragments: []Automata{},
// 		},
// 	}
// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			input := test.inputString[1:] // to cut '['
// 			parser := NewParser(NewLexer(input))

// 			parser.parseBrackets()

// 			if test.expectedError {
// 				assert.NotEmpty(t, parser.Errors)
// 				return
// 			} else {
// 				assert.Empty(t, parser.Errors)
// 			}
// 			assert.Equal(t, test.expectedFragments, unloadStack(parser.fragmentsStack))
// 		})
// 	}
// }

// func TestBuildAST(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		inputString string
// 		expected    Automata
// 		expectErr   bool
// 	}{
// 		{
// 			name:        "simple",
// 			inputString: "abc",
// 			expected: &BinaryOpAutomata{
// 				Operation: Concat,
// 				Left:      &CharAutomata{Character: 'a'},
// 				Right: &BinaryOpAutomata{
// 					Operation: Concat,
// 					Left:      &CharAutomata{Character: 'b'},
// 					Right:     &CharAutomata{Character: 'c'},
// 				},
// 			},
// 			expectErr: false,
// 		},
// 		{
// 			name:        "simple2",
// 			inputString: "a|bc",
// 			expected: &BinaryOpAutomata{
// 				Operation: Or,
// 				Right: &BinaryOpAutomata{
// 					Operation: Concat,
// 					Left:      &CharAutomata{Character: 'b'},
// 					Right:     &CharAutomata{Character: 'c'},
// 				},
// 				Left: &CharAutomata{Character: 'a'},
// 			},
// 			expectErr: false,
// 		},
// 		{
// 			name:        "harder",
// 			inputString: "[a-z]+",
// 			expected: &UnaryOpAutomata{
// 				Operation: PositiveClosure,
// 				Child: &CharacterRangeAutomata{
// 					From: 'a',
// 					To:   'z',
// 				},
// 			},
// 		},
// 		{
// 			name:        "email",
// 			inputString: "[a-z]+@ya(ndex){,1}%.ru",
// 			expected: &BinaryOpAutomata{
// 				Operation: Concat,
// 				Left: &UnaryOpAutomata{
// 					Operation: PositiveClosure,
// 					Child: &CharacterRangeAutomata{
// 						From: 'a',
// 						To:   'z',
// 					},
// 				},
// 				Right: &BinaryOpAutomata{
// 					Operation: Concat,
// 					Left: &CharAutomata{
// 						Character: '@',
// 					},
// 					Right: &BinaryOpAutomata{
// 						Operation: Concat,
// 						Left: &CharAutomata{
// 							Character: 'y',
// 						},
// 						Right: &BinaryOpAutomata{
// 							Operation: Concat,
// 							Left:      &CharAutomata{Character: 'a'},
// 							Right: &BinaryOpAutomata{
// 								Operation: Concat,
// 								Left: &RangeRepeatAutomata{
// 									From: 0,
// 									To:   1,
// 									Child: {
// 										Number: -1,
// 										Child: &BinaryOpAutomata{
// 											Operation: Concat,
// 											Left:      &CharAutomata{Character: 'n'},
// 											Right: &BinaryOpAutomata{
// 												Operation: Concat,
// 												Left:      &CharAutomata{Character: 'd'},
// 												Right: &BinaryOpAutomata{
// 													Operation: Concat,
// 													Left:      &CharAutomata{Character: 'e'},
// 													Right:     &CharAutomata{Character: 'x'},
// 												},
// 											},
// 										},
// 									},
// 								},
// 								Right: &BinaryOpAutomata{
// 									Operation: Concat,
// 									Left:      &CharAutomata{Character: '.'},
// 									Right: &BinaryOpAutomata{
// 										Operation: Concat,
// 										Left:      &CharAutomata{Character: 'r'},
// 										Right:     &CharAutomata{Character: 'u'},
// 									},
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}
// 	for _, test := range tests {
// 		if test.name != "email" {
// 			continue
// 		}
// 		t.Run(test.name, func(t *testing.T) {
// 			input := test.inputString
// 			parser := NewParser(NewLexer(input))
// 			ast, err := parser.BuildNFA()
// 			if test.expectErr {
// 				require.Error(t, err)
// 				return
// 			}
// 			require.NoError(t, err)
// 			assert.Equal(t, test.expected, ast)
// 		})
// 	}
// }
