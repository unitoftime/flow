package tile

import (
	"github.com/unitoftime/flow/glm"
	"github.com/unitoftime/flow/phy2"
)

// TODO
// tile.Chunk?
// tile.Chunkmap?
// chunk.Chunk?
// chunk.Map?

type Chunk[T any] struct {
	TileSize   [2]int // In pixels
	tiles      [][]T
	math       FlatRectMath
	Offset     glm.Vec2 // In world space positioning
	TileOffset Position
}

func NewChunk[T any](tiles [][]T, tileSize [2]int, math FlatRectMath) *Chunk[T] {
	return &Chunk[T]{
		TileSize: tileSize,
		tiles:    tiles,
		math:     math,
		Offset:   glm.Vec2{},
	}
}

// This returns the underlying array, not a copy
// TODO - should I just make tiles public?
func (t *Chunk[T]) Tiles() [][]T {
	return t.tiles
}

func (t *Chunk[T]) Width() int {
	return len(t.tiles)
}

func (t *Chunk[T]) Height() int {
	// TODO - Assumes the tilemap is a square and is larger than size 0
	return len(t.tiles[0])
}

// func (t *Chunk[T]) Bounds() Rect {
// 	min := t.PositionToTile(0, 0)
// 	return R(
// 		min.X,
// 		min.Y,
// 		min.X + t.Width(),
// 		min.Y + t.Height(),
// 	)
// }

func (t *Chunk[T]) Get(pos TilePosition) (T, bool) {
	if pos.X < 0 || pos.X >= len(t.tiles) || pos.Y < 0 || pos.Y >= len(t.tiles[pos.X]) {
		var ret T
		return ret, false
	}

	return t.tiles[pos.X][pos.Y], true
}

func (t *Chunk[T]) unsafeGet(pos TilePosition) T {
	return t.tiles[pos.X][pos.Y]
}

func (t *Chunk[T]) Set(pos TilePosition, tile T) bool {
	if pos.X < 0 || pos.X >= len(t.tiles) || pos.Y < 0 || pos.Y >= len(t.tiles[pos.X]) {
		return false
	}

	t.tiles[pos.X][pos.Y] = tile
	return true
}

func (t *Chunk[T]) TileToPosition(tilePos TilePosition) (float64, float64) {
	x, y := t.math.Position(tilePos.X, tilePos.Y)
	return (x + t.Offset.X), (y + t.Offset.Y)
}

func (t *Chunk[T]) PositionToTile(x, y float64) TilePosition {
	x -= t.Offset.X
	y -= t.Offset.Y
	tX, tY := t.math.PositionToTile(x, y)
	return TilePosition{tX, tY}
}

func (t *Chunk[T]) GetEdgeNeighbors(x, y int) []TilePosition {
	return []TilePosition{
		TilePosition{x + 1, y},
		TilePosition{x - 1, y},
		TilePosition{x, y + 1},
		TilePosition{x, y - 1},
	}
}

// TODO - this might not work for pointy-top tilemaps
func (t *Chunk[T]) BoundsAt(pos TilePosition) (float64, float64, float64, float64) {
	x, y := t.TileToPosition(pos)
	return float64(x) - float64(t.TileSize[0]/2), float64(y) - float64(t.TileSize[1]/2), float64(x) + float64(t.TileSize[0]/2), float64(y) + float64(t.TileSize[1]/2)
}

// Adds an entity to the chunk
// func (t *Chunk[T]) AddEntity(id ecs.Id, collider *Collider, pos *glm.Pos) {
// 	tilePos := t.PositionToTile(float32(pos.X), float32(pos.Y))

// 	for x := tilePos.X; x < tilePos.X + collider.Width; x++ {
// 		for y := tilePos.Y; y < tilePos.Y + collider.Height; y++ {
// 			// TODO - Just using this as a bounds check
// 			tile, ok := t.Get(TilePosition{x, y})
// 			if ok {
// 				t.tiles[x][y].Entity = id // Store the entity
// 			}
// 		}
// 	}
// }

// func (t *Chunk[T]) ClearEntities() {
// 	// Clear Entities
// 	for x := range t.tiles {
// 		for y := range t.tiles[x] {
// 			t.tiles[x][y].Entity = ecs.InvalidEntity
// 		}
// 	}
// }

// // Recalculates all of the entities that are on tiles based on tile colliders
// func (t *Chunk[T]) RecalculateEntities(world *ecs.World) {
// 	t.ClearEntities()

// 	// Recompute all entities with TileColliders
// 	ecs.Map2(world, func(id ecs.Id, collider *Collider, pos *glm.Pos) {
// 		tilePos := t.PositionToTile(float32(pos.X), float32(pos.Y))

// 		for x := tilePos.X; x < tilePos.X + collider.Width; x++ {
// 			for y := tilePos.Y; y < tilePos.Y + collider.Height; y++ {
// 				t.tiles[x][y].Entity = id // Store the entity
// 			}
// 		}
// 	})
// }

// Returns a list of tiles that are overlapping the collider at a position
func (t *Chunk[T]) GetOverlappingTiles(x, y float64, collider *phy2.CircleCollider) []TilePosition {
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
