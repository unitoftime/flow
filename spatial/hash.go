package spatial

import (
	"math"
	"slices"

	"github.com/unitoftime/flow/glm"
)

type arrayMap[T any] struct {
	topRight [][]*T
	topLeft  [][]*T
	botRight [][]*T
	botLeft  [][]*T
}

func newArrayMap[T any](size int) *arrayMap[T] {
	size = size / 2 // Note: We cut in half b/c we use 4 quadrants
	m := &arrayMap[T]{
		topRight: make([][]*T, size),
		topLeft:  make([][]*T, size),
		botRight: make([][]*T, size),
		botLeft:  make([][]*T, size),
	}

	for i := range m.topRight {
		m.topRight[i] = make([]*T, size)
	}
	for i := range m.topLeft {
		m.topLeft[i] = make([]*T, size)
	}
	for i := range m.botRight {
		m.botRight[i] = make([]*T, size)
	}
	for i := range m.botLeft {
		m.botLeft[i] = make([]*T, size)
	}

	return m
}

func (m *arrayMap[T]) safePut(slice [][]*T, x, y int, t *T) [][]*T {
	if x < len(slice) {
		if y < len(slice[x]) {
			slice[x][y] = t
			return slice
		} else {
			slice[x] = slices.Grow(slice[x], y-len(slice[x])+1)
			slice[x] = slice[x][:y+1]
			slice[x][y] = t
			return slice
		}
	}

	slice = slices.Grow(slice, x-len(slice)+1)
	slice = slice[:x+1]

	slice[x] = slices.Grow(slice[x], y-len(slice[x])+1)
	slice[x] = slice[x][:y+1]
	slice[x][y] = t
	return slice
}

func (m *arrayMap[T]) Put(x, y int, t *T) {
	if x >= 0 {
		if y >= 0 {
			m.topRight = m.safePut(m.topRight, x, y, t)
		} else {
			m.botRight = m.safePut(m.botRight, x, -y, t)
		}
	} else {
		if y >= 0 {
			m.topLeft = m.safePut(m.topLeft, -x, y, t)
		} else {
			m.botLeft = m.safePut(m.botLeft, -x, -y, t)
		}
	}
}

func (m *arrayMap[T]) safeGet(slice [][]*T, x, y int) (*T, bool) {
	if x >= len(slice) {
		return nil, false
	}
	if y >= len(slice[x]) {
		return nil, false
	}

	isNil := (slice[x][y] == nil)
	return slice[x][y], !isNil
}

func (m *arrayMap[T]) Get(x, y int) (*T, bool) {
	if x >= 0 {
		if y >= 0 {
			return m.safeGet(m.topRight, x, y)
		} else {
			return m.safeGet(m.botRight, x, -y)
		}
	} else {
		if y >= 0 {
			return m.safeGet(m.topLeft, -x, y)
		} else {
			return m.safeGet(m.botLeft, -x, -y)
		}
	}
}

func (m *arrayMap[T]) ForEachValue(lambda func(t *T)) {
	for x := range m.topRight {
		for y := range m.topRight[x] {
			val := m.topRight[x][y]
			if val == nil {
				continue
			}
			lambda(val)
		}
	}
	for x := range m.topLeft {
		for y := range m.topLeft[x] {
			val := m.topLeft[x][y]
			if val == nil {
				continue
			}
			lambda(val)
		}
	}
	for x := range m.botLeft {
		for y := range m.botLeft[x] {
			val := m.botLeft[x][y]
			if val == nil {
				continue
			}
			lambda(val)
		}
	}
	for x := range m.botRight {
		for y := range m.botRight[x] {
			val := m.botRight[x][y]
			if val == nil {
				continue
			}
			lambda(val)
		}
	}
}

type Index struct {
	X, Y int
}

type BucketItem[T comparable] struct {
	shape Shape
	item  T
}

type Bucket[T comparable] struct {
	List []BucketItem[T]
}

func NewBucket[T comparable]() *Bucket[T] {
	return &Bucket[T]{
		List: make([]BucketItem[T], 0),
	}
}
func (b *Bucket[T]) Add(shape Shape, val T) {
	b.List = append(b.List, BucketItem[T]{
		shape: shape,
		item:  val,
	})
}
func (b *Bucket[T]) Remove(val T) {
	b.List = slices.DeleteFunc(b.List, func(a BucketItem[T]) bool {
		return (a.item == val)
	})
}

func (b *Bucket[T]) Clear() {
	b.List = b.List[:0]
}
func (b *Bucket[T]) Check(colSet *CollisionSet[T], shape Shape) {
	for i := range b.List {
		if shape.Intersects(b.List[i].shape) {
			colSet.Add(b.List[i].item)
		}
	}
}
func (b *Bucket[T]) Collides(shape Shape) bool {
	for i := range b.List {
		if shape.Intersects(b.List[i].shape) {
			return true
		}
	}
	return false
}

func (b *Bucket[T]) FindClosest(shape Shape) (BucketItem[T], bool) {
	center := shape.Bounds.Center()
	distSquared := math.MaxFloat64
	ret := BucketItem[T]{}
	set := false
	for i := range b.List {
		if shape.Intersects(b.List[i].shape) {
			ds := center.DistSq(b.List[i].shape.Bounds.Center())
			if ds < distSquared {
				distSquared = ds
				ret = b.List[i]
				set = true
			}
		}
	}
	return ret, set
}

// --------------------------------------------------------------------------------
type PositionHasher struct {
	size      [2]int
	sizeOver2 [2]int
	div       [2]int
}

func NewPositionHasher(size [2]int) PositionHasher {
	divX := int(math.Log2(float64(size[0])))
	divY := int(math.Log2(float64(size[1])))
	if (1<<divX) != size[0] || (1<<divY) != size[1] {
		panic("Spatial maps must have a chunksize that is a power of 2!")
	}
	return PositionHasher{
		size:      size,
		sizeOver2: [2]int{size[0] / 2, size[1] / 2},
		div:       [2]int{divX, divY},
	}
}

func (h *PositionHasher) PositionToIndex(pos glm.Vec2) Index {
	x := pos.X
	y := pos.Y
	xPos := (int(x)) >> h.div[0]
	yPos := (int(y)) >> h.div[1]

	// TODO: I dont think I need this b/c of how shift right division works (negatives stay negative)
	// // Adjust for negatives
	// if x < float64(-h.sizeOver2[0]) {
	// 	xPos -= 1
	// }
	// if y < float64(-h.sizeOver2[1]) {
	// 	yPos -= 1
	// }

	return Index{xPos, yPos}
}

// --------------------------------------------------------------------------------
// TODO: rename? ColliderMap?
type Hashmap[T comparable] struct {
	PositionHasher

	allBuckets []*Bucket[T]
	Bucket     *arrayMap[Bucket[T]]
}

func NewHashmap[T comparable](chunksize [2]int, startingSize int) *Hashmap[T] {
	return &Hashmap[T]{
		PositionHasher: NewPositionHasher(chunksize),

		allBuckets: make([]*Bucket[T], 0, 1024),
		Bucket:     newArrayMap[Bucket[T]](startingSize),
	}
}

func (h *Hashmap[T]) Clear() {
	for i := range h.allBuckets {
		h.allBuckets[i].Clear()
	}
}

func (h *Hashmap[T]) GetBucket(index Index) *Bucket[T] {
	bucket, ok := h.Bucket.Get(index.X, index.Y)
	if !ok {
		bucket = NewBucket[T]()
		h.allBuckets = append(h.allBuckets, bucket)
		h.Bucket.Put(index.X, index.Y, bucket)
	}
	return bucket
}

func (h *Hashmap[T]) Add(shape Shape, val T) {
	min := h.PositionToIndex(shape.Bounds.Min)
	max := h.PositionToIndex(shape.Bounds.Max)

	for x := min.X; x <= max.X; x++ {
		for y := min.Y; y <= max.Y; y++ {
			bucket := h.GetBucket(Index{x, y})
			bucket.Add(shape, val)
		}
	}
}

// Warning: This is a relatively slow operation
func (h *Hashmap[T]) Remove(val T) {
	// Just try and remove the id from all buckets
	for i := range h.allBuckets {
		h.allBuckets[i].Remove(val)
	}
}

// Finds collisions and adds them directly into your collision set
func (h *Hashmap[T]) Check(colSet *CollisionSet[T], shape Shape) {
	// shape := AABB(rect)
	min := h.PositionToIndex(shape.Bounds.Min)
	max := h.PositionToIndex(shape.Bounds.Max)

	for x := min.X; x <= max.X; x++ {
		for y := min.Y; y <= max.Y; y++ {
			isBorderChunk := (x == min.X || x == max.X) || (y == min.Y || y == max.Y)
			// bucket := h.GetBucket(Index{x, y})
			bucket, ok := h.Bucket.Get(x, y)
			if !ok {
				continue
			}

			if isBorderChunk {
				// For border chunks, we need to do narrow phase too
				bucket.Check(colSet, shape)
			} else {
				// For inner chunks, we can just add everything from the bucket (much faster)
				for i := range bucket.List {
					colSet.Add(bucket.List[i].item)
				}
			}
		}
	}
}

func (h *Hashmap[T]) Collides(rect glm.Rect) bool {
	shape := AABB(rect)
	min := h.PositionToIndex(shape.Bounds.Min)
	max := h.PositionToIndex(shape.Bounds.Max)

	for x := min.X; x <= max.X; x++ {
		for y := min.Y; y <= max.Y; y++ {
			isBorderChunk := (x == min.X || x == max.X) || (y == min.Y || y == max.Y)
			// bucket := h.GetBucket(Index{x, y})
			bucket, ok := h.Bucket.Get(x, y)
			if !ok {
				continue
			}

			if isBorderChunk {
				// For border chunks, we need to do narrow phase too
				if bucket.Collides(shape) {
					return true
				}
			} else {
				// For inner chunks, we can just assume everything collides
				if len(bucket.List) > 0 {
					return true
				}
			}
		}
	}
	return false
}

func (h *Hashmap[T]) FindClosest(rect glm.Rect) (T, bool) {
	shape := AABB(rect)
	min := h.PositionToIndex(shape.Bounds.Min)
	max := h.PositionToIndex(shape.Bounds.Max)

	center := shape.Bounds.Center()
	distSquared := math.MaxFloat64
	var ret T
	set := false

	for x := min.X; x <= max.X; x++ {
		for y := min.Y; y <= max.Y; y++ {
			bucket, ok := h.Bucket.Get(x, y)
			if !ok {
				continue
			}

			// For border chunks, we need to do narrow phase too
			closest, ok := bucket.FindClosest(shape)
			if ok {
				ds := center.DistSq(closest.shape.Bounds.Center())
				if ds < distSquared {
					distSquared = ds
					ret = closest.item
					set = true
				}
			}
		}
	}

	return ret, set
}

// Adds the collisions directly into your collision set. This one doesnt' do any narrow phase detection. It returns all objects that collide with the same chunk
func (h *Hashmap[T]) BroadCheck(colSet CollisionSet[T], rect glm.Rect) {
	shape := AABB(rect)
	min := h.PositionToIndex(shape.Bounds.Min)
	max := h.PositionToIndex(shape.Bounds.Max)

	for x := min.X; x <= max.X; x++ {
		for y := min.Y; y <= max.Y; y++ {
			bucket := h.GetBucket(Index{x, y})
			for i := range bucket.List {
				colSet.Add(bucket.List[i].item)
			}
		}
	}
}

// TODO - there's probably more efficient ways to deduplicate than a map here?
// type CollisionSet[T comparable] map[T]struct{}
// func NewCollisionSet[T comparable](cap int) CollisionSet[T] {
// 	return make(CollisionSet[T], cap)
// }
// func (s CollisionSet[T]) Add(t T) {
// 	s[t] = struct{}{}
// }
// func (s CollisionSet[T]) Clear() {
// 	// clear(s) // TODO: Is this slow?
// 	// Clearing Optimization: https://go.dev/doc/go1.11#performance-compiler
// 	for k := range s {
// 		delete(s, k)
// 	}
// }

type CollisionSet[T comparable] struct {
	List []T
}

func NewCollisionSet[T comparable](cap int) *CollisionSet[T] {
	return &CollisionSet[T]{
		List: make([]T, cap),
	}
}
func (s *CollisionSet[T]) Add(t T) {
	for i := range s.List {
		if s.List[i] == t {
			return // Already added
		}
	}
	s.List = append(s.List, t)
}
func (s *CollisionSet[T]) Clear() {
	s.List = s.List[:0]
}
