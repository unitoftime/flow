package autotile

import (
	"math/rand"

	"github.com/unitoftime/flow/tile"
)

type Pattern uint8
const (
	Top Pattern = 0b00000001
	Right       = 0b00000010
	Bottom      = 0b00000100
	Left        = 0b00001000
)

func (p Pattern) Top() bool {
	return (p & Top) == Top
}
func (p Pattern) Right() bool {
	return (p & Right) == Right
}
func (p Pattern) Bottom() bool {
	return (p & Bottom) == Bottom
}
func (p Pattern) Left() bool {
	return (p & Left) == Left
}

type Tilemap[T any] interface {
	GetTile(pos tile.TilePosition) (T, bool)
}

type Rule[T any] interface {
	Execute(Tilemap[T], tile.TilePosition) int
	ExecuteFull(T, T, T, T, T, T, T, T, T) int
}

// type RawEightRule[T any] struct {
// 	Top []T
// 	Bottom []T
// 	Left []T
// 	Right []T
// }

type BlobmapRule[T any] struct {
	Match func(a, b T) bool
}
func (rule BlobmapRule[T]) ExecuteFull(m, t, b, l, r, tl, tr, bl, br T) int {
	if !rule.Match(m, m) {
		return -1
	}

	pattern := PackedBlobmapNumber(
		rule.Match(m, t),
		rule.Match(m, b),
		rule.Match(m, l),
		rule.Match(m, r),
		rule.Match(m, tl),
		rule.Match(m, tr),
		rule.Match(m, bl),
		rule.Match(m, br),
	)

	return int(pattern)
}

func (rule BlobmapRule[T]) Execute(tilemap Tilemap[T], pos tile.TilePosition) int {
	tile, ok := tilemap.GetTile(pos)
	if !ok {
		return -1
	}
	if !rule.Match(tile, tile) {
		return -1
	}

	t, b, l, r, tl, tr, bl, br := getEightNeighbors(tilemap, pos)


	pattern := PackedBlobmapNumber(
		rule.Match(tile, t),
		rule.Match(tile, b),
		rule.Match(tile, l),
		rule.Match(tile, r),
		rule.Match(tile, tl),
		rule.Match(tile, tr),
		rule.Match(tile, bl),
		rule.Match(tile, br),
	)

	return int(pattern)
}

type PipemapRule[T any] struct {
	Match func(a, b T) bool
}
func (rule PipemapRule[T]) ExecuteFull(m, t, b, l, r, tl, tr, bl, br T) int {
	if !rule.Match(m, m) {
		return -1
	}

	pattern := PackedPipemapNumber(
		rule.Match(m, t),
		rule.Match(m, b),
		rule.Match(m, l),
		rule.Match(m, r),
	)

	return int(pattern)

}
func (rule PipemapRule[T]) Execute(tilemap Tilemap[T], pos tile.TilePosition) int {
	tile, ok := tilemap.GetTile(pos)
	if !ok {
		return -1
	}
	if !rule.Match(tile, tile) {
		return -1
	}

	t, b, l, r := getEdgeNeighbors(tilemap, pos)

	pattern := PackedPipemapNumber(
		rule.Match(tile, t),
		rule.Match(tile, b),
		rule.Match(tile, l),
		rule.Match(tile, r),
	)

	return int(pattern)
}


type LambdaRule[T any] struct {
	Func func(Pattern) int
	Match func(T, T) bool
}
func (rule LambdaRule[T]) ExecuteFull(m, t, b, l, r, tl, tr, bl, br T) int {
	if !rule.Match(m, m) {
		return -1
	}

	pattern := PackedRawEightNumber(
		rule.Match(m, t),
		rule.Match(m, b),
		rule.Match(m, l),
		rule.Match(m, r),
		rule.Match(m, tl),
		rule.Match(m, tr),
		rule.Match(m, bl),
		rule.Match(m, br),
	)

	return rule.Func(Pattern(pattern))
}

func (rule LambdaRule[T]) Execute(tilemap Tilemap[T], pos tile.TilePosition) int {
	tile, ok := tilemap.GetTile(pos)
	if !ok {
		return -1
	}
	if !rule.Match(tile, tile) {
		return -1
	}

	t, b, l, r, tl, tr, bl, br := getEightNeighbors(tilemap, pos)

	pattern := PackedRawEightNumber(
		rule.Match(tile, t),
		rule.Match(tile, b),
		rule.Match(tile, l),
		rule.Match(tile, r),
		rule.Match(tile, tl),
		rule.Match(tile, tr),
		rule.Match(tile, bl),
		rule.Match(tile, br),
	)

	return rule.Func(Pattern(pattern))
}

// type DualGridRule[T any] struct {
// 	Match func(a, b T) bool
// }
// func (rule DualGridRule[T]) Execute(tilemap Tilemap[T], pos tile.TilePosition) int {
// 	tile, ok := tilemap.GetTile(pos)
// 	if !ok {
// 		return -1
// 	}
// 	if !rule.Match(tile, tile) {
// 		return -1
// 	}

// 	t, b, l, r := getEdgeNeighbors(tilemap, pos)

// 	pattern := PackedPipemapNumber(
// 		rule.Match(tile, t),
// 		rule.Match(tile, b),
// 		rule.Match(tile, l),
// 		rule.Match(tile, r),
// 	)

// 	return int(pattern)
// }

type Set[T, S any] struct {
	// mapping map[uint8][]int
	// Rule func(Pattern)int
	Rule Rule[T]
	Tiles [][]S
}

func (s *Set[T,S]) GetFull(m, t, b, l, r, tl, tr, bl, br T) (S, bool) {
	variant := s.Rule.ExecuteFull(m, t, b, l, r, tl, tr, bl, br)
	if variant < 0 {
		var ret S
		return ret, false
	}
	idx := rand.Intn(len(s.Tiles[variant]))
	return s.Tiles[variant][idx], true
}

func (s *Set[T,S]) Get(tilemap Tilemap[T], pos tile.TilePosition) (S, bool) {
	variant := s.Rule.Execute(tilemap, pos)
	if variant < 0 {
		var ret S
		return ret, false
	}
	idx := rand.Intn(len(s.Tiles[variant]))
	return s.Tiles[variant][idx], true
}

// func (s *Set[T]) Get(val Pattern) T {
// 	variant := s.Rule(val)
// 	idx := rand.Intn(len(s.Tiles[variant]))
// 	return s.Tiles[variant][idx]
// }

func getEdgeNeighbors[T any](tilemap Tilemap[T], pos tile.TilePosition) (T, T, T, T) {
	t, _ := tilemap.GetTile(tile.TilePosition{pos.X, pos.Y + 1})
	b, _ := tilemap.GetTile(tile.TilePosition{pos.X, pos.Y - 1})
	l, _ := tilemap.GetTile(tile.TilePosition{pos.X - 1, pos.Y})
	r, _ := tilemap.GetTile(tile.TilePosition{pos.X + 1, pos.Y})
	return t, b, l, r
}

func getEightNeighbors[T any](tilemap Tilemap[T], pos tile.TilePosition) (T, T, T, T, T, T, T, T) {
	t, _ := tilemap.GetTile(tile.TilePosition{pos.X, pos.Y + 1})
	b, _ := tilemap.GetTile(tile.TilePosition{pos.X, pos.Y - 1})
	l, _ := tilemap.GetTile(tile.TilePosition{pos.X - 1, pos.Y})
	r, _ := tilemap.GetTile(tile.TilePosition{pos.X + 1, pos.Y})
	tl, _ := tilemap.GetTile(tile.TilePosition{pos.X - 1, pos.Y + 1})
	tr, _ := tilemap.GetTile(tile.TilePosition{pos.X + 1, pos.Y + 1})
	bl, _ := tilemap.GetTile(tile.TilePosition{pos.X - 1, pos.Y - 1})
	br, _ := tilemap.GetTile(tile.TilePosition{pos.X + 1, pos.Y - 1})

	return t, b, l, r, tl, tr, bl, br
}

func PackedRawEightNumber(t, b, l, r, tl, tr, bl, br bool) uint8 {
	total := uint8(0)
	if t	{ total	+= (1 << 0) }
	if r	{ total	+= (1 << 1) }
	if b	{ total	+= (1 << 2) }
	if l	{ total	+= (1 << 3) }
	if tr { total += (1 << 4) }
	if tl { total += (1 << 5) }
	if br { total += (1 << 6) }
	if bl { total += (1 << 7) }
	return total
}
