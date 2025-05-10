package util

type Queue[T any] struct {
	Items []T
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{}
}

func (s *Queue[T]) Push(item T) {
	s.Items = append(s.Items, item)
}

func (s *Queue[T]) Pop() (T, bool) {
	if len(s.Items) == 0 {
		var zeroValue T
		return zeroValue, false
	}
	value := s.Items[0]
	s.Items = s.Items[1:]
	return value, true
}

func (s *Queue[T]) Peek() (T, bool) {
	if len(s.Items) == 0 {
		var zeroValue T
		return zeroValue, false
	}
	return s.Items[0], true
}

func (s *Queue[T]) IsEmpty() bool {
	return len(s.Items) == 0
}

func (s *Queue[T]) Size() int {
	return len(s.Items)
}
