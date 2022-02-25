package render

import (
	"image/color"

	"github.com/unitoftime/glitch"

	"github.com/unitoftime/ecs"
	"github.com/unitoftime/flow/physics"
)

// Represents multiple sprites
type MultiSprite struct {
	Sprites []Sprite
}
func NewMultiSprite(sprites ...Sprite) MultiSprite {
	m := MultiSprite{
		Sprites: make([]Sprite, len(sprites)),
	}
	for i := range sprites {
		m.Sprites[i] = sprites[i]
	}
	return m
}

type Sprite struct {
	*glitch.Sprite
	Color color.NRGBA // TODO - performance on interfaces vs structs?
	Scale glitch.Vec2
}
func NewSprite(sprite *glitch.Sprite) Sprite {
	return Sprite{
		Sprite: sprite,
		Color: color.NRGBA{255, 255, 255, 255},
		Scale: glitch.Vec2{1, 1},
	}
}

type Keybinds struct {
	Up, Down, Left, Right glitch.Key
}

// Note: val should probably be between 0 and 1
func Interpolate(A, B glitch.Vec2, lowerBound, upperBound float32) glitch.Vec2 {
	delta := B.Sub(A)
	dMag := delta.Len()

	interpValue := float32(0.0)
	if dMag > upperBound {
		interpValue = 1.0
	} else if dMag > lowerBound {
		// y - y1 = m(x - x1)
		slope := 1/(upperBound - lowerBound)
		interpValue = slope * (dMag - lowerBound) + 0
	}

	deltaScaled := delta.Scaled(interpValue)
	return A.Add(deltaScaled)
}

// TODO - interpolate based off of the time till the next fixedTimeStep?
// func InterpolateSpritePositions(world *ecs.World, dt time.Duration) {
// 	view := ecs.ViewAll(world, &Sprite{}, &physics.Transform{})
// 	view.Map(func(id ecs.Id, comp ...interface{}) {
// 		sprite := comp[0].(*Sprite)
// 		transform := comp[1].(*physics.Transform)
// 		// ecs.Each(engine, Sprite{}, func(id ecs.Id, a interface{}) {
// 		// 	sprite := a.(Sprite)

// 		// transform := physics.Transform{}
// 		// ok := ecs.Read(engine, id, &transform)
// 		// if !ok { return }
// 		physicsPosition := pixel.V(transform.X, transform.Y)

// 		// TODO - make configurable
// 		// sprite.Position = physicsPosition
// 		sprite.Position = Interpolate(sprite.Position, physicsPosition, 1.0, 16.0)
// 		// ecs.Write(engine, id, sprite)
// 	})
// }

// TODO - how to do optional components? with some default val?
func DrawSprites(pass *glitch.RenderPass, world *ecs.World) {
	ecs.Map2(world, func(id ecs.Id, sprite *Sprite, t *physics.Transform) {
		mat := glitch.Mat4Ident
		mat.Scale(sprite.Scale[0], sprite.Scale[1], 1.0).Translate(float32(t.X), float32(t.Y), 0)

		// TODO - I think there's some mistakes here with premultiplied vs non premultiplied alpha
		col := glitch.RGBA{float32(sprite.Color.R)/255.0, float32(sprite.Color.G)/255.0, float32(sprite.Color.B)/255.0, float32(sprite.Color.A)/255.0}
		sprite.DrawColorMask(pass, mat, col)
	})
}

func DrawMultiSprites(pass *glitch.RenderPass, world *ecs.World) {
	ecs.Map2(world, func(id ecs.Id, mSprite *MultiSprite, t *physics.Transform) {
		for _, sprite := range mSprite.Sprites {
			mat := glitch.Mat4Ident
			mat.Scale(sprite.Scale[0], sprite.Scale[1], 1.0).Translate(float32(t.X), float32(t.Y), 0)

			// TODO - I think there's some mistakes here with premultiplied vs non premultiplied alpha
			col := glitch.RGBA{float32(sprite.Color.R)/255.0, float32(sprite.Color.G)/255.0, float32(sprite.Color.B)/255.0, float32(sprite.Color.A)/255.0}
			sprite.DrawColorMask(pass, mat, col)
		}
	})
}

func CaptureInput(win *glitch.Window, world *ecs.World) {
	ecs.Map2(world, func(id ecs.Id, keybinds *Keybinds, input *physics.Input) {
	// view := ecs.ViewAll(world, &Keybinds{}, &physics.Input{})
	// view.Map(func(id ecs.Id, comp ...interface{}) {
	// ecs.Each(engine, Keybinds{}, func(id ecs.Id, a interface{}) {
	// 	keybinds := a.(Keybinds)

		// input := physics.Input{}
		// ok := ecs.Read(engine, id, &input)
		// if !ok { return }

		// keybinds := comp[0].(*Keybinds)
		// input := comp[1].(*physics.Input)

		input.Left = false
		input.Right = false
		input.Up = false
		input.Down = false

		if win.Pressed(keybinds.Left) {
			input.Left = true
		}
		if win.Pressed(keybinds.Right) {
			input.Right = true
		}
		if win.Pressed(keybinds.Up) {
			input.Up = true
		}
		if win.Pressed(keybinds.Down) {
			input.Down = true
		}

		// ecs.Write(engine, id, input)
	})
}
