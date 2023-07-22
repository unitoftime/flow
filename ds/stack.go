package ds

type Stack[T any] struct {
	Buffer []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		Buffer: make([]T, 0), // TODO: configurable?
	}
}

func (s *Stack[T]) Len() int {
	return len(s.Buffer)
}

func (s *Stack[T]) Add(t T) {
	s.Buffer = append(s.Buffer, t)
}

// func (s *Stack[T]) Peek() (T, bool) {
// 	if q.ReadIdx == q.WriteIdx {
// 		var ret T
// 		return ret, false
// 	}
// 	return q.Buffer[q.ReadIdx], true
// }
// func (s *Stack[T]) PeekLast() (T, bool) {
// 	if q.ReadIdx == q.WriteIdx {
// 		var ret T
// 		return ret, false
// 	}
// 	idx := (q.WriteIdx + len(q.Buffer) - 1) % len(q.Buffer)
// 	return q.Buffer[idx], true
// }
func (s *Stack[T]) Remove() (T, bool) {
	if len(s.Buffer) == 0 {
		var ret T
		return ret, false
	}

	last := len(s.Buffer)-1
	val := s.Buffer[last]
	s.Buffer = s.Buffer[:last]

	return val, true
}
