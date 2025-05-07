package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueue(t *testing.T) {
	q := NewQueue[int]()
	q.Push(10)
	out, _ := q.Pop()
	assert.Equal(t, out, 10)
	_, ok := q.Pop()
	assert.False(t, ok)
}
