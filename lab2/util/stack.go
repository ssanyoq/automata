package util

type Stack[T any] struct {
	Items []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

func (s *Stack[T]) Push(item T) {
	s.Items = append(s.Items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.Items) == 0 {
		var zeroValue T
		return zeroValue, false
	}
	lastIndex := len(s.Items) - 1
	value := s.Items[lastIndex]
	s.Items = s.Items[:lastIndex]
	return value, true
}

func (s *Stack[T]) Peek() (T, bool) {
	if len(s.Items) == 0 {
		var zeroValue T
		return zeroValue, false
	}
	return s.Items[len(s.Items)-1], true
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.Items) == 0
}

func (s *Stack[T]) Size() int {
	return len(s.Items)
}
