package tilemap

import (
	"github.com/unitoftime/ecs"
	"github.com/unitoftime/flow/physics"
)

type TileType uint8

type Tile struct {
	Type TileType
	Height float32
	// TODO - Should entity inclusion be held somewhere else? What if two entities occupy the same tile?
	Entity ecs.Id // This holds the entity Id of the object that is placed here
}

type TilePosition struct {
	X, Y int
}

type Collider struct {
	Width, Height int // Size of the collider in terms of tiles
}

type Tilemap struct {
	TileSize [2]int // In pixels
	tiles [][]Tile
	math Math
}

func New(tiles [][]Tile, tileSize [2]int, math Math) *Tilemap {
	return &Tilemap{
		TileSize: tileSize,
		tiles: tiles,
		math: math,
	}
}

func (t *Tilemap) Width() int {
	return len(t.tiles)
}

func (t *Tilemap) Height() int {
	// TODO - Assumes the tilemap is a square and is larger than size 0
	return len(t.tiles[0])
}

// Migrate to use TilePosition
func (t *Tilemap) Get(x, y int) (Tile, bool) {
	if x < 0 || x >= len(t.tiles) || y < 0 || y >= len(t.tiles[x]) {
		return Tile{}, false
	}

	return t.tiles[x][y], true
}

func (t *Tilemap) TileToPosition(tilePos TilePosition) (float32, float32) {
	return t.math.Position(tilePos.X, tilePos.Y, t.TileSize)
}

func (t *Tilemap) PositionToTile(x, y float32) TilePosition {
	tX, tY := t.math.PositionToTile(x, y, t.TileSize)
	return TilePosition{tX, tY}
}

func (t *Tilemap) GetEdgeNeighbors(x, y int) []TilePosition {
	return []TilePosition{
		TilePosition{x+1, y},
		TilePosition{x-1, y},
		TilePosition{x, y+1},
		TilePosition{x, y-1},
	}
}

// Recalculates all of the entities that are on tiles based on tile colliders
func (t *Tilemap) RecalculateEntities(world *ecs.World) {
	// Clear Entities
	for x := range t.tiles {
		for y := range t.tiles[x] {
			t.tiles[x][y].Entity = ecs.InvalidEntity
		}
	}

	// Recompute all entities with TileColliders
	ecs.Map2(world, func(id ecs.Id, collider *Collider, transform *physics.Transform) {
		tilePos := t.PositionToTile(float32(transform.X), float32(transform.Y))

		for x := tilePos.X; x < tilePos.X + collider.Width; x++ {
			for y := tilePos.Y; y < tilePos.Y + collider.Height; y++ {
				t.tiles[x][y].Entity = id // Store the entity
			}
		}
	})
}
