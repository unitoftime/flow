package phy2

import (
	"math"
)

//cod:struct
type Pos Vec2

func (v Pos) Add(v2 Pos) Pos {
	return Pos(Vec2(v).Add(Vec2(v2)))
}

func (v Pos) Sub(v2 Pos) Pos {
	return Pos(Vec2(v).Sub(Vec2(v2)))
}

func (v Pos) Norm() Pos {
	return Pos(Vec2(v).Norm())
}

func (v Pos) Dot(u Pos) float64 {
	return Vec2(v).Dot(Vec2(u))
}

func (v Pos) Dist(u Pos) float64 {
	return Vec2(v).Dist(Vec2(u))
}

func (v Pos) Len() float64 {
	return Vec2(v).Len()
}

func (v Pos) Scaled(s float64) Pos {
	return Pos(Vec2(v).Scaled(s))
}

func (v Pos) Rotated(radians float64) Pos {
	return Pos(Vec2(v).Rotated(radians))
}

func (v Pos) Angle() float64 {
	return math.Atan2(v.Y, v.X)
}

//cod:struct
type Vel Vec2

func (v Vel) Add(v2 Vel) Vel {
	return Vel(Vec2(v).Add(Vec2(v2)))
}

func (v Vel) Sub(v2 Vel) Vel {
	return Vel(Vec2(v).Sub(Vec2(v2)))
}

func (v Vel) Norm() Vel {
	return Vel(Vec2(v).Norm())
}

func (v Vel) Dot(u Vel) float64 {
	return Vec2(v).Dot(Vec2(u))
}

func (v Vel) Len() float64 {
	return Vec2(v).Len()
}

func (v Vel) Scaled(s float64) Vel {
	return Vel(Vec2(v).Scaled(s))
}

func (v Vel) Rotated(radians float64) Vel {
	return Vel(Vec2(v).Rotated(radians))
}

func (v Vel) Angle() float64 {
	return math.Atan2(v.Y, v.X)
}

//cod:struct
type Scale Vec2

//cod:struct
type Rotation float64

//cod:struct
type Rigidbody struct {
	Mass float64
	Velocity Vec2
}

// // Applies rigidbody physics
// func RigidbodyPhysics(world *ecs.World, dt time.Duration) {
// 	ecs.Map2(world, func(id ecs.Id, transform *Transform, rigidbody *Rigidbody) {
// 		newTransform := Vec2{transform.X, transform.Y}
// 		delta := rigidbody.Velocity.Scaled(dt.Seconds())
// 		newTransform = newTransform.Add(delta)
// 		transform.X = newTransform.X
// 		transform.Y = newTransform.Y
// 	})
// }
