package tile

import (
	"math"
)

type Math interface {
	Position(x, y int, size [2]int) (float64, float64)
	PositionToTile(x, y float64, size [2]int) (int, int)
}

type FlatRectMath struct {}
func (t FlatRectMath) Position(x, y int, size [2]int) (float64, float64) {
	return float64(x * size[0]), float64(y * size[1])
}

func (t FlatRectMath) PositionToTile(x, y float64, size [2]int) (int, int) {
	so2x := size[0] >> 1 //Note: Same as: / 2
	so2y := size[1] >> 1 // Note: Same as: / 2
	xPos := (int(x) + (so2x)) / size[0]
	yPos := (int(y) + (so2y)) / size[1]

	// Adjust for negatives
	if x < float64(-so2x) {
		xPos -= 1
	}
	if y < float64(-so2y) {
		yPos -= 1
	}
	return xPos, yPos
	// return (int(x) + (size[0]/2)) / size[0],
	// (int(y) + (size[1]/2))/ size[1]
}

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
