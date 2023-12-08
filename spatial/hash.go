package spatial

import (
	"slices"

	"github.com/unitoftime/flow/phy2"
)

// TODO: eventually use shapes
// type ShapeType uint8
// const (
// 	ShapeRect ShapeType = iota
// 	ShapeCircle
// )

// type Shape struct {
// 	Type ShapeType
// 	Bounds phy2.Rect
// }
// func Rect(rect phy2.Rect) Shape {
// 	return Shape{
// 		Type: ShapeRect,
// 		Bounds: rect,
// 	}
// }
// func Circle(rect phy2.Rect) Shape {
// 	return Shape{
// 		Type: ShapeCircle,
// 		Bounds: rect,
// 	}
// }
// func (s Shape) Rect() phy2.Rect {
// 	return s.Bounds
// }
// func (s Shape) Circle() phy2.Circle {
// 	phy2.NewCircle(s.Bounds.Center())
// }

// func (s Shape) Intersects(s2 Shape) bool {
// 	if s.Type == ShapeRect {
// 		if s2.Type == ShapeRect {
// 			return s.Rect().Intersects(s2.Bounds)
// 		} else if s2.Type == ShapeCircle {
// 			return s.Rect().IntersectCircle(s2.Circle())
// 		}
// 	} else if s.Type == ShapeCircle {
// 		if s2.Type == ShapeRect {
// 			return s.Bounds.Intersects(s2.Bounds)
// 		} else if s2.Type == ShapeCircle {
// 			return s.Circle().IntersectCircle(s2.Circle())
// 		}
// 	}
// }

type arrayMap[T any] struct {
	topRight [][]*T
	topLeft [][]*T
	botRight [][]*T
	botLeft [][]*T
}
func newArrayMap[T any]() *arrayMap[T] {
	return &arrayMap[T]{
		topRight: make([][]*T, 0),
		topLeft: make([][]*T, 0),
		botRight: make([][]*T, 0),
		botLeft: make([][]*T, 0),
	}
}

func (m *arrayMap[T]) safePut(slice [][]*T, x, y int, t *T) [][]*T {
	if x < len(slice) {
		if y < len(slice[x]) {
			slice[x][y] = t
			return slice
		} else {
			slice[x] = slices.Grow(slice[x], y - len(slice[x]) + 1)
			slice[x] = slice[x][:y+1]
			slice[x][y] = t
			return slice
		}
	}

	slice = slices.Grow(slice, x - len(slice) + 1)
	slice = slice[:x+1]

	slice[x] = slices.Grow(slice[x], y - len(slice[x]) + 1)
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
			// m.botRight[x][y] = t
		}
	} else {
		if y >= 0 {
			m.topLeft = m.safePut(m.topLeft, -x, y, t)
			// m.topLeft[x][y] = t
		} else {
			m.botLeft = m.safePut(m.botLeft, -x, -y, t)
			// m.botLeft[x][y] = t
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

	// if slice[x][y] == nil {
	// 	return nil, false
	// }

	isNil := (slice[x][y] == nil)
	return slice[x][y], !isNil
}

func (m *arrayMap[T]) Get(x, y int) (*T, bool) {
	if x >= 0 {
		if y >= 0 {
			return m.safeGet(m.topRight, x, y)
		} else {
			return m.safeGet(m.botRight, x, -y)
			// m.botRight[x][y] = t
		}
	} else {
		if y >= 0 {
			return m.safeGet(m.topLeft, -x, y)
			// m.topLeft[x][y] = t
		} else {
			return m.safeGet(m.botLeft, -x, -y)
			// m.botLeft[x][y] = t
		}
	}
}

func (m *arrayMap[T]) ForEachValue(lambda func(t *T)) {
	for x := range m.topRight {
		for y := range m.topRight[x] {
			val := m.topRight[x][y]
			if val == nil { continue }
			lambda(val)
		}
	}
	for x := range m.topLeft {
		for y := range m.topLeft[x] {
			val := m.topLeft[x][y]
			if val == nil { continue }
			lambda(val)
		}
	}
	for x := range m.botLeft {
		for y := range m.botLeft[x] {
			val := m.botLeft[x][y]
			if val == nil { continue }
			lambda(val)
		}
	}
	for x := range m.botRight {
		for y := range m.botRight[x] {
			val := m.botRight[x][y]
			if val == nil { continue }
			lambda(val)
		}
	}
}

type Index struct {
	X, Y int
}

func PositionToIndex(chunksize [2]int, pos phy2.Pos) Index {
	x := pos.X
	y := pos.Y
	xPos := (int(x) + (chunksize[0]/2)) / chunksize[0]
	yPos := (int(y) + (chunksize[1]/2)) / chunksize[1]

	// Adjust for negatives
	if x < float64(-chunksize[0] / 2) {
		xPos -= 1
	}
	if y < float64(-chunksize[1] / 2) {
		yPos -= 1
	}

	return Index{xPos, yPos}
}


type BucketItem[T comparable] struct {
	shape phy2.Rect
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
func (b *Bucket[T]) Add(shape phy2.Rect, val T) {
	b.List = append(b.List, BucketItem[T]{
		shape: shape,
		item: val,
	})
}
func (b *Bucket[T]) Clear() {
	b.List = b.List[:0]
}
func (b *Bucket[T]) Check(colSet *CollisionSet[T], shape phy2.Rect) {
	for i := range b.List {
		if shape.Intersects(b.List[i].shape) {
			colSet.Add(b.List[i].item)
		}
	}
}

// TODO: rename? ColliderMap?
type Hashmap[T comparable] struct {
	allBuckets []*Bucket[T]
	Bucket *arrayMap[Bucket[T]]
	// Bucket map[Index]*Bucket[T]
	chunksize [2]int
}

func NewHashmap[T comparable](chunksize [2]int) *Hashmap[T] {
	return &Hashmap[T]{
		allBuckets: make([]*Bucket[T], 0, 1024),
		// Bucket: make(map[Index]*Bucket[T]),
		Bucket: newArrayMap[Bucket[T]](),
		chunksize: chunksize,
	}
}

func (h *Hashmap[T]) Clear() {
	for _, b := range h.allBuckets {
		b.Clear()
	}
	// h.Bucket.ForEachValue(func(b *Bucket[T]) {
	// 	b.Clear()
	// })
	// for _, b := range h.Bucket {
	// 	b.Clear()
	// }
}

func (h *Hashmap[T]) GetBucket(index Index) *Bucket[T] {
	bucket, ok := h.Bucket.Get(index.X, index.Y)
	if !ok {
		bucket = NewBucket[T]()
		h.allBuckets = append(h.allBuckets, bucket)
		h.Bucket.Put(index.X, index.Y, bucket)
	}
	return bucket
	// bucket, ok := h.Bucket[index]
	// if !ok {
	// 	bucket = NewBucket[T]()
	// 	h.Bucket[index] = bucket
	// }
	// return bucket
}

func (h *Hashmap[T]) Add(shape phy2.Rect, val T) {
	min := PositionToIndex(h.chunksize, phy2.Pos(shape.Min))
	max := PositionToIndex(h.chunksize, phy2.Pos(shape.Max))

	for x := min.X; x <= max.X; x++ {
		for y := min.Y; y <= max.Y; y++ {
			bucket := h.GetBucket(Index{x, y})
			bucket.Add(shape, val)
		}
	}
}

// Finds collisions and adds them directly into your collision set
func (h *Hashmap[T]) Check(colSet *CollisionSet[T], shape phy2.Rect) {
	min := PositionToIndex(h.chunksize, phy2.Pos(shape.Min))
	max := PositionToIndex(h.chunksize, phy2.Pos(shape.Max))

	for x := min.X; x <= max.X; x++ {
		for y := min.Y; y <= max.Y; y++ {
			isBorderChunk := (x == min.X || x == max.X) && (y == min.Y || y == max.Y)
			bucket := h.GetBucket(Index{x, y})
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

// Adds the collisions directly into your collision set. This one doesnt' do any narrow phase detection. It returns all objects that collide with the same chunk
func (h *Hashmap[T]) BroadCheck(colSet CollisionSet[T], shape phy2.Rect) {
	min := PositionToIndex(h.chunksize, phy2.Pos(shape.Min))
	max := PositionToIndex(h.chunksize, phy2.Pos(shape.Max))

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

type CollisionSet[T comparable] struct{
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
