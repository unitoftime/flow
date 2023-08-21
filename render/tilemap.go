package render

import (
	"github.com/unitoftime/glitch"

	"github.com/unitoftime/flow/tile"
)

type TileDraw struct {
	Sprite *glitch.Sprite
	Depth float64
}

type Chunkmap[T any] struct {
	chunkmap *tile.Chunkmap[T]
	tilemapRender *TilemapRender[T]
	chunks map[tile.ChunkPosition]*glitch.Batch
}

func NewChunkmap[T any](chunkmap *tile.Chunkmap[T], tileToSprite func(t T)[]TileDraw) *Chunkmap[T] {
	return &Chunkmap[T]{
		chunkmap: chunkmap,
		tilemapRender: NewTilemapRender[T](tileToSprite),
		chunks: make(map[tile.ChunkPosition]*glitch.Batch),
	}
}

// Returns the request chunk batch
func (c *Chunkmap[T]) GetChunk(chunkPos tile.ChunkPosition) *glitch.Batch {
	batch, ok := c.chunks[chunkPos]
	if !ok {
		c.RebatchChunk(chunkPos)
	}
	batch, ok = c.chunks[chunkPos]
	if !ok {
		panic("Programmer bug")
	}
	return batch
}

// Rebatches a specific chunk (ie signal that the chunk has changed and needs to be rebatched
func (c *Chunkmap[T]) RebatchChunk(chunkPos tile.ChunkPosition) {
	batch, ok := c.chunks[chunkPos]
	if !ok {
		batch = glitch.NewBatch()
	}

	chunk, ok := c.chunkmap.GetChunk(chunkPos)
	if ok {
		// Chunk exists, rebatch it
		batch.Clear()
		c.tilemapRender.Draw(chunk, batch)
	} else {
		// Chunk doesn't exist, so just store an empty batch there
		batch.Clear()
	}
	c.chunks[chunkPos] = batch
}

type TilemapRender[T any] struct {
	tileToSprite func(t T)[]TileDraw
	// tileToSprite map[tile.TileType]*glitch.Sprite
}

func NewTilemapRender[T any](tileToSprite func(t T)[]TileDraw) *TilemapRender[T] {
// func NewTilemapRender(tileToSprite map[tile.TileType]*glitch.Sprite) *TilemapRender {
	// Note: Assumes that all sprites share the same spritesheet
	return &TilemapRender[T]{
		tileToSprite: tileToSprite,
	}
}

func (r *TilemapRender[T]) Draw(tmap *tile.Chunk[T], batch *glitch.Batch) {
	for x := 0; x < tmap.Width(); x++ {
		for y := tmap.Height(); y >= 0; y-- {
			t, ok := tmap.Get(tile.TilePosition{x, y})
			if !ok { continue }

			// pos := r.Math.Position(x, y, t.TileSize)
			xPos, yPos := tmap.TileToPosition(tile.TilePosition{x, y})
			pos := glitch.Vec2{xPos, yPos}// .Add(glitch.Vec2{8, 8})

			// TODO!!! - This should get captured in maybe some extra offset function?
			// pos[1] += t.Height * float32(tmap.TileSize[1])

			// Normal grid
			// pos := glitch.Vec2{float32(x * t.TileSize[0]), float32(y * t.TileSize[1])}

			// Isometric grid
			// pos := glitch.Vec2{
			// 	// If y goes up, then xPos must go downward a bit
			// 	-float32((x * t.TileSize[0] / 2) - (y * t.TileSize[0] / 2)),
			// 	// If x goes up, then yPos must go up a bit as well
			// 	-float32((y * t.TileSize[1] / 2) + (x * t.TileSize[1] / 2))}

			tileDraw := r.tileToSprite(t)
			for _, d := range tileDraw {
				if d.Sprite == nil {
					continue // Skip if the sprite is nil
				}


				mat := glitch.Mat4Ident
				mat.Translate(pos[0], pos[1], d.Depth)
				d.Sprite.Draw(batch, mat)
			}
		}
	}
}
