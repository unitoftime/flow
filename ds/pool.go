package ds

import (
	"github.com/zergon321/mempool"
)

type erasable[T any] struct {
	val T
}

func newErasable[T any](v T) erasable[T] {
	return erasable[T]{v}
}

func (s erasable[T]) Erase() error {
	return nil
}

type SlicePool[T any] struct {
	inner *mempool.Pool[erasable[[]T]]
}

func NewSlicePool[T any](defaultSliceSize int) SlicePool[T] {
	inner, err := mempool.NewPool(func() erasable[[]T] {
		return newErasable(make([]T, 0, defaultSliceSize))
	})
	if err != nil {
		panic(err) // TODO: The only failure case should be caused by options
	}

	return SlicePool[T]{
		inner: inner,
	}
}

func (p SlicePool[T]) Put(slice []T) {
	slice = slice[:0] // Erase

	err := p.inner.Put(newErasable(slice))
	if err != nil {
		panic(err)
	}
}

func (p SlicePool[T]) Get() []T {
	s := p.inner.Get()
	return s.val
}

type MapPool[K comparable, V any] struct {
	inner *mempool.Pool[erasable[map[K]V]]
}

func NewMapPool[K comparable, V any](defaultSize int) MapPool[K, V] {
	inner, err := mempool.NewPool(func() erasable[map[K]V] {
		return newErasable(make(map[K]V, defaultSize))
	})

	if err != nil {
		panic(err)
	}

	return MapPool[K, V]{
		inner: inner,
	}
}

func (p MapPool[K, V]) Put(m map[K]V) {
	// Erase
	for k := range m {
		delete(m, k)
	}

	err := p.inner.Put(newErasable(m))
	if err != nil {
		panic(err)
	}
}

func (p MapPool[K, V]) Get() map[K]V {
	s := p.inner.Get()
	return s.val
}

func (p MapPool[K, V]) Clone(og map[K]V) map[K]V {
	m := p.Get()
	for k, v := range og {
		m[k] = v
	}
	return m
}
