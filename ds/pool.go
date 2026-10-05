package ds

import (
	"sync"
)

// A pool of reusable values that is safe for concurrent use
// Note: Values travel through sync.Pool inside of recycled *T boxes, so that neither Get nor Put allocates
type pool[T any] struct {
	full  sync.Pool // Boxes that hold a pooled value
	empty sync.Pool // Boxes that are waiting to hold a pooled value
}

func newPool[T any](constructor func() T) *pool[T] {
	return &pool[T]{
		full: sync.Pool{
			New: func() any {
				box := new(T)
				*box = constructor()
				return box
			},
		},
		empty: sync.Pool{
			New: func() any {
				return new(T)
			},
		},
	}
}

func (p *pool[T]) Get() T {
	box := p.full.Get().(*T)
	val := *box

	// Clear the box so that it doesn't keep the value alive while it waits
	var zero T
	*box = zero
	p.empty.Put(box)

	return val
}

func (p *pool[T]) Put(val T) {
	box := p.empty.Get().(*T)
	*box = val
	p.full.Put(box)
}

type SlicePool[T any] struct {
	inner *pool[[]T]
}

func NewSlicePool[T any](defaultSliceSize int) SlicePool[T] {
	return SlicePool[T]{
		inner: newPool(func() []T {
			return make([]T, 0, defaultSliceSize)
		}),
	}
}

func (p SlicePool[T]) Put(slice []T) {
	slice = slice[:0] // Erase

	p.inner.Put(slice)
}

func (p SlicePool[T]) Get() []T {
	return p.inner.Get()
}

type MapPool[K comparable, V any] struct {
	inner *pool[map[K]V]
}

func NewMapPool[K comparable, V any](defaultSize int) MapPool[K, V] {
	return MapPool[K, V]{
		inner: newPool(func() map[K]V {
			return make(map[K]V, defaultSize)
		}),
	}
}

func (p MapPool[K, V]) Put(m map[K]V) {
	// Erase
	for k := range m {
		delete(m, k)
	}

	p.inner.Put(m)
}

func (p MapPool[K, V]) Get() map[K]V {
	return p.inner.Get()
}

func (p MapPool[K, V]) Clone(og map[K]V) map[K]V {
	m := p.Get()
	for k, v := range og {
		m[k] = v
	}
	return m
}
