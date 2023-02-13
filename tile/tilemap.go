package tile

import (
	"github.com/unitoftime/ecs"
	"github.com/unitoftime/flow/phy2"
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

// TODO! - Replace with actual iterator pattern
func (r Rect) Iter() []TilePosition {
	ret := make([]TilePosition, 0)
	for x := r.Min.X; x <= r.Max.X; x++ {
		for y := r.Min.Y; y <= r.Max.Y; y++ {
			ret = append(ret, TilePosition{x, y})
		}
	}
	return ret
}

type Collider struct {
	Width, Height int // Size of the collider in terms of tiles
}

type Tilemap struct {
	TileSize [2]int // In pixels
	tiles [][]Tile
	math Math
	Offset phy2.Vec2 // In world space positioning
}

func New(tiles [][]Tile, tileSize [2]int, math Math) *Tilemap {
	return &Tilemap{
		TileSize: tileSize,
		tiles: tiles,
		math: math,
		Offset: phy2.Vec2{},
	}
}

// This returns the underlying array, not a copy
// TODO - should I just make tiles public?
func (t *Tilemap) Tiles() [][]Tile {
	return t.tiles
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

func (t *Tilemap) Set(pos TilePosition, tile Tile) bool {
	if pos.X < 0 || pos.X >= len(t.tiles) || pos.Y < 0 || pos.Y >= len(t.tiles[pos.X]) {
		return false
	}

	t.tiles[pos.X][pos.Y] = tile
	return true
}

func (t *Tilemap) TileToPosition(tilePos TilePosition) (float64, float64) {
	x, y := t.math.Position(tilePos.X, tilePos.Y, t.TileSize)
	return (x + float64(t.Offset.X)), (y + float64(t.Offset.Y))
}

func (t *Tilemap) PositionToTile(x, y float64) TilePosition {
	x -= t.Offset.X
	y -= t.Offset.Y
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

// TODO - this might not work for pointy-top tilemaps
func (t *Tilemap) BoundsAt(pos TilePosition) (float64, float64, float64, float64) {
	x, y := t.TileToPosition(pos)
	return float64(x) - float64(t.TileSize[0]/2), float64(y) - float64(t.TileSize[1]/2), float64(x) + float64(t.TileSize[0]/2), float64(y) + float64(t.TileSize[1]/2)
}

func (t *Tilemap) ClearEntities() {
	// Clear Entities
	for x := range t.tiles {
		for y := range t.tiles[x] {
			t.tiles[x][y].Entity = ecs.InvalidEntity
		}
	}
}

// Recalculates all of the entities that are on tiles based on tile colliders
func (t *Tilemap) RecalculateEntities(world *ecs.World) {
	t.ClearEntities()

	query := ecs.Query2[Collider, phy2.Pos](world)
	// Recompute all entities with TileColliders
	query.MapId(func(id ecs.Id, collider *Collider, pos *phy2.Pos) {
		tilePos := t.PositionToTile(pos.X, pos.Y)

		for x := tilePos.X; x < tilePos.X + collider.Width; x++ {
			for y := tilePos.Y; y < tilePos.Y + collider.Height; y++ {
				t.tiles[x][y].Entity = id // Store the entity
			}
		}
	})
}

// Returns a list of tiles that are overlapping the collider at a position
func (t *Tilemap) GetOverlappingTiles(x, y float64, collider *phy2.CircleCollider) []TilePosition {
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
