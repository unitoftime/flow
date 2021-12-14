package physics

import (
	"time"

	"github.com/ungerik/go3d/float64/vec2"

	"github.com/unitoftime/ecs"
)

type Transform struct {
	X, Y float64
}

type Rigidbody struct {
	Mass float64
	Velocity vec2.T
}

type Input struct {
	Up, Down, Left, Right bool
}

func HandleInput(world *ecs.World, dt time.Duration) {
	view := ecs.ViewAll(world, &Input{}, &Transform{})
	view.Map(func(id ecs.Id, comp ...interface{}) {
		input := comp[0].(*Input)
		transform := comp[1].(*Transform)
		// ecs.Each(engine, Input{}, func(id ecs.Id, a interface{}) {
		// input := a.(Input)
		// // Note: 100 good starting point, 200 seemed like a good max
		// speed := 125.0
		// transform := Transform{}
		// ok := ecs.Read(engine, id, &transform)
		// if !ok { return }

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

		// ecs.Write(engine, id, transform)
	})
}

// Applies rigidbody physics
func RigidbodyPhysics(world *ecs.World, dt time.Duration) {
	view := ecs.ViewAll(world, &Transform{}, &Rigidbody{})
	view.Map(func(id ecs.Id, comp ...interface{}) {
		transform := comp[0].(*Transform)
		rigidbody := comp[1].(*Rigidbody)

		newTransform := vec2.T{transform.X, transform.Y}
		delta := rigidbody.Velocity.Scaled(dt.Seconds())
		newTransform.Add(&delta)
		transform.X = newTransform[0]
		transform.Y = newTransform[1]
	})
}
