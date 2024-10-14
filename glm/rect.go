package glm

import "math"

type Rect struct {
	Min, Max Vec2
}

func R(minX, minY, maxX, maxY float64) Rect {
	// TODO - guarantee min is less than max
	return Rect{
		Min: Vec2{minX, minY},
		Max: Vec2{maxX, maxY},
	}
}

// Creates a centered rect
func CR(radius float64) Rect {
	// TODO - guarantee min is less than max
	return Rect{
		Min: Vec2{-radius, -radius},
		Max: Vec2{radius, radius},
	}
}

// Returns a box that holds this rect. The Z axis is 0
func (r Rect) Box() Box {
	return r.ToBox()
}
func (r Rect) ToBox() Box {
	return Box{
		Min: Vec3{r.Min.X, r.Min.Y, 0},
		Max: Vec3{r.Max.X, r.Max.Y, 0},
	}
}

func (r Rect) W() float64 {
	return r.Max.X - r.Min.X
}

func (r Rect) H() float64 {
	return r.Max.Y - r.Min.Y
}

func (r Rect) Center() Vec2 {
	return Vec2{r.Min.X + (r.W() / 2), r.Min.Y + (r.H() / 2)}
}

//	func (r Rect) CenterAt(v Vec2) Rect {
//		return r.Moved(r.Center().Scaled(-1)).Moved(v)
//	}
func (r Rect) WithCenter(v Vec2) Rect {
	w := r.W() / 2
	h := r.H() / 2
	return R(v.X-w, v.Y-h, v.X+w, v.Y+h)
}

// TODO: Should I make a pointer version of this that handles the nil case too?
// Returns the smallest rect which contains both input rects
func (r Rect) Union(s Rect) Rect {
	r = r.Norm()
	s = s.Norm()
	x1 := min(r.Min.X, s.Min.X)
	x2 := max(r.Max.X, s.Max.X)
	y1 := min(r.Min.Y, s.Min.Y)
	y2 := max(r.Max.Y, s.Max.Y)
	return R(x1, y1, x2, y2)
}

func (r Rect) Moved(v Vec2) Rect {
	return Rect{
		Min: r.Min.Add(v),
		Max: r.Max.Add(v),
	}
}

// Calculates the scale required to fit rect r inside r2
func (r Rect) FitScale(r2 Rect) float64 {
	scaleX := r2.W() / r.W()
	scaleY := r2.H() / r.H()

	min := min(scaleX, scaleY)
	return min
}

// Fits rect 'r' into another rect 'r2' with same center but only integer scaled
func (r Rect) FitInt(r2 Rect) Rect {
	scale := math.Floor(r.FitScale(r2))
	return r.Scaled(scale).WithCenter(r2.Center())
}

// Scales rect r uniformly to fit inside rect r2
// TODO This only scales around {0, 0}
func (r Rect) ScaledToFit(r2 Rect) Rect {
	return r.Scaled(r.FitScale(r2))
}

// Returns the largest square that fits inside the rectangle
func (r Rect) SubSquare() Rect {
	w := r.W()
	h := r.H()
	min := min(w, h)
	m2 := min / 2
	return R(-m2, -m2, m2, m2).Moved(r.Center())
}

func (r Rect) CenterScaled(scale float64) Rect {
	c := r.Center()
	w := r.W() * scale / 2.0
	h := r.H() * scale / 2.0
	return R(c.X-w, c.Y-h, c.X+w, c.Y+h)
}

func (r Rect) CenterScaledXY(scaleX, scaleY float64) Rect {
	c := r.Center()
	w := r.W() * scaleX / 2.0
	h := r.H() * scaleY / 2.0
	return R(c.X-w, c.Y-h, c.X+w, c.Y+h)
}

// Note: This scales around the center
// func (r Rect) ScaledXY(scale Vec2) Rect {
// 	c := r.Center()
// 	w := r.W() * scale.X / 2.0
// 	h := r.H() * scale.Y / 2.0
// 	return R(c.X - w, c.Y - h, c.X + w, c.Y + h)
// }

// TODO: I need to deprecate this. This currently just indepentently scales the min and max point which is only useful if the center, min, or max is on (0, 0)
func (r Rect) Scaled(scale float64) Rect {
	// center := r.Center()
	// r = r.Moved(center.Scaled(-1))
	r = Rect{
		Min: r.Min.Scaled(scale),
		Max: r.Max.Scaled(scale),
	}
	// r = r.Moved(center)
	return r
}

func (r Rect) ScaledXY(scale Vec2) Rect {
	r = Rect{
		Min: r.Min.ScaledXY(scale),
		Max: r.Max.ScaledXY(scale),
	}
	return r
}

func (r Rect) Norm() Rect {
	x1, x2 := minMax(r.Min.X, r.Max.X)
	y1, y2 := minMax(r.Min.Y, r.Max.Y)
	return R(x1, y1, x2, y2)
}

func (r Rect) Contains(pos Vec2) bool {
	return pos.X > r.Min.X && pos.X < r.Max.X && pos.Y > r.Min.Y && pos.Y < r.Max.Y
}

// func (r Rect) Contains(x, y float64) bool {
// 	return x > r.Min.X && x < r.Max.X && y > r.Min.Y && y < r.Max.Y
// }

func (r Rect) Intersects(r2 Rect) bool {
	return (r.Min.X <= r2.Max.X &&
		r.Max.X >= r2.Min.X &&
		r.Min.Y <= r2.Max.Y &&
		r.Max.Y >= r2.Min.Y)
}

// Layous out 'n' rectangles horizontally with specified padding between them and returns that rect
// The returned rectangle has a min point of 0,0
func (r Rect) LayoutHorizontal(n int, padding float64) Rect {
	return R(
		0,
		0,
		float64(n)*r.W()+float64(n-1)*padding,
		r.H(),
	)
}

func (r *Rect) CutLeft(amount float64) Rect {
	cutRect := *r
	cutRect.Max.X = cutRect.Min.X + amount
	r.Min.X += amount
	return cutRect
}

func (r *Rect) CutRight(amount float64) Rect {
	cutRect := *r
	cutRect.Min.X = cutRect.Max.X - amount
	r.Max.X -= amount
	return cutRect
}

func (r *Rect) CutBottom(amount float64) Rect {
	cutRect := *r
	cutRect.Max.Y = cutRect.Min.Y + amount
	r.Min.Y += amount
	return cutRect
}

func (r *Rect) CutTop(amount float64) Rect {
	cutRect := *r
	cutRect.Min.Y = cutRect.Max.Y - amount
	r.Max.Y -= amount
	return cutRect
}

// Returns a centered horizontal sliver with height set by amount
func (r Rect) SliceHorizontal(amount float64) Rect {
	r.CutTop((r.H() - amount) / 2)
	return r.CutTop(amount)
}

// Returns a centered vertical sliver with width set by amount
func (r Rect) SliceVertical(amount float64) Rect {
	r.CutRight((r.W() - amount) / 2)
	return r.CutRight(amount)
}

func (r Rect) Snap() Rect {
	r.Min = r.Min.Snap()
	r.Max = r.Max.Snap()
	return r
}

// Adds padding to a rectangle consistently
func (r Rect) PadAll(padding float64) Rect {
	return r.Pad(R(padding, padding, padding, padding))
}

// Adds padding to a rectangle (pads inward if padding is negative)
func (r Rect) Pad(pad Rect) Rect {
	return R(r.Min.X-pad.Min.X, r.Min.Y-pad.Min.Y, r.Max.X+pad.Max.X, r.Max.Y+pad.Max.Y)
}

// Removes padding from a rectangle (pads outward if padding is negative). Essentially calls pad but with negative values
func (r Rect) Unpad(pad Rect) Rect {
	return r.Pad(pad.Scaled(-1))
}

// Takes r2 and places it in r based on the alignment
// TODO - rename to InnerAnchor?
func (r Rect) Anchor(r2 Rect, anchor Vec2) Rect {
	// Anchor point is the position in r that we are anchoring to
	anchorPoint := Vec2{r.Min.X + (anchor.X * r.W()), r.Min.Y + (anchor.Y * r.H())}
	pivotPoint := Vec2{r2.Min.X + (anchor.X * r2.W()), r2.Min.Y + (anchor.Y * r2.H())}

	// fmt.Println("Anchor:", anchorPoint)
	// fmt.Println("Pivot:", pivotPoint)

	a := Vec2{anchorPoint.X - pivotPoint.X, anchorPoint.Y - pivotPoint.Y}
	return R(a.X, a.Y, a.X+r2.W(), a.Y+r2.H()).Norm()
}

// Anchors r2 to r1 based on two anchors, one for r and one for r2
// TODO - rename to Anchor?
func (r Rect) FullAnchor(r2 Rect, anchor, pivot Vec2) Rect {
	anchorPoint := Vec2{r.Min.X + (anchor.X * r.W()), r.Min.Y + (anchor.Y * r.H())}
	pivotPoint := Vec2{r2.Min.X + (pivot.X * r2.W()), r2.Min.Y + (pivot.Y * r2.H())}

	a := Vec2{anchorPoint.X - pivotPoint.X, anchorPoint.Y - pivotPoint.Y}
	return R(a.X, a.Y, a.X+r2.W(), a.Y+r2.H()).Norm()
}

// Move the min point of the rect to a certain position
func (r Rect) MoveMin(pos Vec2) Rect {
	dv := r.Min.Sub(pos)
	return r.Moved(dv)
}

func lerp(a, b float64, t float64) float64 {
	m := b - a // Slope = Rise over run | Note: Run = (1 - 0)
	y := (m * t) + a
	return y
}

// returns the min, max of the two numbers
func minMax(a, b float64) (float64, float64) {
	if a > b {
		return b, a
	}
	return a, b
}

func (r Rect) RectDraw(r2 Rect) Mat4 {
	srcCenter := r.Center()
	dstCenter := r2.Center()
	mat := Mat4Ident
	mat.
		Translate(-srcCenter.X, -srcCenter.Y, 0).
		Scale(r2.W()/r.W(), r2.H()/r.H(), 1).
		Translate(dstCenter.X, dstCenter.Y, 0)
	return mat
}
