package tile

import (
	"fmt"

	"github.com/unitoftime/flow/phy2"
	"github.com/zyedidia/generic/queue"
)

// type ChunkLoader[T any] interface {
// 	LoadChunk(chunkPos ChunkPosition) ([][]T, error)
// 	SaveChunk(chunkmap *Chunkmap[T], chunkPos ChunkPosition) error
// }

type ChunkPosition struct {
	X, Y int32
}

type Chunkmap[T any] struct {
	ChunkMath
	chunks map[ChunkPosition]*Chunk[T]
	// loader ChunkLoader[T]
}

func NewChunkmap[T any](math ChunkMath) *Chunkmap[T] {
	return &Chunkmap[T]{
		ChunkMath: math,
		chunks: make(map[ChunkPosition]*Chunk[T]),
	}
}

// func(c *Chunkmap[T]) SetLoader(loader ChunkLoader[T]) *Chunkmap[T] {
// 	c.loader = loader
// 	return c
// }

// TODO - It might be cool to have a function which returns a rectangle of chunks as a list (To automatically cull out sections we don't want)
func (c *Chunkmap[T]) GetAllChunks() []*Chunk[T] {
	ret := make([]*Chunk[T], 0, c.NumChunks())
	for _, chunk := range c.chunks {
		ret = append(ret, chunk)
	}
	return ret
}
func (c *Chunkmap[T]) GetAllChunkPositions() []ChunkPosition {
	ret := make([]ChunkPosition, 0, c.NumChunks())
	for chunkPos := range c.chunks {
		ret = append(ret, chunkPos)
	}
	return ret
}

func (c *Chunkmap[T]) Bounds() Rect {
	var bounds Rect
	i := 0
	for chunkPos := range c.chunks {
		fmt.Println(c.GetChunkTileRect(chunkPos))
		if i == 0 {
			bounds = c.GetChunkTileRect(chunkPos)
		} else {
			bounds = bounds.Union(c.GetChunkTileRect(chunkPos))
		}
		i++
	}
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
	chunk, ok := c.chunks[chunkPos]
	if ok {
		return chunk, true
	}

	// // If we couldn't load from map, then load from loader
	// if c.loader != nil {
	// 	tiles, err := c.loader.LoadChunk(chunkPos)
	// 	if err != nil {
	// 		return nil, false
	// 	}

	// 	return c.AddChunk(chunkPos, tiles), true
	// }

	return nil, false
}

// This generates a chunk based on the passed in expansionLambda
func (c *Chunkmap[T]) GenerateChunk(chunkPos ChunkPosition, expansionLambda func(x, y int) T) *Chunk[T] {
	chunk, ok := c.GetChunk(chunkPos)
	if ok {
		return chunk // Return the chunk and don't create, if the chunk is already made
	}

	tileOffset := c.ChunkToTile(chunkPos)

	tiles := make([][]T, c.ChunkMath.ChunkSize[0], c.ChunkMath.ChunkSize[1])
	for x := range tiles {
		tiles[x] = make([]T, c.ChunkMath.ChunkSize[0], c.ChunkMath.ChunkSize[1])
		for y := range tiles[x] {
			// fmt.Println(x, y, tileOffset.X, tileOffset.Y)
			if expansionLambda != nil {
				tiles[x][y] = expansionLambda(x + tileOffset.X, y + tileOffset.Y)
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
	chunk := NewChunk[T](tiles, c.ChunkMath.TileSize, c.ChunkMath.Math)

	offX, offY := c.ChunkMath.Math.Position(int(chunkPos.X), int(chunkPos.Y),
		[2]int{c.ChunkMath.TileSize[0]*c.ChunkMath.ChunkSize[0], c.ChunkMath.TileSize[1]*c.ChunkMath.ChunkSize[1]})

	chunk.Offset.X = float64(offX)
	chunk.Offset.Y = float64(offY)

	// Write back
	c.chunks[chunkPos] = chunk
	return chunk
}

func (c *Chunkmap[T]) NumChunks() int {
	return len(c.chunks)
}

func (c *Chunkmap[T]) GetTile(pos TilePosition) (T, bool) {
	chunkPos := c.TileToChunk(pos)
	chunk, ok := c.GetChunk(chunkPos)
	if !ok {
		var ret T
		return ret, false
	}

	tileOffset := c.ChunkToTile(chunkPos)
	localTilePos := TilePosition{pos.X - tileOffset.X, pos.Y - tileOffset.Y}
	// fmt.Println("chunk.Get:", chunkPos, pos, localTilePos)
	return chunk.Get(localTilePos)
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

func (c *Chunkmap[T]) GetNeighborsAtDistance(tilePos TilePosition, dist int) []TilePosition {
	distance := make(map[TilePosition]int)

	q := queue.New[TilePosition]()
	q.Enqueue(tilePos)

	for !q.Empty() {
		current := q.Dequeue()

		d := distance[current]
		if d >= dist { continue } // Don't need to search past our limit

		neighbors := c.GetEdgeNeighbors(current.X, current.Y)
		for _, next := range neighbors {
			_, ok := c.GetTile(next)
			if !ok { continue } // Skip as neighbor doesn't actually exist (ie could be OOB)

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
		if d != dist { continue } // Don't return if distance isn't corect
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
			if !ok { continue } // Skip as neighbor doesn't actually exist (ie could be OOB)

			if !valid(t) { continue } // Skip if the tile isn't valid

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

func (c *Chunkmap[T]) GetPerimeter() map[ChunkPosition]bool {
	perimeter := make(map[ChunkPosition]bool) // List of chunkPositions that are the perimeter
	processed := make(map[ChunkPosition]bool) // List of chunkPositions that we've already processed

	// Just start at some random chunkPosition (whichever is first
	var start ChunkPosition
	for chunkPos := range c.chunks {
		start = chunkPos
		break
	}

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
	if !ok { return 0 }

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
	if !ok { return 0 }

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
	if !ok { return 0 }

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
	ChunkSize [2]int
	TileSize [2]int
	Math Math
}

// Returns the worldspace position of a chunk
func (c *ChunkMath) ToPosition(chunkPos ChunkPosition) phy2.Vec2 {
	offX, offY := c.Math.Position(int(chunkPos.X), int(chunkPos.Y),
		[2]int{c.TileSize[0]*c.ChunkSize[0], c.TileSize[1]*c.ChunkSize[1]})

	offset := phy2.Vec2{
		float64(offX),
		float64(offY) - (0.5 * float64(c.ChunkSize[1]) * float64(c.TileSize[1])) + float64(c.TileSize[1]/2),
	}
	return offset
}

//Note: untested
func (c *ChunkMath) TileToChunkLocalPosition(tilePos Position) phy2.Pos {
	chunkPos := c.TileToChunk(tilePos)
	offsetPos := c.ToPosition(chunkPos)
	pos := c.TileToPosition(tilePos)
	return pos.Sub(phy2.Pos(offsetPos))
}

func (c *ChunkMath) PositionToChunk(x, y float64) ChunkPosition {
	return c.TileToChunk(c.PositionToTile(x, y))
}

func (c *ChunkMath) TileToChunk(tilePos TilePosition) ChunkPosition {
	if tilePos.X < 0 {
		tilePos.X -= (c.ChunkSize[0] - 1)
	}
	if tilePos.Y < 0 {
		tilePos.Y -= (c.ChunkSize[1] - 1)
	}
	chunkX := tilePos.X / c.ChunkSize[0]
	chunkY := tilePos.Y / c.ChunkSize[1]

	return ChunkPosition{int32(chunkX), int32(chunkY)}
}

// Returns the center tile of a chunk
func (c *ChunkMath) ChunkToTile(chunkPos ChunkPosition) TilePosition {
	tileX := int(chunkPos.X) * c.ChunkSize[0]
	tileY := int(chunkPos.Y) * c.ChunkSize[1]

	tilePos := TilePosition{tileX, tileY}

	return tilePos
}

func (c *ChunkMath) TileToPosition(tilePos TilePosition) phy2.Pos {
	x, y := c.Math.Position(tilePos.X, tilePos.Y, c.TileSize)
	return phy2.Pos{x, y}
}

func (c *ChunkMath) PositionToTile(x, y float64) TilePosition {
	tX, tY := c.Math.PositionToTile(x, y, c.TileSize)
	return TilePosition{tX, tY}
}
func (c *ChunkMath) PositionToTile2(pos phy2.Vec) TilePosition {
	tX, tY := c.Math.PositionToTile(pos.X, pos.Y, c.TileSize)
	return TilePosition{tX, tY}
}

func (c *ChunkMath) GetChunkTileRect(chunkPos ChunkPosition) Rect {
	center := c.ChunkToTile(chunkPos)
	return R(
		center.X,
		center.Y,
		center.X + (c.ChunkSize[0]) - 1,
		center.Y + (c.ChunkSize[1]) - 1,
	)
}

// Returns a rect including all of the tiles.
// Centered on edge tiles
func (c *ChunkMath) RectToWorldRect(r Rect) phy2.Rect {
	min := c.TileToPosition(r.Min)
	max := c.TileToPosition(r.Max)
	return phy2.Rect{phy2.Vec(min), phy2.Vec(max)}
}

func (c *ChunkMath) GetEdgeNeighbors(x, y int) []TilePosition {
	return []TilePosition{
		TilePosition{x+1, y},
		TilePosition{x-1, y},
		TilePosition{x, y+1},
		TilePosition{x, y-1},
	}
}

func (c *ChunkMath) GetChunkEdgeNeighbors(pos ChunkPosition) []ChunkPosition {
	return []ChunkPosition{
		ChunkPosition{pos.X+1, pos.Y},
		ChunkPosition{pos.X-1, pos.Y},
		ChunkPosition{pos.X, pos.Y+1},
		ChunkPosition{pos.X, pos.Y-1},
	}
}

func (c *ChunkMath) GetNeighbors(pos TilePosition) []TilePosition {
	x := pos.X
	y := pos.Y
	return []TilePosition{
		// Edges
		TilePosition{x+1, y},
		TilePosition{x-1, y},
		TilePosition{x, y+1},
		TilePosition{x, y-1},

		// Corners
		TilePosition{x-1, y-1},
		TilePosition{x-1, y+1},
		TilePosition{x+1, y-1},
		TilePosition{x+1, y+1},
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

// // --------------------------------------------------------------------------------
// // - ECS function (TODO pull out)
// // --------------------------------------------------------------------------------

// func (c *Chunkmap[T]) AddEntity(id ecs.Id, pos phy2.Pos, collider Collider) {
// 	tilePos := c.PositionToTile(float32(pos.X), float32(pos.Y))
// 	for x := tilePos.X; x < tilePos.X + collider.Width; x++ {
// 		for y := tilePos.Y; y < tilePos.Y + collider.Height; y++ {
// 			chunkPos := c.TileToChunk(TilePosition{x, y})
// 			// TODO - I feel like adding an entity on a chunk edge shouldn't cause us to reload chunks. This could cause a waterfall effect of chunks triggering more chunks just because entities are placed on edges
// 			chunk, ok := c.GetChunk(chunkPos)
// 			if !ok { continue } // Skip: The chunk doesn't exist
// 			localTilePosition := chunk.PositionToTile(float32(pos.X), float32(pos.Y))
// 			chunk.tiles[localTilePosition.X][localTilePosition.Y].Entity = id // Store the entity
// 		}
// 	}
// }

// // TODO - optimize in a locality sort of way
// // Recalculates all of the entities that are on tiles based on tile colliders
// func (c *Chunkmap[T]) RecalculateEntities(world *ecs.World) {
// 	for _, chunk := range c.chunks {
// 		// Clear Entities
// 		for x := range chunk.tiles {
// 			for y := range chunk.tiles[x] {
// 				chunk.tiles[x][y].Entity = ecs.InvalidEntity
// 			}
// 		}
// 	}

// 	// Recompute all entities with TileColliders
// 	ecs.Map2(world, func(id ecs.Id, collider *Collider, pos *phy2.Pos) {
// 		tilePos := c.PositionToTile(float32(pos.X), float32(pos.Y))

// 		for x := tilePos.X; x < tilePos.X + collider.Width; x++ {
// 			for y := tilePos.Y; y < tilePos.Y + collider.Height; y++ {
// 				chunkPos := c.TileToChunk(TilePosition{x, y})
// 				chunk, ok := c.GetChunk(chunkPos)
// 				if !ok { panic("Something has been built on a chunk that doesn't exist!") }
// 				localTilePosition := chunk.PositionToTile(float32(pos.X), float32(pos.Y))
// 				chunk.tiles[localTilePosition.X][localTilePosition.Y].Entity = id // Store the entity
// 			}
// 		}
// 	})
// }
