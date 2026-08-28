package ds

import (
	"testing"
)

func TestStack(t *testing.T) {
	s := NewStack[int]()

	compare(t, s.Len(), 0)
	
	val, ok := s.Peek()
	compare(t, ok, false)

	val, ok = s.Remove()
	compare(t, ok, false)

	s.Add(10)
	compare(t, s.Len(), 1)
	
	val, ok = s.Peek()
	compare(t, ok, true)
	compare(t, val, 10)

	s.Add(20)
	compare(t, s.Len(), 2)
	
	val, ok = s.Peek()
	compare(t, ok, true)
	compare(t, val, 20)

	val, ok = s.Remove()
	compare(t, ok, true)
	compare(t, val, 20)
	compare(t, s.Len(), 1)

	val, ok = s.Peek()
	compare(t, ok, true)
	compare(t, val, 10)

	s.Add(30)
	s.Clear()
	compare(t, s.Len(), 0)
}
