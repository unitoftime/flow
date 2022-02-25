package physics

import (
	"time"

	// "github.com/ungerik/go3d/float64/vec2"

	"github.com/unitoftime/ecs"
)

type Transform struct {
	X, Y float64
}

type Rigidbody struct {
	Mass float64
	Velocity Vec2
}

type Input struct {
	Up, Down, Left, Right bool
}

func HandleInput(world *ecs.World, dt time.Duration) {
	// view := ecs.ViewAll2[Input, Transform](world)
	ecs.Map2(world, func(id ecs.Id, input *Input, transform *Transform) {
		// Note: 100 good starting point, 200 seemed like a good max
		speed := 125.0

		if input.Left {
			transform.X -= speed * dt.Seconds()
		}
		if input.Right {
			transform.X += speed * dt.Seconds()
		}
		if input.Up {
			transform.Y += speed * dt.Seconds()
		}
		if input.Down {
			transform.Y -= speed * dt.Seconds()
		}
	})
}

// Applies rigidbody physics
func RigidbodyPhysics(world *ecs.World, dt time.Duration) {
	ecs.Map2(world, func(id ecs.Id, transform *Transform, rigidbody *Rigidbody) {
		// newTransform := vec2.T{transform.X, transform.Y}
		newTransform := Vec2{transform.X, transform.Y}
		delta := rigidbody.Velocity.Scaled(dt.Seconds())
		newTransform = newTransform.Add(delta)
		transform.X = newTransform.X
		transform.Y = newTransform.Y
	})
}
