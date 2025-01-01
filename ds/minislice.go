package ds

import (
	"iter"
)

type ArrayConstraint[T any] interface {
	[1]T | [2]T | [3]T | [4]T | [5]T | [6]T | [7]T | [8]T | [9]T | [10]T | [11]T | [12]T | [13]T | [14]T | [15]T | [16]T
}


// A mini slice that allocates the first N elements into an array and then heap allocates the remaining into a traditional slice
type MiniSlice[A ArrayConstraint[T], T any] struct {
	Array A
	Slice []T
	nextArrayIdx uint8
}

func (s *MiniSlice[A, T]) Append(val T) {
	arrayLen := len(s.Array)
	if int(s.nextArrayIdx) >= arrayLen {
		// If nextArrayIndex is outside to the fixed-sized array, then just start appending to the slice
		s.Slice = append(s.Slice, val)
		return
	}

	// Else, append to the fixed array and track the index
	s.Array[s.nextArrayIdx] = val
	s.nextArrayIdx++
}

// Returns the number of elements in the slice
func (s *MiniSlice[A, T]) Len() int {
	return int(s.nextArrayIdx) + len(s.Slice)
}

// Clears the array, so the next append will start at index 0
func (s *MiniSlice[A, T]) Clear() {
	s.nextArrayIdx = 0
	s.Slice = s.Slice[:0]
}

// Iterates the slice from 0 to the last element added
func (s *MiniSlice[A, T]) All() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		arrayEnd := int(s.nextArrayIdx)

		for i := 0; i < arrayEnd; i++ {
			if !yield(i, s.Array[i]) {
				return
			}
		}

		for i := range s.Slice {
			if !yield(i + arrayEnd, s.Slice[i]) {
				return
			}
		}
	}
}
