package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestParseBrackets(t *testing.T) {
	tests := []struct {
		name              string
		inputString       string
		expectedError     bool
		expectedFragments []Node
	}{
		{
			name:          "basic example",
			inputString:   "[a-z]",
			expectedError: false,
			expectedFragments: []Node{&CharacterRangeNode{
				From: 'a',
				To:   'z',
			}},
		},
		{
			name:              "bad range",
			inputString:       "[z-a]",
			expectedError:     true,
			expectedFragments: []Node{},
		},
		{
			name:              "not enough arguments",
			inputString:       "[a-]",
			expectedError:     true,
			expectedFragments: []Node{},
		},
		{
			name:              "not enough arguments2",
			inputString:       "[-z]",
			expectedError:     true,
			expectedFragments: []Node{},
		},
		{
			name:              "not enough arguments3",
			inputString:       "[az]",
			expectedError:     true,
			expectedFragments: []Node{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := test.inputString[1:] // to cut '['
			parser := NewParser(NewLexer(input))

			parser.parseBrackets()

			if test.expectedError {
				assert.NotEmpty(t, parser.Errors)
				return
			} else {
				assert.Empty(t, parser.Errors)
			}
			assert.Equal(t, test.expectedFragments, unloadStack(parser.fragmentsStack))
		})
	}
}

func TestBuildAST(t *testing.T) {
	tests := []struct {
		name        string
		inputString string
		expected    Node
		expectErr   bool
	}{
		{
			name:        "simple",
			inputString: "abc",
			expected: &BinaryOpNode{
				Operation: Concat,
				Left:      &CharNode{Character: 'a'},
				Right: &BinaryOpNode{
					Operation: Concat,
					Left:      &CharNode{Character: 'b'},
					Right:     &CharNode{Character: 'c'},
				},
			},
			expectErr: false,
		},
		{
			name:        "simple2",
			inputString: "a|bc",
			expected: &BinaryOpNode{
				Operation: Or,
				Right: &BinaryOpNode{
					Operation: Concat,
					Left:      &CharNode{Character: 'b'},
					Right:     &CharNode{Character: 'c'},
				},
				Left: &CharNode{Character: 'a'},
			},
			expectErr: false,
		},
		{
			name:        "harder",
			inputString: "[a-z]+",
			expected: &UnaryOpNode{
				Operation: PositiveClosure,
				Child: &CharacterRangeNode{
					From: 'a',
					To:   'z',
				},
			},
		},
		{
			name:        "email",
			inputString: "[a-z]+@ya(ndex){,1}%.ru",
			expected: &BinaryOpNode{
				Operation: Concat,
				Left: &UnaryOpNode{
					Operation: PositiveClosure,
					Child: &CharacterRangeNode{
						From: 'a',
						To:   'z',
					},
				},
				Right: &BinaryOpNode{
					Operation: Concat,
					Left: &CharNode{
						Character: '@',
					},
					Right: &BinaryOpNode{
						Operation: Concat,
						Left: &CharNode{
							Character: 'y',
						},
						Right: &BinaryOpNode{
							Operation: Concat,
							Left:      &CharNode{Character: 'a'},
							Right: &BinaryOpNode{
								Operation: Concat,
								Left: &RangeRepeatNode{
									From: 0,
									To:   1,
									Child: &CaptureGroupNode{
										Number: -1,
										Child: &BinaryOpNode{
											Operation: Concat,
											Left:      &CharNode{Character: 'n'},
											Right: &BinaryOpNode{
												Operation: Concat,
												Left:      &CharNode{Character: 'd'},
												Right: &BinaryOpNode{
													Operation: Concat,
													Left:      &CharNode{Character: 'e'},
													Right:     &CharNode{Character: 'x'},
												},
											},
										},
									},
								},
								Right: &BinaryOpNode{
									Operation: Concat,
									Left:      &CharNode{Character: '.'},
									Right: &BinaryOpNode{
										Operation: Concat,
										Left:      &CharNode{Character: 'r'},
										Right:     &CharNode{Character: 'u'},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, test := range tests {
		if test.name != "email" {
			continue
		}
		t.Run(test.name, func(t *testing.T) {
			input := test.inputString
			parser := NewParser(NewLexer(input))
			ast, err := parser.BuildAST()
			if test.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expected, ast)
		})
	}
}
