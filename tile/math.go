package tile

import (
	"math"
)

type Math interface {
	Position(x, y int, size [2]int) (float64, float64)
	PositionToTile(x, y float64, size [2]int) (int, int)
}

type FlatRectMath struct {
	size [2]int
	sizeOver2 [2]int
	div [2]int
}
func NewFlatRectMath(size [2]int) FlatRectMath {
	divX := int(math.Log2(float64(size[0])))
	divY := int(math.Log2(float64(size[1])))
	if (1 << divX) != size[0] || (1 << divY) != size[1] {
		panic("Tile maps must have a chunksize and tilesize that is a power of 2!")
	}
	return FlatRectMath{
		size: size,
		sizeOver2: [2]int{size[0]/2, size[1]/2},
		div: [2]int{divX, divY},
	}
}
func (m FlatRectMath) Position(x, y int) (float64, float64) {
	return float64(x * m.size[0]), float64(y * m.size[1])
}

func (m FlatRectMath) PositionToTile(x, y float64) (int, int) {
	// xPos := (int(x) + (m.sizeOver2[0])) / size[0]
	// yPos := (int(y) + (m.sizeOver2[1])) / size[1]
	xPos := (int(x) + m.sizeOver2[0]) >> m.div[0]
	yPos := (int(y) + m.sizeOver2[1]) >> m.div[1]

	// xPos := (int(x) + m.sizeOver2[0]) >> m.div[0]
	// yPos := (int(y) + m.sizeOver2[1]) >> m.div[1]

	// // Adjust for negatives
	// if x < float64(-m.sizeOver2[0]) {
	// 	xPos -= 1
	// }
	// if y < float64(-m.sizeOver2[1]) {
	// 	yPos -= 1
	// }
	return xPos, yPos
	// return (int(x) + (size[0]/2)) / size[0],
	// (int(y) + (size[1]/2))/ size[1]
}

// func (t FlatRectMath) Position(x, y int, size [2]int) (float64, float64) {
// 	return float64(x * size[0]), float64(y * size[1])
// }

// func (t FlatRectMath) PositionToTile(x, y float64, size [2]int) (int, int) {
// 	so2x := size[0] >> 1 //Note: Same as: / 2
// 	so2y := size[1] >> 1 // Note: Same as: / 2
// 	xPos := (int(x) + (so2x)) / size[0]
// 	yPos := (int(y) + (so2y)) / size[1]

// 	// Adjust for negatives
// 	if x < float64(-so2x) {
// 		xPos -= 1
// 	}
// 	if y < float64(-so2y) {
// 		yPos -= 1
// 	}
// 	return xPos, yPos
// 	// return (int(x) + (size[0]/2)) / size[0],
// 	// (int(y) + (size[1]/2))/ size[1]
// }

type PointyRectMath struct {}
func (t PointyRectMath) Position(x, y int, size [2]int) (float64, float64) {
	// If y goes up, then xPos must go downward a bit
	return -float64((x * size[0] / 2) - (y * size[0] / 2)),
		// If x goes up, then yPos must go up a bit as well
	-float64((y * size[1] / 2) + (x * size[1] / 2))
}

func (t PointyRectMath) PositionToTile(x, y float64, size [2]int) (int, int) {
	dx := float64(size[0]) / 2.0
	dy := float64(size[1]) / 2.0

	tx := -((x/dx) + (y/dy))/2
	ty := -((y/dy) + tx)
	return int(math.Round(float64(tx))), int(math.Round(float64(ty)))
}
