package tile

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

func ManhattanDistance(a, b TilePosition) int {
	dx := a.X - b.X
	dy := a.Y - b.Y
	if dx < 0 { dx = -dx }
	if dy < 0 { dy = -dy }
	return dx + dy
}

type Rect struct {
	Min, Max TilePosition
}
func R(minX, minY, maxX, maxY int) Rect {
	return Rect{
		TilePosition{minX, minY},
		TilePosition{maxX, maxY},
	}
}
func (r Rect) Contains(pos TilePosition) bool {
	return pos.X < r.Max.X && pos.X > r.Min.X && pos.Y < r.Max.Y && pos.Y > r.Min.Y
}


type Collider struct {
	Width, Height int // Size of the collider in terms of tiles
}

type Tilemap struct {
	TileSize [2]int // In pixels
	tiles [][]Tile
	math Math
	Offset physics.Vec2 // In world space positioning
}

func New(tiles [][]Tile, tileSize [2]int, math Math) *Tilemap {
	return &Tilemap{
		TileSize: tileSize,
		tiles: tiles,
		math: math,
		Offset: physics.Vec2{},
	}
}

func (t *Tilemap) Width() int {
	return len(t.tiles)
}

func (t *Tilemap) Height() int {
	// TODO - Assumes the tilemap is a square and is larger than size 0
	return len(t.tiles[0])
}

func (t *Tilemap) Get(pos TilePosition) (Tile, bool) {
	if pos.X < 0 || pos.X >= len(t.tiles) || pos.Y < 0 || pos.Y >= len(t.tiles[pos.X]) {
		return Tile{}, false
	}

	return t.tiles[pos.X][pos.Y], true
}

func (t *Tilemap) TileToPosition(tilePos TilePosition) (float32, float32) {
	x, y := t.math.Position(tilePos.X, tilePos.Y, t.TileSize)
	return (x + float32(t.Offset.X)), (y + float32(t.Offset.Y))
}

func (t *Tilemap) PositionToTile(x, y float32) TilePosition {
	x -= float32(t.Offset.X)
	y -= float32(t.Offset.Y)
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
