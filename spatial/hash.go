package spatial

import (
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


type Index struct {
	X, Y int
}

func PositionToIndex(chunksize [2]int, pos phy2.Pos) Index {
	x := pos.X
	y := pos.Y
	xPos := (int(x) + (chunksize[0]/2)) / chunksize[0]
	yPos := (int(y) + (chunksize[1]/2))/ chunksize[1]

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
func (b *Bucket[T]) Check(colSet map[T]struct{}, shape phy2.Rect) {
	for i := range b.List {
		if shape.Intersects(b.List[i].shape) {
			colSet[b.List[i].item] = struct{}{}
		}
	}
}

// TODO: rename? ColliderMap?
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
func (h *Hashmap[T]) Check(colSet CollisionSet[T], shape phy2.Rect) {
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
					colSet[bucket.List[i].item] = struct{}{}
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
				colSet[bucket.List[i].item] = struct{}{}
			}
		}
	}
}

// TODO - there's probably more efficient ways to deduplicate than a map here?
type CollisionSet[T comparable] map[T]struct{}
func NewCollisionSet[T comparable](cap int) CollisionSet[T] {
	return make(CollisionSet[T], cap)
}
func (s CollisionSet[T]) Clear() {
	clear(s)
}
