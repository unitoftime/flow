package render

import (
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
			life.Remaining -= dt

			if life.Remaining <= 0 {
				ecs.Delete(world, id)
			}
		})
	}

	// Color
	{
		ecs.Map3(world, func(id ecs.Id, life *particle.Lifetime, color *particle.Color, sprite *Sprite) {
			sprite.Color = color.Get(life.Ratio())
		})
	}

	// Size
	{
		ecs.Map3(world, func(id ecs.Id, life *particle.Lifetime, size *particle.Size, sprite *Sprite) {
			newSize := size.Get(life.Ratio())
			spriteBounds := sprite.Bounds()
			sprite.Scale = glitch.Vec2{float32(newSize.X) / spriteBounds.W(), float32(newSize.Y) / spriteBounds.H()}
		})
	}
}
