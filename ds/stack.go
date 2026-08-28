package ds

type Stack[T any] struct {
	Buffer []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		Buffer: make([]T, 0),
	}
}

func (s *Stack[T]) Len() int {
	return len(s.Buffer)
}

func (s *Stack[T]) Add(t T) {
	s.Buffer = append(s.Buffer, t)
}

func (s *Stack[T]) Peek() (T, bool) {
	if len(s.Buffer) == 0 {
		var ret T
		return ret, false
	}
	return s.Buffer[len(s.Buffer)-1], true
}

func (s *Stack[T]) Remove() (T, bool) {
	if len(s.Buffer) == 0 {
		var ret T
		return ret, false
	}
	idx := len(s.Buffer) - 1
	ret := s.Buffer[idx]
	var zero T
	s.Buffer[idx] = zero
	s.Buffer = s.Buffer[:idx]
	return ret, true
}

func (s *Stack[T]) Clear() {
	s.Buffer = s.Buffer[:0]
}
