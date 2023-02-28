package spatial

import (
	"github.com/unitoftime/flow/phy2"
)

type BucketItem[T comparable] struct {
	bounds phy2.Rect
	item T
}

type Bucket[T comparable] struct {
	List []BucketItem[T]
}

func NewBucket[T comparable]() *Bucket[T] {
	return &Bucket[T]{
		List: make([]BucketItem[T], 0),
	}
}
func (b *Bucket[T]) Add(bounds phy2.Rect, val T) {
	b.List = append(b.List, BucketItem[T]{
		bounds: bounds,
		item: val,
	})
}
func (b *Bucket[T]) Clear() {
	b.List = b.List[:0]
}
func (b *Bucket[T]) Check(bounds phy2.Rect) []BucketItem[T] {
	ret := make([]BucketItem[T], 0)
	for i := range b.List {
		if bounds.Intersects(b.List[i].bounds) {
			ret = append(ret, b.List[i])
		}
	}
	return ret
}


type Index struct {
	X, Y int
}

type Hashmap[T comparable] struct {
	Bucket map[Index]*Bucket[T]
	chunksize [2]int
}

func NewHashmap[T comparable](chunksize [2]int) *Hashmap[T] {
	return &Hashmap[T]{
		Bucket: make(map[Index]*Bucket[T]),
		chunksize: chunksize,
	}
}

func (h *Hashmap[T]) Clear() {
	for _, b := range h.Bucket {
		b.Clear()
	}
}

func (h *Hashmap[T]) GetBucket(index Index) *Bucket[T] {
	bucket, ok := h.Bucket[index]
	if !ok {
		bucket = NewBucket[T]()
		h.Bucket[index] = bucket
	}
	return bucket
}

func (h *Hashmap[T]) Add(bounds phy2.Rect, val T) {
	min := h.PositionToIndex(phy2.Pos(bounds.Min))
	max := h.PositionToIndex(phy2.Pos(bounds.Max))

	for x := min.X; x <= max.X; x++ {
		for y := min.Y; y <= max.Y; y++ {
			bucket := h.GetBucket(Index{x, y})
			bucket.Add(bounds, val)
		}
	}
}

func (h *Hashmap[T]) PositionToIndex(pos phy2.Pos) Index {
	x := pos.X
	y := pos.Y
	xPos := (int(x) + (h.chunksize[0]/2)) / h.chunksize[0]
	yPos := (int(y) + (h.chunksize[1]/2))/ h.chunksize[1]

	// Adjust for negatives
	if x < float64(-h.chunksize[0] / 2) {
		xPos -= 1
	}
	if y < float64(-h.chunksize[1] / 2) {
		yPos -= 1
	}

	return Index{xPos, yPos}
}

func (h *Hashmap[T]) Check(bounds phy2.Rect) []T {
	min := h.PositionToIndex(phy2.Pos(bounds.Min))
	max := h.PositionToIndex(phy2.Pos(bounds.Max))

	// TODO - there's probably more efficient ways to deduplicate than a map here
	set := make(map[T]bool)

	for x := min.X; x <= max.X; x++ {
		for y := min.Y; y <= max.Y; y++ {
			bucket := h.GetBucket(Index{x, y})
			collisions := bucket.Check(bounds)
			for i := range collisions {
				set[collisions[i].item] = true
			}
		}
	}

	ret := make([]T, 0, len(set))
	for k := range set {
		ret = append(ret, k)
	}

	return ret
}
