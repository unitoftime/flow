package ds

import (
	"iter"
)

type ArrayConstraint[T any] interface {
	[1]T | [2]T | [3]T | [4]T | [5]T | [6]T | [7]T | [8]T | [9]T | [10]T | [11]T | [12]T | [13]T | [14]T | [15]T | [16]T
}

// A mini slice that allocates the first N elements into an array and then heap allocates the remaining into a traditional slice
type MiniSlice[A ArrayConstraint[T], T comparable] struct {
	Array        A
	Slice        []T
	nextArrayIdx uint8
}

func (s *MiniSlice[A, T]) getIdx(idx int) (innerIdx int, array bool) {
	arrayLen := len(s.Array)
	if idx >= arrayLen {
		innerIdx := idx - arrayLen
		return innerIdx, false
	}

	return idx, true
}

func (s *MiniSlice[A, T]) Get(idx int) T {
	innerIdx, isArray := s.getIdx(idx)
	if isArray {
		return s.Array[innerIdx]
	} else {
		return s.Slice[innerIdx]
	}
}

func (s *MiniSlice[A, T]) Set(idx int, val T) {
	innerIdx, isArray := s.getIdx(idx)
	if isArray {
		s.Array[innerIdx] = val
	} else {
		s.Slice[innerIdx] = val
	}
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

// Find and return the index of the first element, else return -1
func (s *MiniSlice[A, T]) Find(searchVal T) int {
	for i, val := range s.All() {
		if searchVal == val {
			return i
		}
	}
	return -1
}

// Removes the element at the supplied index, swapping the element that was at the last index to
// the supplied index
func (s *MiniSlice[A, T]) Delete(idx int) {
	if idx < 0 {
		return
	}
	if idx > s.Len() {
		return
	}

	lastVal := s.Get(s.Len() - 1)
	s.Set(idx, lastVal)
	s.SliceLast()
}

// Slices the last element
func (s *MiniSlice[A, T]) SliceLast() {
	innerIdx, isArray := s.getIdx(s.Len() - 1)
	if isArray {
		s.nextArrayIdx--
	} else {
		s.Slice = s.Slice[:innerIdx]
	}
}

// // Returns the last index, or returns -1 if empty
// func (s *MiniSlice[A, T]) Last() int {
// 	// If last element is on slice
// 	if len(s.Slice) > 0 {
// 		return len(s.Slice)-1
// 		// val := s.Slice[lastIdx]
// 		// s.Slice = s.Slice[:lastIdx]
// 		// return val
// 	}

// 	// Else last element is on array
// 	return int(s.nextArrayIdx) - 1
// }

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
			if !yield(i+arrayEnd, s.Slice[i]) {
				return
			}
		}
	}
}
