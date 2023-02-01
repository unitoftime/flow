package phy2

import (
	"math"
	"github.com/unitoftime/ecs"
)

type CollisionLayer uint8

func (c CollisionLayer) Mask(layer CollisionLayer) bool {
	return (c & layer) > 0 // One layer must overlap for the layermask to work
}


// This tracks the list of current collisions
type ColliderCache struct {
	Current []ecs.Id
	Last []ecs.Id
	NewCollisions []ecs.Id // This list contains all new collisions
}

func NewColliderCache() ColliderCache {
	return ColliderCache{
		Current: make([]ecs.Id, 0),
		Last: make([]ecs.Id, 0),
		NewCollisions: make([]ecs.Id, 0),
	}
}

func (c *ColliderCache) Add(id ecs.Id) {
	c.Current = append(c.Current, id)

	for i := range c.Last {
		if c.Last[i] == id { return } // Exit early, because this one was in the last frame
	}
	// Else if we get here, then the id wasn't in the last frame list
	c.NewCollisions = append(c.NewCollisions, id)
}

func (c *ColliderCache) Clear() {
	last := c.Last
	current := c.Current

	c.Last = current
	c.Current = last[:0]

	c.NewCollisions = c.NewCollisions[:0]
}

type CircleCollider struct {
	CenterX, CenterY float64 // TODO - right now this holds the entire position of the circle (relative to world space). You might consider stripping that out though
	Radius float64
	HitLayer CollisionLayer
	Layer CollisionLayer
	Disabled bool // If set true, this collider won't collide with anything
}

func NewCircleCollider(radius float64) CircleCollider {
	return CircleCollider{
		Radius: radius,
	}
}

func (c *CircleCollider) LayerMask(layer CollisionLayer) bool {
	return (c.HitLayer & layer) > 0 // One layer must overlap for the layermask to work
}

func (c *CircleCollider) Bounds() Rect {
	return Rect{
		Min: Vec{c.CenterX - c.Radius, c.CenterY - c.Radius},
		Max: Vec{c.CenterX + c.Radius, c.CenterY + c.Radius},
	}
}

func (c *CircleCollider) Contains(yProjection float64, pos Pos) bool {
	// dx := transform.X - c.CenterX
	// dy := transform.Y - c.CenterY
	// dist := math.Hypot(dx, yProjection * dy)
	// return dist < c.Radius
	dx := pos.X - c.CenterX
	dy := pos.Y - c.CenterY
	dist := math.Hypot(dx, yProjection * dy)
	return dist < c.Radius
}

// TODO - Maybe pass position based delta into collider?
func (c *CircleCollider) Collides(yProjection float64, c2 *CircleCollider) bool {
	return !c.Disabled && !c2.Disabled && c.Overlaps(yProjection, c2)
}

func (c *CircleCollider) Overlaps(yProjection float64, c2 *CircleCollider) bool {
	dx := c2.CenterX - c.CenterX
	dy := c2.CenterY - c.CenterY
	dist := math.Hypot(dx, yProjection * dy)
	return dist < (c.Radius + c2.Radius)
}


type HashPosition struct {
	X, Y int32
}

type SpatialBucket struct {
	position HashPosition
	ids []ecs.Id
}
func NewSpatialBucket(hashPos HashPosition) *SpatialBucket {
	return &SpatialBucket{
		position: hashPos,
		ids: make([]ecs.Id, 0),
	}
}

// This holds a spatial hash of objects placed inside
type SpatialHash struct {
	bucketSize float64
	buckets map[HashPosition]*SpatialBucket
}

// TODO - pass in world dimensions?
// TODO - 2d bucket sizes?
func NewSpatialHash(bucketSize float64) *SpatialHash {
	return &SpatialHash{
		bucketSize: bucketSize,
		buckets: make(map[HashPosition]*SpatialBucket),
	}
}

func (s *SpatialHash) AddCircle(id ecs.Id, circle *CircleCollider) {
	hashPos := s.ToHashPosition(circle.CenterX, circle.CenterY)

	bucket, ok := s.buckets[hashPos]
	if !ok {
		bucket = NewSpatialBucket(hashPos)
		s.buckets[hashPos] = bucket
	}

	bucket.ids = append(bucket.ids, id)
}

func (s *SpatialHash) ToHashPosition(x, y float64) HashPosition {
	bucketX := int32(math.Floor(x/s.bucketSize))
	bucketY := int32(math.Floor(y/s.bucketSize))

	return HashPosition{bucketX, bucketY}
}
