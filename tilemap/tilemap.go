package tilemap

type TileType uint8

type Tile struct {
	Type TileType
	Height float32
}

type TilePosition struct {
	X, Y int
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
