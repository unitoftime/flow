package render

import (
	"fmt"
	"time"

	"github.com/unitoftime/ecs"

	"github.com/unitoftime/flow/particle"
	"github.com/unitoftime/glitch"
)

// TODO - Any way to do optional so that I can do all of this in a single loop?
func InterpolateParticles(world *ecs.World, dt time.Duration) {
	// Lifetime
	{
		ecs.Map(world, func(id ecs.Id, life *particle.Lifetime) {
		// view := ecs.ViewAll(world, &particle.Lifetime{})
		// view.Map(func(id ecs.Id, comp ...interface{}) {
		// 	life := comp[0].(*particle.Lifetime)
			life.Remaining -= dt

			if life.Remaining <= 0 {
				fmt.Println("NEED TO DELETE PARTICLE")
				// ecs.Tag(world, id, "delete")// TODO - Delete particles
			}
		})
	}

	// Color
	{
		ecs.Map3(world, func(id ecs.Id, life *particle.Lifetime, color *particle.Color, sprite *Sprite) {
		// view := ecs.ViewAll(world, &particle.Lifetime{}, &particle.Color{}, &Sprite{})
		// view.Map(func(id ecs.Id, comp ...interface{}) {
			// life := comp[0].(*particle.Lifetime)
			// color := comp[1].(*particle.Color)
			// sprite := comp[2].(*Sprite)

			sprite.Color = color.Get(life.Ratio())
		})
	}

	// Size
	{
		ecs.Map3(world, func(id ecs.Id, life *particle.Lifetime, size *particle.Size, sprite *Sprite) {
		// view := ecs.ViewAll(world, &particle.Lifetime{}, &particle.Size{}, &Sprite{})
		// view.Map(func(id ecs.Id, comp ...interface{}) {
			// life := comp[0].(*particle.Lifetime)
			// size := comp[1].(*particle.Size)
			// sprite := comp[2].(*Sprite)

			newSize := size.Get(life.Ratio())
			spriteBounds := sprite.Bounds()
			sprite.Scale = glitch.Vec2{float32(newSize[0]) / spriteBounds.W(), float32(newSize[1]) / spriteBounds.H()}
		})
	}
}
