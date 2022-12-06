package render

import (
	"image/color"

	"github.com/unitoftime/glitch"
	"github.com/unitoftime/flow/phy2"
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
	Layer uint8
}
func NewSprite(sprite *glitch.Sprite) Sprite {
	return Sprite{
		Sprite: sprite,
		Color: color.NRGBA{255, 255, 255, 255},
		Scale: glitch.Vec2{1, 1},
		Layer: glitch.DefaultLayer,
	}
}

func (sprite *Sprite) Draw(pass *glitch.RenderPass, pos *phy2.Pos) {
	mat := glitch.Mat4Ident
	mat.Scale(sprite.Scale[0], sprite.Scale[1], 1.0).Translate(float32(pos.X), float32(pos.Y), 0)

	// TODO - I think there's some mistakes here with premultiplied vs non premultiplied alpha
	col := glitch.RGBA{float32(sprite.Color.R)/255.0, float32(sprite.Color.G)/255.0, float32(sprite.Color.B)/255.0, float32(sprite.Color.A)/255.0}
	pass.SetLayer(sprite.Layer)
	sprite.DrawColorMask(pass, mat, col)
}

// type Keybinds struct {
// 	Up, Down, Left, Right glitch.Key
// }

// // Note: val should probably be between 0 and 1
// func Interpolate(A, B glitch.Vec2, lowerBound, upperBound float32) glitch.Vec2 {
// 	delta := B.Sub(A)
// 	dMag := delta.Len()

// 	interpValue := float32(0.0)
// 	if dMag > upperBound {
// 		interpValue = 1.0
// 	} else if dMag > lowerBound {
// 		// y - y1 = m(x - x1)
// 		slope := 1/(upperBound - lowerBound)
// 		interpValue = slope * (dMag - lowerBound) + 0
// 	}

// 	deltaScaled := delta.Scaled(interpValue)
// 	return A.Add(deltaScaled)
// }

// // TODO - how to do optional components? with some default val?
// func DrawSprites(pass *glitch.RenderPass, world *ecs.World) {
// 	ecs.Map2(world, func(id ecs.Id, sprite *Sprite, t *phy2.Transform) {
// 		mat := glitch.Mat4Ident
// 		mat.Scale(sprite.Scale[0], sprite.Scale[1], 1.0).Translate(float32(t.X), float32(t.Y + t.Height), 0)

// 		// TODO - I think there's some mistakes here with premultiplied vs non premultiplied alpha
// 		col := glitch.RGBA{float32(sprite.Color.R)/255.0, float32(sprite.Color.G)/255.0, float32(sprite.Color.B)/255.0, float32(sprite.Color.A)/255.0}
// 		pass.SetLayer(sprite.Layer)
// 		sprite.DrawColorMask(pass, mat, col)
// 	})
// }

// func DrawMultiSprites(pass *glitch.RenderPass, world *ecs.World) {
// 	ecs.Map2(world, func(id ecs.Id, mSprite *MultiSprite, t *phy2.Transform) {
// 		for _, sprite := range mSprite.Sprites {
// 			mat := glitch.Mat4Ident
// 			mat.Scale(sprite.Scale[0], sprite.Scale[1], 1.0).Translate(float32(t.X), float32(t.Y + t.Height), 0)

// 			// TODO - I think there's some mistakes here with premultiplied vs non premultiplied alpha
// 			col := glitch.RGBA{float32(sprite.Color.R)/255.0, float32(sprite.Color.G)/255.0, float32(sprite.Color.B)/255.0, float32(sprite.Color.A)/255.0}
// 			pass.SetLayer(sprite.Layer)
// 			sprite.DrawColorMask(pass, mat, col)
// 		}
// 	})
// }

// func CaptureInput(win *glitch.Window, world *ecs.World) {
// 	// TODO - technically this should only run for the player Ids?
// 	ecs.Map2(world, func(id ecs.Id, keybinds *Keybinds, input *phy2.Input) {
// 		input.Left = false
// 		input.Right = false
// 		input.Up = false
// 		input.Down = false

// 		if win.Pressed(keybinds.Left) {
// 			input.Left = true
// 		}
// 		if win.Pressed(keybinds.Right) {
// 			input.Right = true
// 		}
// 		if win.Pressed(keybinds.Up) {
// 			input.Up = true
// 		}
// 		if win.Pressed(keybinds.Down) {
// 			input.Down = true
// 		}
// 	})
// }
