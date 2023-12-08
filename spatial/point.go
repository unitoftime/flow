package spatial

import (
	"github.com/unitoftime/flow/phy2"
)

type PointBucket[T comparable] struct {
	List []T
}

func NewPointBucket[T comparable]() *PointBucket[T] {
	return &PointBucket[T]{
		List: make([]T, 0),
	}
}
func (b *PointBucket[T]) Add(val T) {
	b.List = append(b.List, val)
}
func (b *PointBucket[T]) Clear() {
	b.List = b.List[:0]
}


// TODO: rename? ColliderMap?
type Pointmap[T comparable] struct {
	PositionHasher

	Bucket *arrayMap[PointBucket[T]]
	allBuckets []*PointBucket[T]
}

func NewPointmap[T comparable](chunksize [2]int) *Pointmap[T] {
	return &Pointmap[T]{
		allBuckets: make([]*PointBucket[T], 0, 1024),
		Bucket: newArrayMap[PointBucket[T]](),
		PositionHasher: NewPositionHasher(chunksize),
	}
}

func (h *Pointmap[T]) Clear() {
	for _, b := range h.allBuckets {
		b.Clear()
	}
}

func (h *Pointmap[T]) GetBucket(index Index) *PointBucket[T] {
	bucket, ok := h.Bucket.Get(index.X, index.Y)
	if !ok {
		bucket = NewPointBucket[T]()
		h.allBuckets = append(h.allBuckets, bucket)
		h.Bucket.Put(index.X, index.Y, bucket)
	}
	return bucket
}

func (h *Pointmap[T]) Add(pos phy2.Pos, val T) {
	idx := h.PositionToIndex(pos)
	bucket := h.GetBucket(idx)
	bucket.Add(val)
}

// Adds the collisions directly into your collision list. Items are deduplicated by nature of them only existing once in this Pointmap. (ie if you add multiple of the same thing, you might get multiple out)
func (h *Pointmap[T]) Check(list []T, bounds phy2.Rect) []T {
	min := h.PositionToIndex(phy2.Pos(bounds.Min))
	max := h.PositionToIndex(phy2.Pos(bounds.Max))

	// TODO: Might be nice if this spirals from inside to outside, that way its roughly sorted by distance?
	for x := min.X; x <= max.X; x++ {
		for y := min.Y; y <= max.Y; y++ {
			bucket := h.GetBucket(Index{x, y})
			for i := range bucket.List {
				list = append(list, bucket.List[i])
			}
		}
	}

	return list
}
