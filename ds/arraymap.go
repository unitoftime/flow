package ds

import (
	"iter"
)

type arrayMapData[K comparable, V any] struct {
	key   K
	value V
}

func newArrayMapData[K comparable, V any](k K, v V) arrayMapData[K, V] {
	return arrayMapData[K, V]{
		key:   k,
		value: v,
	}
}

// Acts like a map, but is backed by an array. Can provide better iteration speed at the cost of slower lookups
type ArrayMap[K comparable, V any] struct {
	slice []arrayMapData[K, V] // TODO: Maybe faster with two slices? one for key, another for value?
}

func NewArrayMap[K comparable, V any]() ArrayMap[K, V] {
	return ArrayMap[K, V]{
		slice: make([]arrayMapData[K, V], 0),
	}
}

func (m *ArrayMap[K, V]) append(key K, val V) {
	m.slice = append(m.slice, newArrayMapData(key, val))
}

// returns the index of the value, or -1 if the key does not exist
func (m *ArrayMap[K, V]) find(key K) int {
	for i := range m.slice {
		if m.slice[i].key != key {
			continue // Skip: Wrong key
		}

		return i
	}

	return -1
}

func (m *ArrayMap[K, V]) Put(key K, val V) {
	idx := m.find(key)
	if idx < 0 {
		// Can't find key, so just append
		m.append(key, val)
	} else {
		// Else, just update the current key
		m.slice[idx].value = val
	}
}

func (m *ArrayMap[K, V]) Get(key K) (V, bool) {
	idx := m.find(key)
	if idx < 0 {
		var v V
		return v, false
	} else {
		return m.slice[idx].value, true
	}
}

// Delete a specific index. Note this will move the last index into the hole
func (m *ArrayMap[K, V]) Delete(key K) {
	idx := m.find(key)
	if idx < 0 {
		return // Nothing to do, does not exist
	}

	m.slice[idx] = m.slice[len(m.slice)-1]
	m.slice = m.slice[:len(m.slice)-1]
}

// Clear the entire slice
func (m *ArrayMap[K, V]) Clear() {
	m.slice = m.slice[:0]
}

// Iterate through the entire slice
func (m *ArrayMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, d := range m.slice {
			if !yield(d.key, d.value) {
				break
			}
		}
	}
}
