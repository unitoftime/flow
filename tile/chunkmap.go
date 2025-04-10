package tile

import (
	// "fmt"

	"iter"

	"github.com/unitoftime/flow/glm"
	"github.com/unitoftime/flow/phy2"
	"github.com/unitoftime/intmap"
	"github.com/zyedidia/generic/queue"
)

// type ChunkLoader[T any] interface {
// 	LoadChunk(chunkPos ChunkPosition) ([][]T, error)
// 	SaveChunk(chunkmap *Chunkmap[T], chunkPos ChunkPosition) error
// }

type ChunkPosition struct {
	X, Y int16
}

func (c *ChunkPosition) hash() uint32 {
	return (uint32(uint16(c.X)) << 16) | uint32(uint16(c.Y))
}
func fromHash(hash uint32) ChunkPosition {
	return ChunkPosition{
		X: int16(uint16((hash >> 16) & 0xFFFF)),
		Y: int16(uint16(hash & 0xFFFF)),
	}
}

type Chunkmap[T any] struct {
	ChunkMath
	// chunks map[ChunkPosition]*Chunk[T]
	chunks *intmap.Map[uint32, *Chunk[T]]
}

func NewChunkmap[T any](math ChunkMath) *Chunkmap[T] {
	return &Chunkmap[T]{
		ChunkMath: math,
		// chunks: make(map[ChunkPosition]*Chunk[T]),
		chunks: intmap.New[uint32, *Chunk[T]](0),
	}
}

// func(c *Chunkmap[T]) SetLoader(loader ChunkLoader[T]) *Chunkmap[T] {
// 	c.loader = loader
// 	return c
// }

// TODO - It might be cool to have a function which returns a rectangle of chunks as a list (To automatically cull out sections we don't want)
func (c *Chunkmap[T]) GetAllChunks() []*Chunk[T] {
	ret := make([]*Chunk[T], 0, c.NumChunks())
	// for _, chunk := range c.chunks {
	c.chunks.ForEach(func(_ uint32, chunk *Chunk[T]) {
		ret = append(ret, chunk)
	})
	return ret
}
func (c *Chunkmap[T]) GetAllChunkPositions() []ChunkPosition {
	ret := make([]ChunkPosition, 0, c.NumChunks())
	c.chunks.ForEach(func(chunkHash uint32, chunk *Chunk[T]) {
		chunkPos := fromHash(chunkHash)
		// for chunkPos := range c.chunks {
		ret = append(ret, chunkPos)
	})
	return ret
}

func (c *Chunkmap[T]) Bounds() Rect {
	var bounds Rect
	i := 0
	c.chunks.ForEach(func(chunkHash uint32, chunk *Chunk[T]) {
		chunkPos := fromHash(chunkHash)
		// for chunkPos := range c.chunks {
		if i == 0 {
			bounds = c.GetChunkTileRect(chunkPos)
		} else {
			bounds = bounds.Union(c.GetChunkTileRect(chunkPos))
		}
		i++
	})
	return bounds
	// var bounds Rect
	// i := 0
	// for _, chunk := range c.chunks {
	// 	if i == 0 {
	// 		bounds = chunk.Bounds()
	// 	} else {
	// 		bounds = bounds.Union(chunk.Bounds())
	// 	}
	// 	i++
	// }
	// return bounds
}

func (c *Chunkmap[T]) GetChunk(chunkPos ChunkPosition) (*Chunk[T], bool) {
	chunk, ok := c.chunks.Get(chunkPos.hash())
	if ok {
		return chunk, true
	}

	return nil, false
}

// This generates a chunk based on the passed in expansionLambda
func (c *Chunkmap[T]) GenerateChunk(chunkPos ChunkPosition, expansionLambda func(x, y int) T) *Chunk[T] {
	chunk, ok := c.GetChunk(chunkPos)
	if ok {
		return chunk // Return the chunk and don't create, if the chunk is already made
	}

	tileOffset := c.ChunkToTile(chunkPos)

	tiles := make([][]T, c.ChunkMath.chunkmath.size[0], c.ChunkMath.chunkmath.size[1])
	for x := range tiles {
		tiles[x] = make([]T, c.ChunkMath.chunkmath.size[0], c.ChunkMath.chunkmath.size[1])
		for y := range tiles[x] {
			// fmt.Println(x, y, tileOffset.X, tileOffset.Y)
			if expansionLambda != nil {
				tiles[x][y] = expansionLambda(x+tileOffset.X, y+tileOffset.Y)
			}
		}
	}

	return c.AddChunk(chunkPos, tiles)
}

// func (c *Chunkmap[T]) SaveChunk(chunkPos ChunkPosition) error {
// 	if c.loader == nil {
// 		return fmt.Errorf("Chunkmap loader is nil")
// 	}

// 	// TODO - I feel like I need some way to dump a chunk out of memory. like, SaveCHunk(...), then RemoveFromMemory(...) - OR - PersistChunk (...) which just does both of those

// 	return c.loader.SaveChunk(c, chunkPos)
// }

func (c *Chunkmap[T]) AddChunk(chunkPos ChunkPosition, tiles [][]T) *Chunk[T] {
	chunk := NewChunk[T](tiles, c.ChunkMath.tilemath.size, c.ChunkMath.tilemath)

	// offX, offY := c.ChunkMath.math.Position(int(chunkPos.X), int(chunkPos.Y),
	// 	[2]int{c.ChunkMath.tileSize[0]*c.ChunkMath.chunkSize[0], c.ChunkMath.tileSize[1]*c.ChunkMath.chunkSize[1]})
	// chunk.Offset.X = float64(offX)
	// chunk.Offset.Y = float64(offY)
	chunk.Offset = c.ChunkMath.ToPosition(chunkPos)

	chunk.TileOffset = c.ChunkToTile(chunkPos)

	// Write back
	// c.chunks[chunkPos] = chunk
	c.chunks.Put(chunkPos.hash(), chunk)
	return chunk
}

func (c *Chunkmap[T]) NumChunks() int {
	return c.chunks.Len()
	// return len(c.chunks)
}

func (c *Chunkmap[T]) GetTile(pos TilePosition) (T, bool) {
	chunkPos := c.TileToChunk(pos)
	chunk, ok := c.GetChunk(chunkPos)
	if !ok {
		var ret T
		return ret, false
	}

	localTilePos := TilePosition{pos.X - chunk.TileOffset.X, pos.Y - chunk.TileOffset.Y}
	return chunk.unsafeGet(localTilePos), true
}

// Adds a tile at the position, if the chunk doesnt exist, then it will be created
func (c *Chunkmap[T]) AddTile(pos TilePosition, tile T) {
	success := c.SetTile(pos, tile)
	if !success {
		// TODO: in my game, I have the default T value being an empty tile, but it might be nice to modify Chunkmap struct to have a 'defaultTile' or something that people can use to fill blank spaces
		chunkPos := c.TileToChunk(pos)
		c.GenerateChunk(chunkPos, nil)
		success := c.SetTile(pos, tile)
		if !success {
			panic("programmer error")
		}
	}
}

// Tries to set the tile, returns false if the chunk does not exist
func (c *Chunkmap[T]) SetTile(pos TilePosition, tile T) bool {
	chunkPos := c.TileToChunk(pos)
	chunk, ok := c.GetChunk(chunkPos)
	if !ok {
		return false
	}

	tileOffset := c.ChunkToTile(chunkPos)
	localTilePos := TilePosition{pos.X - tileOffset.X, pos.Y - tileOffset.Y}
	// fmt.Println("chunk.Get:", chunkPos, pos, localTilePos)
	ok = chunk.Set(localTilePos, tile)
	if !ok {
		panic("Programmer error")
	}
	return true
}

func (c *Chunkmap[T]) GetEightNeighbors(pos TilePosition) (T, T, T, T, T, T, T, T) {
	t, _ := c.GetTile(Position{pos.X, pos.Y + 1})
	b, _ := c.GetTile(Position{pos.X, pos.Y - 1})
	l, _ := c.GetTile(Position{pos.X - 1, pos.Y})
	r, _ := c.GetTile(Position{pos.X + 1, pos.Y})
	tl, _ := c.GetTile(Position{pos.X - 1, pos.Y + 1})
	tr, _ := c.GetTile(Position{pos.X + 1, pos.Y + 1})
	bl, _ := c.GetTile(Position{pos.X - 1, pos.Y - 1})
	br, _ := c.GetTile(Position{pos.X + 1, pos.Y - 1})

	return t, b, l, r, tl, tr, bl, br
}

func (c *Chunkmap[T]) GetNeighborsAtDistance(tilePos TilePosition, dist int) []TilePosition {
	distance := make(map[TilePosition]int)

	q := queue.New[TilePosition]()
	q.Enqueue(tilePos)

	for !q.Empty() {
		current := q.Dequeue()

		d := distance[current]
		if d >= dist {
			continue
		} // Don't need to search past our limit

		neighbors := c.GetEdgeNeighbors(current.X, current.Y)
		for _, next := range neighbors {
			_, ok := c.GetTile(next)
			if !ok {
				continue
			} // Skip as neighbor doesn't actually exist (ie could be OOB)

			// If we haven't already walked over this neighbor, then enqueue it and add it to our path
			_, exists := distance[next]
			if !exists {
				q.Enqueue(next)
				distance[next] = 1 + distance[current]
			}
		}
	}

	// Pull out all of the tiles that are at the correct distance
	ret := make([]TilePosition, 0)
	for pos, d := range distance {
		if d != dist {
			continue
		} // Don't return if distance isn't corect
		ret = append(ret, pos)
	}
	return ret
}

func (c *Chunkmap[T]) BreadthFirstSearch(tilePos TilePosition, valid func(t T) bool) []TilePosition {
	distance := make(map[TilePosition]int)

	q := queue.New[TilePosition]()
	q.Enqueue(tilePos)

	for !q.Empty() {
		current := q.Dequeue()

		neighbors := c.GetEdgeNeighbors(current.X, current.Y)
		for _, next := range neighbors {
			t, ok := c.GetTile(next)
			if !ok {
				continue
			} // Skip as neighbor doesn't actually exist (ie could be OOB)

			if !valid(t) {
				continue
			} // Skip if the tile isn't valid

			// If we haven't already walked over this neighbor, then enqueue it and add it to our path
			_, exists := distance[next]
			if !exists {
				q.Enqueue(next)
				distance[next] = 1 + distance[current]
			}
		}
	}

	// Pull out all of the tiles that are at the correct distance
	ret := make([]TilePosition, 0)
	for pos := range distance {
		ret = append(ret, pos)
	}
	return ret
}

func (c *Chunkmap[T]) IterBreadthFirst(tilePos TilePosition, valid func(t T) bool) iter.Seq[Position] {
	return func(yield func(Position) bool) {

		distance := make(map[TilePosition]int)
		q := queue.New[TilePosition]()
		q.Enqueue(tilePos)

		for !q.Empty() {
			current := q.Dequeue()

			if !yield(current) {
				break
			}

			neighbors := c.GetEdgeNeighbors(current.X, current.Y)
			for _, next := range neighbors {
				t, ok := c.GetTile(next)
				if !ok {
					continue
				} // Skip as neighbor doesn't actually exist (ie could be OOB)

				if !valid(t) {
					continue
				} // Skip if the tile isn't valid

				// If we haven't already walked over this neighbor, then enqueue it and add it to our path
				_, exists := distance[next]
				if !exists {
					q.Enqueue(next)
					distance[next] = 1 + distance[current]
				}
			}
		}
	}
}

func (c *Chunkmap[T]) GetPerimeter() map[ChunkPosition]bool {
	perimeter := make(map[ChunkPosition]bool) // List of chunkPositions that are the perimeter
	processed := make(map[ChunkPosition]bool) // List of chunkPositions that we've already processed

	// Just start at some random chunkPosition (whichever is first)
	var start ChunkPosition
	// TODO: originally this function would just use the first in the for loop. but we cant break out of a lambda func
	c.chunks.ForEach(func(chunkHash uint32, chunk *Chunk[T]) {
		chunkPos := fromHash(chunkHash)
		// for chunkPos := range c.chunks {
		start = chunkPos
		// break
	})

	q := queue.New[ChunkPosition]()
	q.Enqueue(start)

	for !q.Empty() {
		current := q.Dequeue()

		neighbors := c.GetChunkEdgeNeighbors(current)
		for _, next := range neighbors {
			_, ok := c.GetChunk(next)
			if ok {
				// If the chunk's neighbor exists, then add it and keep processing
				_, alreadyProcessed := processed[next]
				if !alreadyProcessed {
					q.Enqueue(next)
					processed[next] = true
				}
				continue
			}

			perimeter[next] = true
		}
	}

	return perimeter
}

func (c *Chunkmap[T]) CalculateBlobmapVariant(pos TilePosition, same func(a T, b T) bool) uint8 {
	tile, ok := c.GetTile(pos)
	if !ok {
		return 0
	}

	t, _ := c.GetTile(TilePosition{pos.X, pos.Y + 1})
	b, _ := c.GetTile(TilePosition{pos.X, pos.Y - 1})
	l, _ := c.GetTile(TilePosition{pos.X - 1, pos.Y})
	r, _ := c.GetTile(TilePosition{pos.X + 1, pos.Y})
	tl, _ := c.GetTile(TilePosition{pos.X - 1, pos.Y + 1})
	tr, _ := c.GetTile(TilePosition{pos.X + 1, pos.Y + 1})
	bl, _ := c.GetTile(TilePosition{pos.X - 1, pos.Y - 1})
	br, _ := c.GetTile(TilePosition{pos.X + 1, pos.Y - 1})

	return PackedBlobmapNumber(
		same(tile, t),
		same(tile, b),
		same(tile, l),
		same(tile, r),
		same(tile, tl),
		same(tile, tr),
		same(tile, bl),
		same(tile, br),
	)
}

func (c *Chunkmap[T]) CalculatePipemapVariant(pos TilePosition, same func(a T, b T) bool) uint8 {
	tile, ok := c.GetTile(pos)
	if !ok {
		return 0
	}

	t, _ := c.GetTile(TilePosition{pos.X, pos.Y + 1})
	b, _ := c.GetTile(TilePosition{pos.X, pos.Y - 1})
	l, _ := c.GetTile(TilePosition{pos.X - 1, pos.Y})
	r, _ := c.GetTile(TilePosition{pos.X + 1, pos.Y})

	return PackedPipemapNumber(
		same(tile, t),
		same(tile, b),
		same(tile, l),
		same(tile, r),
	)
}

func (c *Chunkmap[T]) CalculateRawEightVariant(pos TilePosition, same func(a T, b T) bool) uint8 {
	tile, ok := c.GetTile(pos)
	if !ok {
		return 0
	}

	t, _ := c.GetTile(TilePosition{pos.X, pos.Y + 1})
	b, _ := c.GetTile(TilePosition{pos.X, pos.Y - 1})
	l, _ := c.GetTile(TilePosition{pos.X - 1, pos.Y})
	r, _ := c.GetTile(TilePosition{pos.X + 1, pos.Y})
	tl, _ := c.GetTile(TilePosition{pos.X - 1, pos.Y + 1})
	tr, _ := c.GetTile(TilePosition{pos.X + 1, pos.Y + 1})
	bl, _ := c.GetTile(TilePosition{pos.X - 1, pos.Y - 1})
	br, _ := c.GetTile(TilePosition{pos.X + 1, pos.Y - 1})

	return PackedRawEightNumber(
		same(tile, t),
		same(tile, b),
		same(tile, l),
		same(tile, r),
		same(tile, tl),
		same(tile, tr),
		same(tile, bl),
		same(tile, br),
	)
}

// --------------------------------------------------------------------------------
// - Math functions
// --------------------------------------------------------------------------------

type ChunkMath struct {
	// chunkSize [2]int
	// chunkSizeOver2 [2]int
	// chunkDiv [2]int
	// tileSize [2]int
	// tileSizeOver2 [2]int
	// tileDiv [2]int
	// math FlatRectMath

	globalmath FlatRectMath
	chunkmath  FlatRectMath
	tilemath   FlatRectMath
}

func NewChunkmath(chunkSize int, tileSize int) ChunkMath {
	return ChunkMath{
		globalmath: NewFlatRectMath([2]int{tileSize * chunkSize, tileSize * chunkSize}),
		chunkmath:  NewFlatRectMath([2]int{chunkSize, chunkSize}),
		tilemath:   NewFlatRectMath([2]int{tileSize, tileSize}),
	}
	// chunkDiv := int(math.Log2(float64(chunkSize)))
	// if (1 << chunkDiv) != chunkSize {
	// 	panic("Chunk maps must have a chunksize that is a power of 2!")
	// }

	// tileDiv := int(math.Log2(float64(tileSize)))
	// if (1 << tileDiv) != tileSize {
	// 	panic("Chunk maps must have a tileSize that is a power of 2!")
	// }

	// return ChunkMath{
	// 	chunkSize: [2]int{chunkSize, chunkSize},
	// 	chunkDiv: [2]int{chunkDiv, chunkDiv},
	// 	chunkSizeOver2: [2]int{chunkSize/2, chunkSize/2},
	// 	tileSize: [2]int{tileSize, tileSize},
	// 	tileSizeOver2: [2]int{tileSize/2, tileSize/2},
	// 	tileDiv: [2]int{tileDiv, tileDiv},
	// 	math: FlatRectMath{},
	// }
}

// Returns the worldspace position of a chunk
func (c *ChunkMath) ToPosition(chunkPos ChunkPosition) glm.Vec2 {
	offX, offY := c.globalmath.Position(int(chunkPos.X), int(chunkPos.Y))
	// offX, offY := c.math.Position(int(chunkPos.X), int(chunkPos.Y),
	// 	[2]int{c.tileSize[0]*c.chunkSize[0], c.tileSize[1]*c.chunkSize[1]})

	offset := glm.Vec2{
		X: float64(offX),
		// Y: float64(offY) - (0.5 * float64(c.chunkSize[1]) * float64(c.tileSize[1])) + float64(c.tileSize[1]/2),
		// Y: float64(offY) - float64(c.chunkSizeOver2[1] * c.tileSize[1]) + float64(c.tileSizeOver2[1]),
		Y: float64(offY) - float64(c.chunkmath.sizeOver2[1]*c.tilemath.size[1]) + float64(c.tilemath.sizeOver2[1]),
	}
	return offset
}

// Note: untested
func (c *ChunkMath) TileToChunkLocalPosition(tilePos Position) glm.Vec2 {
	chunkPos := c.TileToChunk(tilePos)
	offsetPos := c.ToPosition(chunkPos)
	pos := c.TileToPosition(tilePos)
	return pos.Sub(offsetPos)
}

func (c *ChunkMath) PositionToChunk(x, y float64) ChunkPosition {
	return c.TileToChunk(c.PositionToTile(x, y))
}

func (c *ChunkMath) TileToChunk(tilePos TilePosition) ChunkPosition {
	xPos := tilePos.X >> c.chunkmath.div[0]
	yPos := tilePos.Y >> c.chunkmath.div[1]

	return ChunkPosition{int16(xPos), int16(yPos)}

	// xPos := (int(tilePos.X) + c.chunkmath.sizeOver2[0]) >> c.chunkmath.div[0]
	// yPos := (int(tilePos.Y) + c.chunkmath.sizeOver2[1]) >> c.chunkmath.div[1]

	// // Adjust for negatives
	// if tilePos.X < -c.chunkmath.sizeOver2[0] {
	// 	xPos -= 1
	// }
	// if tilePos.Y < -c.chunkmath.sizeOver2[1] {
	// 	yPos -= 1
	// }
	// return ChunkPosition{int16(xPos), int16(yPos)}

	// if tilePos.X < 0 {
	// 	tilePos.X -= (c.chunkmath.size[0] - 1)
	// }
	// if tilePos.Y < 0 {
	// 	tilePos.Y -= (c.chunkmath.size[1] - 1)
	// }
	// chunkX := tilePos.X / c.chunkmath.size[0]
	// chunkY := tilePos.Y / c.chunkmath.size[1]
	// // chunkX := tilePos.X >> c.chunkmath.div[0]
	// // chunkY := tilePos.Y >> c.chunkmath.div[1]

	// return ChunkPosition{int16(chunkX), int16(chunkY)}
}

// Returns the center tile of a chunk
func (c *ChunkMath) ChunkToTile(chunkPos ChunkPosition) TilePosition {
	tileX := int(chunkPos.X) * c.chunkmath.size[0]
	tileY := int(chunkPos.Y) * c.chunkmath.size[1]

	tilePos := TilePosition{tileX, tileY}

	return tilePos
}

func (c *ChunkMath) TileToPosition(tilePos TilePosition) glm.Vec2 {
	x, y := c.tilemath.Position(tilePos.X, tilePos.Y)
	return glm.Vec2{x, y}
}

func (c *ChunkMath) PositionToTile(x, y float64) TilePosition {
	tX, tY := c.tilemath.PositionToTile(x, y)
	return TilePosition{tX, tY}
}
func (c *ChunkMath) PositionToTile2(pos glm.Vec2) TilePosition {
	tX, tY := c.tilemath.PositionToTile(pos.X, pos.Y)
	return TilePosition{tX, tY}
}

func (c *ChunkMath) GetChunkTileRect(chunkPos ChunkPosition) Rect {
	center := c.ChunkToTile(chunkPos)
	return R(
		center.X,
		center.Y,
		center.X+(c.chunkmath.size[0])-1,
		center.Y+(c.chunkmath.size[1])-1,
	)
}

// // Returns the worldspace position of a chunk
// func (c *ChunkMath) ToPosition(chunkPos ChunkPosition) glm.Vec22 {
// 	offX, offY := c.math.Position(int(chunkPos.X), int(chunkPos.Y),
// 		[2]int{c.tileSize[0]*c.chunkSize[0], c.tileSize[1]*c.chunkSize[1]})

// 	offset := glm.Vec22{
// 		X: float64(offX),
// 		// Y: float64(offY) - (0.5 * float64(c.chunkSize[1]) * float64(c.tileSize[1])) + float64(c.tileSize[1]/2),
// 		Y: float64(offY) - float64(c.chunkSizeOver2[1] * c.tileSize[1]) + float64(c.tileSizeOver2[1]),
// 	}
// 	return offset
// }

// //Note: untested
// func (c *ChunkMath) TileToChunkLocalPosition(tilePos Position) glm.Pos {
// 	chunkPos := c.TileToChunk(tilePos)
// 	offsetPos := c.ToPosition(chunkPos)
// 	pos := c.TileToPosition(tilePos)
// 	return pos.Sub(glm.Pos(offsetPos))
// }

// func (c *ChunkMath) PositionToChunk(x, y float64) ChunkPosition {
// 	return c.TileToChunk(c.PositionToTile(x, y))
// }

// func (c *ChunkMath) TileToChunk(tilePos TilePosition) ChunkPosition {
// 	if tilePos.X < 0 {
// 		tilePos.X -= (c.chunkSize[0] - 1)
// 	}
// 	if tilePos.Y < 0 {
// 		tilePos.Y -= (c.chunkSize[1] - 1)
// 	}
// 	chunkX := tilePos.X / c.chunkSize[0]
// 	chunkY := tilePos.Y / c.chunkSize[1]

// 	return ChunkPosition{int16(chunkX), int16(chunkY)}
// }

// // Returns the center tile of a chunk
// func (c *ChunkMath) ChunkToTile(chunkPos ChunkPosition) TilePosition {
// 	tileX := int(chunkPos.X) * c.chunkSize[0]
// 	tileY := int(chunkPos.Y) * c.chunkSize[1]

// 	tilePos := TilePosition{tileX, tileY}

// 	return tilePos
// }

// func (c *ChunkMath) TileToPosition(tilePos TilePosition) glm.Pos {
// 	x, y := c.math.Position(tilePos.X, tilePos.Y, c.tileSize)
// 	return glm.Pos{x, y}
// }

// func (c *ChunkMath) PositionToTile(x, y float64) TilePosition {
// 	tX, tY := c.math.PositionToTile(x, y, c.tileSize)
// 	return TilePosition{tX, tY}
// }
// func (c *ChunkMath) PositionToTile2(pos glm.Vec2) TilePosition {
// 	tX, tY := c.math.PositionToTile(pos.X, pos.Y, c.tileSize)
// 	return TilePosition{tX, tY}
// }

// func (c *ChunkMath) GetChunkTileRect(chunkPos ChunkPosition) Rect {
// 	center := c.ChunkToTile(chunkPos)
// 	return R(
// 		center.X,
// 		center.Y,
// 		center.X + (c.chunkSize[0]) - 1,
// 		center.Y + (c.chunkSize[1]) - 1,
// 	)
// }

func (c *ChunkMath) WorldToTileRect(r glm.Rect) Rect {
	min := c.PositionToTile2(r.Min)
	max := c.PositionToTile2(r.Max)
	return Rect{min, max}
}

// Returns a rect including all of the tiles.
// Centered on edge tiles
func (c *ChunkMath) RectToWorldRect(r Rect) glm.Rect {
	min := c.TileToPosition(r.Min)
	max := c.TileToPosition(r.Max)
	return glm.Rect{glm.Vec2(min), glm.Vec2(max)}
}

func (c *ChunkMath) GetEdgeNeighbors(x, y int) []TilePosition {
	return []TilePosition{
		TilePosition{x + 1, y},
		TilePosition{x - 1, y},
		TilePosition{x, y + 1},
		TilePosition{x, y - 1},
	}
}

func (c *ChunkMath) GetNeighbors(pos TilePosition) []TilePosition {
	x := pos.X
	y := pos.Y
	return []TilePosition{
		// Edges
		TilePosition{x + 1, y},
		TilePosition{x - 1, y},
		TilePosition{x, y + 1},
		TilePosition{x, y - 1},

		// Corners
		TilePosition{x - 1, y - 1},
		TilePosition{x - 1, y + 1},
		TilePosition{x + 1, y - 1},
		TilePosition{x + 1, y + 1},
	}
}

func (c *ChunkMath) GetChunkEdgeNeighbors(pos ChunkPosition) []ChunkPosition {
	return []ChunkPosition{
		{pos.X + 1, pos.Y},
		{pos.X - 1, pos.Y},
		{pos.X, pos.Y + 1},
		{pos.X, pos.Y - 1},
	}
}

func (c *ChunkMath) GetChunkNeighbors(pos ChunkPosition) []ChunkPosition {
	return []ChunkPosition{
		// Edge
		{pos.X + 1, pos.Y},
		{pos.X - 1, pos.Y},
		{pos.X, pos.Y + 1},
		{pos.X, pos.Y - 1},

		// Corners
		{pos.X + 1, pos.Y + 1},
		{pos.X - 1, pos.Y + 1},
		{pos.X - 1, pos.Y - 1},
		{pos.X + 1, pos.Y - 1},
	}
}

// Returns a list of tiles that are overlapping the collider at a position
func (t *ChunkMath) GetOverlappingTiles(x, y float64, collider *phy2.CircleCollider) []TilePosition {
	minX := x - collider.Radius
	maxX := x + collider.Radius
	minY := y - collider.Radius
	maxY := y + collider.Radius

	min := t.PositionToTile(minX, minY)
	max := t.PositionToTile(maxX, maxY)

	ret := make([]TilePosition, 0)
	for tx := min.X; tx <= max.X; tx++ {
		for ty := min.Y; ty <= max.Y; ty++ {
			ret = append(ret, TilePosition{tx, ty})
		}
	}
	return ret
}

func (t *ChunkMath) GetOverlappingTiles2(ret []TilePosition, x, y float64, radius float64) []TilePosition {
	minX := x - float64(radius)
	maxX := x + float64(radius)
	minY := y - float64(radius)
	maxY := y + float64(radius)

	min := t.PositionToTile(minX, minY)
	max := t.PositionToTile(maxX, maxY)

	ret = ret[:0]
	for tx := min.X; tx <= max.X; tx++ {
		for ty := min.Y; ty <= max.Y; ty++ {
			ret = append(ret, TilePosition{tx, ty})
		}
	}
	return ret
}
