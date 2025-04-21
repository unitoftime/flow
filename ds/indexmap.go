package ds

import (
	"iter"

	"golang.org/x/exp/constraints"
)

// Acts like a map, but backed by an integer indexed slice instead of a map
// This is mostly for use cases where you have a small list of increasing numbers that you want to store in an array, but you dont want to worry about ensuring the bounds are always correct
// Note: If you put a huge key in here, the slice will allocate a ton of space.
type IndexMap[K constraints.Integer, V any] struct {
	// TODO: set slice should just be a bitmask
	set   []bool // Tracks whether or not the data at a location is set or empty
	slice []V    // Tracks the data
}

func NewIndexMap[K constraints.Integer, V any]() IndexMap[K, V] {
	return IndexMap[K, V]{
		set:   make([]bool, 0),
		slice: make([]V, 0),
	}
}

func (m *IndexMap[K, V]) grow(idx K) {
	requiredLength := idx + 1
	growAmount := requiredLength - K(len(m.set))
	if growAmount <= 0 {
		return // No need to grow if the sliceIdx is already in bounds
	}

	m.set = append(m.set, make([]bool, growAmount)...)
	m.slice = append(m.slice, make([]V, growAmount)...)
}

func (m *IndexMap[K, V]) Put(idx K, val V) {
	if idx < 0 {
		return
	}

	m.grow(idx) // Ensure index is within bounds

	m.set[idx] = true
	m.slice[idx] = val
}

func (m *IndexMap[K, V]) Get(idx K) (V, bool) {
	if idx < 0 || idx >= K(len(m.set)) {
		var v V
		return v, false
	}

	return m.slice[idx], m.set[idx]
}

// Delete a specific index
func (m *IndexMap[K, V]) Delete(idx K) {
	m.set[idx] = false
}

// Clear the entire slice
func (m *IndexMap[K, V]) Clear() {
	m.slice = m.slice[:0]
}

// Iterate through the entire slice
func (m *IndexMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for i, v := range m.slice {
			// Ensure that the map key is set
			if !m.set[i] {
				continue
			}

			if !yield(K(i), v) {
				break
			}
		}
	}
}
