package spatial

import (
	"slices"

	"github.com/unitoftime/flow/glm"
)

type PointBucketItem[T comparable] struct {
	point glm.Vec2
	item  T
}

type PointBucket[T comparable] struct {
	List []PointBucketItem[T]
}

func NewPointBucket[T comparable]() *PointBucket[T] {
	return &PointBucket[T]{
		List: make([]PointBucketItem[T], 0),
	}
}

func (b *PointBucket[T]) Add(point glm.Vec2, val T) {
	b.List = append(b.List, PointBucketItem[T]{
		point: point,
		item:  val,
	})
}

func (b *PointBucket[T]) Remove(point glm.Vec2, val T) {
	itemToRemove := PointBucketItem[T]{
		point: point,
		item:  val,
	}
	indexToRemove := slices.Index(b.List, itemToRemove)
	if indexToRemove < 0 {
		return
	} // Nothing to remove

	b.RemoveIndex(indexToRemove)
}

func (b *PointBucket[T]) RemoveIndex(idx int) {
	lastIdx := len(b.List) - 1
	b.List[idx] = b.List[lastIdx]
	b.List = b.List[:lastIdx]
}

// Remove every point that collides with the bucket
func (b *PointBucket[T]) RemoveCollides(bounds glm.Rect) {
	for i := 0; i < len(b.List); i++ {
		if !bounds.Contains(b.List[i].point) { continue } // skip if doesnt collide

		b.RemoveIndex(i)
		i--
	}
}

func (b *PointBucket[T]) Clear() {
	b.List = b.List[:0]
}

// TODO: rename? ColliderMap?
type Pointmap[T comparable] struct {
	PositionHasher

	Bucket     *arrayMap[PointBucket[T]]
	allBuckets []*PointBucket[T]
}

func NewPointmap[T comparable](chunksize [2]int, startingSize int) *Pointmap[T] {
	return &Pointmap[T]{
		allBuckets:     make([]*PointBucket[T], 0, 1024),
		Bucket:         newArrayMap[PointBucket[T]](startingSize),
		PositionHasher: NewPositionHasher(chunksize),
	}
}

func (h *Pointmap[T]) Clear() {
	for i := range h.allBuckets {
		h.allBuckets[i].Clear()
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

func (h *Pointmap[T]) Add(pos glm.Vec2, val T) {
	idx := h.PositionToIndex(pos)
	bucket := h.GetBucket(idx)
	bucket.Add(pos, val)
}

func (h *Pointmap[T]) Remove(pos glm.Vec2, val T) {
	idx := h.PositionToIndex(pos)
	bucket := h.GetBucket(idx)
	bucket.Remove(pos, val)
}

// TODO: Right now this does a broad phased check
// Adds the collisions directly into your collision list. Items are deduplicated by nature of them only existing once in this Pointmap. (ie if you add multiple of the same thing, you might get multiple out)
func (h *Pointmap[T]) BroadCheck(list []T, bounds glm.Rect) []T {
	min := h.PositionToIndex(bounds.Min)
	max := h.PositionToIndex(bounds.Max)

	// TODO: Might be nice if this spirals from inside to outside, that way its roughly sorted by distance?
	for x := min.X; x <= max.X; x++ {
		for y := min.Y; y <= max.Y; y++ {
			bucket, ok := h.Bucket.Get(x, y)
			if !ok {
				continue
			}
			// bucket := h.GetBucket(Index{x, y})
			for i := range bucket.List {
				list = append(list, bucket.List[i].item)
			}
		}
	}

	return list
}

// TODO: I think I'd rather the default for this be called "Check" then have the other be called CheckBroad or something
func (h *Pointmap[T]) NarrowCheck(list []T, bounds glm.Rect) []T {
	min := h.PositionToIndex(bounds.Min)
	max := h.PositionToIndex(bounds.Max)

	// TODO: Might be nice if this spirals from inside to outside, that way its roughly sorted by distance?
	for x := min.X; x <= max.X; x++ {
		for y := min.Y; y <= max.Y; y++ {
			bucket, ok := h.Bucket.Get(x, y)
			if !ok {
				continue
			}

			for i := range bucket.List {
				if bounds.Contains(bucket.List[i].point) {
					list = append(list, bucket.List[i].item)
				}
			}
		}
	}

	return list
}

// TODO: This only does a broadphase check. no narrow phase
// Returns true if the bounds collides with anything
func (h *Pointmap[T]) Collides(bounds glm.Rect) bool {
	min := h.PositionToIndex(bounds.Min)
	max := h.PositionToIndex(bounds.Max)

	// TODO: Might be nice if this spirals from inside to outside, that way its roughly sorted by distance?
	for x := min.X; x <= max.X; x++ {
		for y := min.Y; y <= max.Y; y++ {
			bucket, ok := h.Bucket.Get(x, y)
			if !ok {
				continue
			}
			if len(bucket.List) > 0 {
				return true
			}
		}
	}

	return false
}

// Remove every point that collides with the supplied bounds
func (h *Pointmap[T]) RemoveCollides(bounds glm.Rect) {
	min := h.PositionToIndex(bounds.Min)
	max := h.PositionToIndex(bounds.Max)

	// TODO: Might be nice if this spirals from inside to outside, that way its roughly sorted by distance?
	for x := min.X; x <= max.X; x++ {
		for y := min.Y; y <= max.Y; y++ {
			bucket, ok := h.Bucket.Get(x, y)
			if !ok {
				continue
			}

			bucket.RemoveCollides(bounds)
		}
	}
}
