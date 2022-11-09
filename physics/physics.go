package physics

import (
	"time"
	"math"

	"github.com/unitoftime/ecs"
)

type Transform struct {
	X, Y float64
	Height float64
}

// TODO - combine this with vec2 primitives
func (a *Transform) DistanceTo(b *Transform) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return math.Sqrt(dx * dx + dy * dy)
}

type Rigidbody struct {
	Mass float64
	Velocity Vec2
}

type Input struct {
	Up, Down, Left, Right bool
}

// Applies rigidbody physics
func RigidbodyPhysics(world *ecs.World, dt time.Duration) {
	ecs.Map2(world, func(id ecs.Id, transform *Transform, rigidbody *Rigidbody) {
		newTransform := Vec2{transform.X, transform.Y}
		delta := rigidbody.Velocity.Scaled(dt.Seconds())
		newTransform = newTransform.Add(delta)
		transform.X = newTransform.X
		transform.Y = newTransform.Y
	})
}
