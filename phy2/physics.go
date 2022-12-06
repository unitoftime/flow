package phy2

import (
	// "time"
	// "math"

	// "github.com/unitoftime/ecs"
)

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

func (v Pos) Len() float64 {
	return Vec2(v).Len()
}

func (v Pos) Scaled(s float64) Pos {
	return Pos(Vec2(v).Scaled(s))
}

func (v Pos) Rotated(radians float64) Pos {
	return Pos(Vec2(v).Rotated(radians))
}

type Scale Vec2
type Rotation float64

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
