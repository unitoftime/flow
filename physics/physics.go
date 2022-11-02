package physics

import (
	"time"

	"github.com/unitoftime/ecs"
)

type Transform struct {
	X, Y float64
	Height float64
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
