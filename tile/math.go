package tile

import (
	"math"
)

type Math interface {
	Position(x, y int, size [2]int) (float32, float32)
	PositionToTile(x, y float32, size [2]int) (int, int)
}

type FlatRectMath struct {}
func (t FlatRectMath) Position(x, y int, size [2]int) (float32, float32) {
	return float32(x * size[0]), float32(y * size[1])
}

func (t FlatRectMath) PositionToTile(x, y float32, size [2]int) (int, int) {
	xPos := (int(x) + (size[0]/2)) / size[0]
	yPos := (int(y) + (size[1]/2))/ size[1]

	// Adjust for negatives
	if x < float32(-size[0] / 2) {
		xPos -= 1
	}
	if y < float32(-size[1] / 2) {
		yPos -= 1
	}
	return xPos, yPos
	// return (int(x) + (size[0]/2)) / size[0],
	// (int(y) + (size[1]/2))/ size[1]
}

type PointyRectMath struct {}
func (t PointyRectMath) Position(x, y int, size [2]int) (float32, float32) {
	// If y goes up, then xPos must go downward a bit
	return -float32((x * size[0] / 2) - (y * size[0] / 2)),
		// If x goes up, then yPos must go up a bit as well
	-float32((y * size[1] / 2) + (x * size[1] / 2))
}

func (t PointyRectMath) PositionToTile(x, y float32, size [2]int) (int, int) {
	dx := float32(size[0]) / 2.0
	dy := float32(size[1]) / 2.0

	tx := -((x/dx) + (y/dy))/2
	ty := -((y/dy) + tx)
	return int(math.Round(float64(tx))), int(math.Round(float64(ty)))
}
