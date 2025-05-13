package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEpsilonClosure(t *testing.T) {
	t.Run("simple no cycles", func(t *testing.T) {
		/*
			          ε    a
					a -> b -> c
		*/
		c := &State{
			isAccepting: true,
		}
		b := &State{
			transitions: []Transition{
				&AlphaTransition{
					s:    c,
					char: 'a',
				},
			},
		}
		a := &State{
			transitions: []Transition{
				&EpsilonTransition{
					s: b,
				},
			},
		}
		out := EpsilonClosure([]*State{a})
		assert.Equal(t, []*State{a, b}, out)
	})

	t.Run("simple with cycles", func(t *testing.T) {
		/*
			          ε
					a -> b
					^	  \
				     \     |
					 ε\____/

		*/
		b := &State{
			isAccepting: true,
		}
		a := &State{
			transitions: []Transition{
				&EpsilonTransition{
					s: b,
				},
			},
		}
		b.transitions = []Transition{
			&EpsilonTransition{
				s: a,
			},
		}
		out := EpsilonClosure([]*State{a})
		assert.Equal(t, []*State{a, b}, out)
	})

	t.Run("simple with cycles and sets", func(t *testing.T) {
		/*
			          a
					a -> b
					^	   \
				     \     /
					 ε -c<-ε

		*/
		b := &State{
			isAccepting: true,
		}
		a := &State{
			transitions: []Transition{
				&AlphaTransition{
					char: 'a',
					s:    b,
				},
			},
		}
		c := &State{
			transitions: []Transition{
				&EpsilonTransition{
					s: a,
				},
			},
		}
		b.transitions = []Transition{
			&EpsilonTransition{
				s: c,
			},
		}
		out := EpsilonClosure([]*State{a, b})
		assert.Equal(t, []*State{a, b, c}, out)
	})
}
