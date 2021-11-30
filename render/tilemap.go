package render

import (
	// "github.com/faiface/pixel"
	// "github.com/faiface/pixel/pixelgl"

	"github.com/jstewart7/glitch"

	"github.com/jstewart7/flow/tilemap"
	"github.com/jstewart7/flow/asset"
)

type TilemapRender struct {
	spritesheet *asset.Spritesheet
	pass *glitch.RenderPass
	tileToSprite map[tilemap.TileType]*glitch.Sprite
}

func NewTilemapRender(spritesheet *asset.Spritesheet, tileToSprite map[tilemap.TileType]*glitch.Sprite, pass *glitch.RenderPass) *TilemapRender {
	// Note: Assumes that all sprites share the same spritesheet
	return &TilemapRender{
		spritesheet: spritesheet,
		pass: pass,
		tileToSprite: tileToSprite,
	}
}

func (r *TilemapRender) Clear() {
	r.pass.Clear()
}

func (r *TilemapRender) Batch(t *tilemap.Tilemap) {
	for x := 0; x < t.Width(); x++ {
		for y := 0; y < t.Height(); y++ {
			tile, ok := t.Get(x, y)
			if !ok { continue }

			pos := glitch.Vec2{float32(x * t.TileSize), float32(y * t.TileSize)}

			sprite, ok := r.tileToSprite[tile.Type]
			if !ok {
				panic("Unable to find TileType")
			}

			mat := glitch.Mat4Ident
			mat.Translate(pos[0], pos[1], 0)
			sprite.Draw(r.pass, mat)

			// mat := pixel.IM.Moved(pos)
			// sprite.Draw(r.batch, mat)
		}
	}
}

func (r *TilemapRender) Draw(win *glitch.Window) {
	r.pass.Draw(win)
}
