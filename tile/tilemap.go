package tile

import (
	"iter"

	"github.com/unitoftime/flow/phy2"
)

// type TileType uint8

// type Tile struct {
// 	Type TileType
// 	Height float32
// 	// TODO - Should entity inclusion be held somewhere else? What if two entities occupy the same tile?
// 	Entity ecs.Id // This holds the entity Id of the object that is placed here
// }

type TilePosition = Position

type Position struct {
	X, Y int
}

func (t Position) Add(v Position) Position {
	return Position{
		t.X + v.X,
		t.Y + v.Y,
	}
}

func (t Position) Sub(v Position) Position {
	return Position{
		t.X - v.X,
		t.Y - v.Y,
	}
}

func (t Position) Div(val int) Position {
	return Position{
		t.X / val,
		t.Y / val,
	}
}

func (t Position) Maximum() int {
	x := t.X
	y := t.Y
	if x < 0 { x = -x }
	if y < 0 { y = -y }

	if x > y {
		return x
	}
	return y
}

func (t Position) Manhattan() int {
	x := t.X
	y := t.Y
	if x < 0 { x = -x }
	if y < 0 { y = -y }
	return x + y
}

func ManhattanDistance(a, b Position) int {
	dx := a.X - b.X
	dy := a.Y - b.Y
	if dx < 0 { dx = -dx }
	if dy < 0 { dy = -dy }
	return dx + dy
}

type Rect struct {
	Min, Max Position
}
func R(minX, minY, maxX, maxY int) Rect {
	return Rect{
		Position{minX, minY},
		Position{maxX, maxY},
	}
}
func (r Rect) WithCenter(v Position) Rect {
	c := r.Center()
	zRect := r.Moved(Position{
		X: -c.X,
		Y: -c.Y,
	})
	return zRect.Moved(v)
}

func (r Rect) Area() int {
	return r.W() * r.H()
}

func (r Rect) W() int {
	return r.Max.X - r.Min.X
}
func (r Rect) H() int {
	return r.Max.Y - r.Min.Y
}

func (r Rect) Center() Position {
	return Position{r.Min.X + (r.W()/2), r.Min.Y + (r.H()/2)}
}

func (r Rect) Moved(v Position) Rect {
	return Rect{
		Min: r.Min.Add(v),
		Max: r.Max.Add(v),
	}
}

func (r Rect) Contains(pos Position) bool {
	return pos.X <= r.Max.X && pos.X >= r.Min.X && pos.Y <= r.Max.Y && pos.Y >= r.Min.Y
}

func (r Rect) Intersects(r2 Rect) bool {
	return (
		r.Min.X <= r2.Max.X &&
			r.Max.X >= r2.Min.X &&
			r.Min.Y <= r2.Max.Y &&
			r.Max.Y >= r2.Min.Y)
}

func (r Rect) Norm() Rect {
	x1, x2 := minMax(r.Min.X, r.Max.X)
	y1, y2 := minMax(r.Min.Y, r.Max.Y)
	return R(x1, y1, x2, y2)
}

func (r Rect) Union(s Rect) Rect {
	r = r.Norm()
	s = s.Norm()
	x1, _ := minMax(r.Min.X, s.Min.X)
	_, x2 := minMax(r.Max.X, s.Max.X)
	y1, _ := minMax(r.Min.Y, s.Min.Y)
	_, y2 := minMax(r.Max.Y, s.Max.Y)
	return R(x1, y1, x2, y2)
}

func minMax(a, b int) (int, int) {
	if a > b {
		return b, a
	}
	return a, b
}

func (r Rect) PadAll(pad int) Rect {
	return R(r.Min.X - pad, r.Min.Y - pad, r.Max.X + pad, r.Max.Y + pad)
}

func (r Rect) UnpadAll(pad int) Rect {
	return r.PadAll(-pad)
}

func (r Rect) Iter() iter.Seq[Position] {
	return func(yield func(Position) bool) {

		for x := r.Min.X; x <= r.Max.X; x++ {
			for y := r.Min.Y; y <= r.Max.Y; y++ {
				if !yield(Position{x, y}) {
					return // Exit the iteration
				}
			}
		}

	}
}


// //cod:struct
// type Collider struct {
// 	Width, Height int // Size of the collider in terms of tiles
// }

type Tilemap[T any] struct {
	TileSize [2]int // In pixels
	tiles [][]T
	math Math
	Offset phy2.Vec2 // In world space positioning
}

func New[T any](tiles [][]T, tileSize [2]int, math Math) *Tilemap[T] {
	return &Tilemap[T]{
		TileSize: tileSize,
		tiles: tiles,
		math: math,
		Offset: phy2.Vec2{},
	}
}

// This returns the underlying array, not a copy
// TODO - should I just make tiles public?
func (t *Tilemap[T]) Tiles() [][]T {
	return t.tiles
}

func (t *Tilemap[T]) Width() int {
	return len(t.tiles)
}

func (t *Tilemap[T]) Height() int {
	// TODO - Assumes the tilemap is a square and is larger than size 0
	return len(t.tiles[0])
}

func (t *Tilemap[T]) GetTile(pos Position) (T, bool) {
	if pos.X < 0 || pos.X >= len(t.tiles) || pos.Y < 0 || pos.Y >= len(t.tiles[pos.X]) {
		var ret T
		return ret, false
	}

	return t.tiles[pos.X][pos.Y], true
}

func (t *Tilemap[T]) SetTile(pos Position, tile T) bool {
	if pos.X < 0 || pos.X >= len(t.tiles) || pos.Y < 0 || pos.Y >= len(t.tiles[pos.X]) {
		return false
	}

	t.tiles[pos.X][pos.Y] = tile
	return true
}

func (t *Tilemap[T]) TileToPosition(tilePos Position) (float64, float64) {
	x, y := t.math.Position(tilePos.X, tilePos.Y, t.TileSize)
	return (x + float64(t.Offset.X)), (y + float64(t.Offset.Y))
}

func (t *Tilemap[T]) PositionToTile(x, y float64) Position {
	x -= t.Offset.X
	y -= t.Offset.Y
	tX, tY := t.math.PositionToTile(x, y, t.TileSize)
	return Position{tX, tY}
}

func (t *Tilemap[T]) GetEdgeNeighbors(x, y int) []Position {
	return []Position{
		Position{x+1, y},
		Position{x-1, y},
		Position{x, y+1},
		Position{x, y-1},
	}
}

// TODO - this might not work for pointy-top tilemaps
func (t *Tilemap[T]) BoundsAt(pos Position) (float64, float64, float64, float64) {
	x, y := t.TileToPosition(pos)
	return float64(x) - float64(t.TileSize[0]/2), float64(y) - float64(t.TileSize[1]/2), float64(x) + float64(t.TileSize[0]/2), float64(y) + float64(t.TileSize[1]/2)
}

// Returns a list of tiles that are overlapping the collider at a position
func (t *Tilemap[T]) GetOverlappingTiles(x, y float64, collider *phy2.CircleCollider) []Position {
	minX := x - collider.Radius
	maxX := x + collider.Radius
	minY := y - collider.Radius
	maxY := y + collider.Radius

	min := t.PositionToTile(minX, minY)
	max := t.PositionToTile(maxX, maxY)

	ret := make([]Position, 0)
	for tx := min.X; tx <= max.X; tx++ {
		for ty := min.Y; ty <= max.Y; ty++ {
			ret = append(ret, Position{tx, ty})
		}
	}
	return ret
}
