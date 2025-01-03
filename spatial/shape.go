package spatial

import (
	"github.com/unitoftime/flow/glm"
)

type ShapeType uint8

// TODO: Tagged Union?
const (
	ShapeAABB ShapeType = iota // A rectangle, not rotated nor scaled
	ShapeRect                  // Can be rotated or scaled
	// ShapeCircle // TODO: A circle
	// ShapeRing // TODO: A ring (Circle with excluded middle segment
	// ShapeEllipses // TODO: Maybe combine with circle?
	// ShapePolygon // TODO: An arbitrary convex polygon (Question: How to support arbitrary length points)
)

type Shape struct {
	Type    ShapeType
	Bounds  glm.Rect    // The bounding AABB
	Vectors [4]glm.Vec2 // Vector data which is stored differently depending on the shape type
}

func AABB(rect glm.Rect) Shape {
	return Shape{
		Type:   ShapeAABB,
		Bounds: rect,
		Vectors: [4]glm.Vec2{
			rect.Min,
			rect.TL(),
			rect.Max,
			rect.BR(),
		},
	}
}
func Rect(r glm.Rect, mat glm.Mat4) Shape {
	tl := mat.Apply(glm.Vec3{r.Min.X, r.Max.Y, 0}).Vec2()
	br := mat.Apply(glm.Vec3{r.Max.X, r.Min.Y, 0}).Vec2()

	bl := mat.Apply(r.Min.Vec3()).Vec2()
	tr := mat.Apply(r.Max.Vec3()).Vec2()

	xMin := min(tl.X, br.X, bl.X, tr.X)
	yMin := min(tl.Y, br.Y, bl.Y, tr.Y)

	xMax := max(tl.X, br.X, bl.X, tr.X)
	yMax := max(tl.Y, br.Y, bl.Y, tr.Y)

	bounds := glm.R(xMin, yMin, xMax, yMax)
	return Shape{
		Type:   ShapeRect,
		Bounds: bounds,
		Vectors: [4]glm.Vec2{
			bl, tl, tr, br,
		},
	}
}

func (s Shape) Intersects(s2 Shape) bool {
	switch s.Type {
	case ShapeAABB:
		{
			switch s2.Type {
			case ShapeAABB:
				return s.Bounds.Intersects(s2.Bounds)
			case ShapeRect:
				if !s.Bounds.Intersects(s2.Bounds) {
					return false // Bounding box must intersect
				}

				return polygonIntersectionCheck(s.Vectors[:], s2.Vectors[:])
			}
		}
	case ShapeRect:
		{
			switch s2.Type {
			case ShapeAABB:
				if !s.Bounds.Intersects(s2.Bounds) {
					return false // Bounding box must intersect
				}

				return polygonIntersectionCheck(s.Vectors[:], s2.Vectors[:])
			case ShapeRect:
				if !s.Bounds.Intersects(s2.Bounds) {
					return false // Bounding box must intersect
				}

				return polygonIntersectionCheck(s.Vectors[:], s2.Vectors[:])
			}
		}
	}

	// Invalid case
	return false
}

func polygonIntersectionCheck(a, b []glm.Vec2) bool {
	// First check if at least one point is inside. we need to do this to handle the case where one polygon is entirely inside the other
	for i := range a {
		if polygonContainsPoint(b, a[i]) {
			return true
		}
	}

	for i := range b {
		if polygonContainsPoint(a, b[i]) {
			return true
		}
	}

	// Else no points are inside, so we test all edge pairs
	lenA := len(a)
	lenB := len(b)
	var l1, l2 glm.Line2
	for i := range a {
		l1.A = a[i]
		l1.B = a[(i+1)%lenA]
		for j := range b {
			l2.A = b[j]
			l2.B = b[(j+1)%lenB]

			if l1.Intersects(l2) {
				return true
			}
		}
	}

	return false
}

// func polygonContainsPoint(poly []glm.Vec2, point glm.Vec2) bool {
// 	count := 0
// 	length := len(poly)
// 	for i := range poly {
// 		a := poly[i]
// 		b := poly[(i+1) % length]

// 		if (point.Y < a.Y) != (point.Y < b.Y) && point.X < a.X + ((point.Y - a.Y)/(b.Y - a.Y)) * (b.X - a.X) {
// 			count++
// 		}
// 	}
// 	return (count % 2) == 1
// }

func polygonContainsPoint(poly []glm.Vec2, point glm.Vec2) bool {
	// Ref: https://stackoverflow.com/questions/22521982/check-if-point-is-inside-a-polygon
	// https://wrf.ecse.rpi.edu/Research/Short_Notes/pnpoly.html

	x := point.X
	y := point.Y

	length := len(poly)

	var inside = false
	for i := range poly {
		a := poly[i]
		b := poly[(i+1)%length]
		xi := a.X
		yi := a.Y
		xj := b.X
		yj := b.Y

		intersect := ((yi > y) != (yj > y)) && (x < (xj-xi)*(y-yi)/(yj-yi)+xi)
		if intersect {
			inside = !inside
		}
	}

	return inside
}

// Note: Originally I was doing this, but it actually takes way more projections than I thought. So I opted to just do a generic line intersection test for all edges of both rectangles. Even though this one may be faster. Maybe I"ll add it back in the future
// TODO: This was almost done. I think all that remained was that I need to
// 1. Repeat for both axes and make sure all segments overlap
// Reference: https://math.stackexchange.com/questions/1278665/how-to-check-if-two-rectangles-intersect-rectangles-can-be-rotated
// 1. calc new axis
// 2. project b points onto new axis
// 3. Check if any point overlaps 'rect A' ranges
// 2. (If only one is rotated) The axis and projection calculations
// func (a RectShape) Intersects(b RectShape) bool {
// 	// return a.Bounds.Intersects(b.Bounds)

// 	// 1. calc new axis
// 	axisXVec := a.BR.Sub(a.BL)
// 	axisYVec := a.TL.Sub(a.BL)
// 	axisX := axisXVec.Norm()
// 	axisY := axisYVec.Norm()

// 	fmt.Println("axisX", axisX)
// 	fmt.Println("axisY", axisY)

// 	// 2. project b points onto new axis
// 	newBL := b.BL.Sub(a.BL)
// 	newTL := b.TL.Sub(a.BL)
// 	newTR := b.TR.Sub(a.BL)
// 	newBR := b.BR.Sub(a.BL)

// 	// X
// 	projBL_X := newBL.Dot(axisX)
// 	projTL_X := newTL.Dot(axisX)
// 	projTR_X := newTR.Dot(axisX)
// 	projBR_X := newBR.Dot(axisX)

// 	// Y
// 	projBL_Y := newBL.Dot(axisY)
// 	projTL_Y := newTL.Dot(axisY)
// 	projTR_Y := newTR.Dot(axisY)
// 	projBR_Y := newBR.Dot(axisY)

// 	fmt.Println("projX", projBL_X, projTL_X, projTR_X, projBR_X)
// 	fmt.Println("projY", projBL_Y, projTL_Y, projTR_Y, projBR_Y)

// 	// 3. Check if any point overlaps 'rect A' ranges
// 	rangeX := axisXVec.Len()
// 	rangeY := axisYVec.Len()

// 	fmt.Println("rangeX:", rangeX)
// 	fmt.Println("rangeY:", rangeY)

// 	// This is recalculated around the new Axis
// 	rectA := glm.R(0, 0, rangeX, rangeY)

// 	// TODO: It might be faster to reorganize the search to do one axis at a time? Idk
// 	// Check Rect B point BL
// 	if rectA.Contains(glm.Vec2{projBL_X, projBL_Y}) {
// 		return true
// 	}

// 	// Check Rect B point TL
// 	if rectA.Contains(glm.Vec2{projTL_X, projTL_Y}) {
// 		return true
// 	}

// 	// Check Rect B point TR
// 	if rectA.Contains(glm.Vec2{projTR_X, projTR_Y}) {
// 		return true
// 	}

// 	// Check Rect B point BR
// 	if rectA.Contains(glm.Vec2{projBR_X, projBR_Y}) {
// 		return true
// 	}

// 	return false
// }
