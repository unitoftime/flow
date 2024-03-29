package render

// import (
// 	"time"

// 	"github.com/unitoftime/ecs"

// 	"github.com/unitoftime/flow/particle"
// 	"github.com/unitoftime/glitch"
// )

// // TODO - Any way to do optional so that I can do all of this in a single loop?
// func InterpolateParticles(world *ecs.World, dt time.Duration) {
// 	// Lifetime
// 	{
// 		query := ecs.Query1[particle.Lifetime](world)
// 		query.MapId(func(id ecs.Id, life *particle.Lifetime) {
// 			life.Remaining -= dt

// 			if life.Remaining <= 0 {
// 				ecs.Delete(world, id)
// 			}
// 		})
// 	}

// 	// Color
// 	{
// 		query := ecs.Query3[particle.Lifetime, particle.Color, Sprite](world)
// 		query.MapId(func(id ecs.Id, life *particle.Lifetime, color *particle.Color, sprite *Sprite) {
// 			sprite.Color = glitch.FromNRGBA(color.Get(life.Ratio()))
// 		})
// 	}

// 	// Size
// 	{
// 		query := ecs.Query3[particle.Lifetime, particle.Size, Sprite](world)
// 		query.MapId(func(id ecs.Id, life *particle.Lifetime, size *particle.Size, sprite *Sprite) {
// 			newSize := size.Get(life.Ratio())
// 			spriteBounds := sprite.Bounds()
// 			sprite.Scale = glitch.Vec2{newSize.X / spriteBounds.W(), newSize.Y / spriteBounds.H()}
// 		})
// 	}
// }
