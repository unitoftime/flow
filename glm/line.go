package glm

type Line2 struct {
	A, B Vec2
}

// Returns true if the two lines intersect
func (l1 Line2) Intersects(l2 Line2) bool {
	// Adapted from: https://www.gorillasun.de/blog/an-algorithm-for-polygon-intersections/
	x1 := l1.A.X
	y1 := l1.A.Y
	x2 := l1.B.X
	y2 := l1.B.Y

	x3 := l2.A.X
	y3 := l2.A.Y
	x4 := l2.B.X
	y4 := l2.B.Y

	// Ensure no line is 0 length
	if ((x1 == x2 && y1 == y2) || (x3 == x4 && y3 == y4)) {
		return false
	}

	den := ((y4 - y3) * (x2 - x1) - (x4 - x3) * (y2 - y1))

	// Lines are parallel
	if den == 0 {
		return false
	}

	ua := ((x4 - x3) * (y1 - y3) - (y4 - y3) * (x1 - x3)) / den
	ub := ((x2 - x1) * (y1 - y3) - (y2 - y1) * (x1 - x3)) / den

  // is the intersection along the segments
  if (ua < 0 || ua > 1 || ub < 0 || ub > 1) {
		return false
  }

	return true

  // // Return a object with the x and y coordinates of the intersection
  // x := x1 + ua * (x2 - x1)
  // y := y1 + ua * (y2 - y1)

  // return {x, y}
}
