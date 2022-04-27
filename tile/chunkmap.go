package tile

import (
	"encoding/binary"

	"github.com/unitoftime/ecs"
	"github.com/unitoftime/flow/physics"
)

// This is basically a reversable hash of the X and Y position of the chunk
type ChunkId uint64
type ChunkPosition struct {
	X, Y int32
}

// TODO - Investigate: Is this slow? Might need endianness correct in case of sending this data accross hardware (save/load files, etc?)

// Based on the chunk X and Y position, calculate the ChunkId
func ChunkPositionToChunk(pos ChunkPosition) ChunkId {
	// Note: On casting int32 to uint32 (memory layout shouldn't change): https://stackoverflow.com/questions/50815512/when-casting-an-int64-to-uint64-is-the-sign-retained
	var buf [8]byte
	binary.LittleEndian.PutUint32(buf[0:4], uint32(pos.X))
	binary.LittleEndian.PutUint32(buf[4:8], uint32(pos.Y))

	return ChunkId(binary.LittleEndian.Uint64(buf[:]))
}
// Based on the ChunkId, calculate the chunk's X and Y position
func ChunkToPosition(id ChunkId) ChunkPosition {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], uint64(id))

	x := int32(binary.LittleEndian.Uint32(buf[0:4]))
	y := int32(binary.LittleEndian.Uint32(buf[4:8]))
	return ChunkPosition{x, y}
}

type Chunkmap struct {
	ChunkSize int // TODO - implies only square chunks
	TileSize [2]int // In pixels
	math Math
	chunks map[ChunkId]*Tilemap
	expansion func(x, y int) Tile // This is a function that deterministically calculates what each tile should be for each position
}

func NewChunkmap(chunkSize int, tileSize [2]int, math Math, expansionLambda func(x, y int) Tile) *Chunkmap {
	return &Chunkmap{
		ChunkSize: chunkSize,
		TileSize: tileSize,
		math: math,
		chunks: make(map[ChunkId]*Tilemap),
		expansion: expansionLambda,
	}
}

func (c *Chunkmap) TileToChunk(tilePos TilePosition) ChunkId {
	if tilePos.X < 0 {
		tilePos.X -= (c.ChunkSize - 1)
	}
	if tilePos.Y < 0 {
		tilePos.Y -= (c.ChunkSize - 1)
	}
	chunkX := tilePos.X / c.ChunkSize
	chunkY := tilePos.Y / c.ChunkSize

	chunkId := ChunkPositionToChunk(ChunkPosition{int32(chunkX), int32(chunkY)})
	return chunkId
}

// TODO - It might be cool to have a function which returns a rectangle of chunks as a list (To automatically cull out sections we don't want)
func (c *Chunkmap) GetAllChunks() []*Tilemap {
	ret := make([]*Tilemap, 0, c.NumChunks())
	for _, chunk := range c.chunks {
		ret = append(ret, chunk)
	}
	return ret
}

func (c *Chunkmap) GetChunk(chunkId ChunkId) (*Tilemap, bool) {
	chunk, ok := c.chunks[chunkId]
	if !ok {
		return nil, false
	}
	return chunk, true
}

func (c *Chunkmap) CreateChunk(chunkPos ChunkPosition) *Tilemap {
	chunkId := ChunkPositionToChunk(chunkPos)
	tiles := make([][]Tile, c.ChunkSize, c.ChunkSize)
	for x := range tiles {
		tiles[x] = make([]Tile, c.ChunkSize, c.ChunkSize)
		for y := range tiles[x] {
			tiles[x][y] = c.expansion(x, y)
		}
	}
	chunk := New(tiles, c.TileSize, c.math)

	// chunkPos := ChunkToPosition(chunkId)
	offX, offY := c.math.Position(int(chunkPos.X), int(chunkPos.Y),
		[2]int{c.TileSize[0]*c.ChunkSize, c.TileSize[1]*c.ChunkSize})

	chunk.Offset.X = float64(offX)
	chunk.Offset.Y = float64(offY)

	// Write back
	c.chunks[chunkId] = chunk
	return chunk
}


func (c *Chunkmap) NumChunks() int {
	return len(c.chunks)
}

func (c *Chunkmap) GetTile(pos TilePosition) (Tile, bool) {
	chunkId := c.TileToChunk(pos)
	chunk, ok := c.GetChunk(chunkId)
	if !ok {
		return Tile{}, false
	}

	return chunk.Get(pos)
}

func (c *Chunkmap) TileToPosition(tilePos TilePosition) (float32, float32) {
	x, y := c.math.Position(tilePos.X, tilePos.Y, c.TileSize)
	return x, y
}

func (c *Chunkmap) PositionToTile(x, y float32) TilePosition {
	tX, tY := c.math.PositionToTile(x, y, c.TileSize)
	return TilePosition{tX, tY}
}

func (c *Chunkmap) GetEdgeNeighbors(x, y int) []TilePosition {
	return []TilePosition{
		TilePosition{x+1, y},
		TilePosition{x-1, y},
		TilePosition{x, y+1},
		TilePosition{x, y-1},
	}
}

// TODO - optimize in a locality sort of way
// Recalculates all of the entities that are on tiles based on tile colliders
func (c *Chunkmap) RecalculateEntities(world *ecs.World) {
	for _, chunk := range c.chunks {
		// Clear Entities
		for x := range chunk.tiles {
			for y := range chunk.tiles[x] {
				chunk.tiles[x][y].Entity = ecs.InvalidEntity
			}
		}
	}

	// Recompute all entities with TileColliders
	ecs.Map2(world, func(id ecs.Id, collider *Collider, transform *physics.Transform) {
		tilePos := c.PositionToTile(float32(transform.X), float32(transform.Y))

		for x := tilePos.X; x < tilePos.X + collider.Width; x++ {
			for y := tilePos.Y; y < tilePos.Y + collider.Height; y++ {
				chunkId := c.TileToChunk(TilePosition{x, y})
				chunk, ok := c.GetChunk(chunkId)
				if !ok { panic("Something has been built on a chunk that doesn't exist!") }
				localTilePosition := chunk.PositionToTile(float32(transform.X), float32(transform.Y))
				chunk.tiles[localTilePosition.X][localTilePosition.Y].Entity = id // Store the entity
			}
		}
	})
}
