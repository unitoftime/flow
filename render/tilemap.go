package render

import (
	"github.com/unitoftime/glitch"

	"github.com/unitoftime/flow/tilemap"
	"github.com/unitoftime/flow/asset"
)

type TilemapRender struct {
	spritesheet *asset.Spritesheet
	pass *glitch.RenderPass
	tileToSprite map[tilemap.TileType]*glitch.Sprite
}

func NewTilemapRender(spritesheet *asset.Spritesheet,
	tileToSprite map[tilemap.TileType]*glitch.Sprite,
	pass *glitch.RenderPass) *TilemapRender {
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

			// pos := r.Math.Position(x, y, t.TileSize)
			xPos, yPos := t.TileToPosition(tilemap.TilePosition{x, y})
			pos := glitch.Vec2{xPos, yPos}

			pos[1] += tile.Height * float32(t.TileSize[1])

			// Normal grid
			// pos := glitch.Vec2{float32(x * t.TileSize[0]), float32(y * t.TileSize[1])}

			// Isometric grid
			// pos := glitch.Vec2{
			// 	// If y goes up, then xPos must go downward a bit
			// 	-float32((x * t.TileSize[0] / 2) - (y * t.TileSize[0] / 2)),
			// 	// If x goes up, then yPos must go up a bit as well
			// 	-float32((y * t.TileSize[1] / 2) + (x * t.TileSize[1] / 2))}

			// fmt.Println(pos)
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

// func (r *TilemapRender) Draw(win *glitch.Window) {
// 	r.pass.Draw(win)
// }
