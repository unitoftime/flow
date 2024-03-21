package phy2

import (
	"math"
)

// TODO - Maybe make a new package to hold all vector math?

type Vec = Vec2

//cod:struct
type Vec2 struct {
	X, Y float64
}

func V2(x, y float64) Vec2 {
	return Vec2{x, y}
}

func (v Vec2) Add(v2 Vec2) Vec2 {
	return V2(v.X + v2.X, v.Y + v2.Y)
}

func (v Vec2) Sub(v2 Vec2) Vec2 {
	return V2(v.X - v2.X, v.Y - v2.Y)
}

func (v Vec2) Norm() Vec2 {
	len := v.Len()
	if len == 0 {
		return Vec2{}
	}
	return V2(v.X/len, v.Y/len)
}

func (v Vec2) Dist(u Vec2) float64 {
	return v.Sub(u).Len()
}

func (v Vec2) DistSq(u Vec2) float64 {
	return v.Sub(u).LenSq()
}

func (v Vec2) Dot(u Vec2) float64 {
	return (v.X * u.X) + (v.Y * u.Y)
}

// Returns the length of the vector
func (v Vec2) Len() float64 {
	return math.Sqrt((v.X * v.X) + (v.Y * v.Y))
}

// Returns the length of the vector squared. Note this is slightly faster b/c it doesn't take the square root
func (v Vec2) LenSq() float64 {
	return (v.X * v.X) + (v.Y * v.Y)
}

func (v Vec2) Scaled(s float64) Vec2 {
	return V2(s * v.X, s * v.Y)
}

func (v Vec2) Rotated(radians float64) Vec2 {
	sin := math.Sin(radians)
	cos := math.Cos(radians)
	return V2(
		v.X * cos - v.Y * sin,
		v.X * sin + v.Y * cos,
	)
}

func (v Vec) Angle() float64 {
	return math.Atan2(v.Y, v.X)
}

// Finds the angle between two vectors
func Angle(a, b Vec) float64 {
	angle := a.Angle() - b.Angle()
	if angle > math.Pi {
		angle -= 2 * math.Pi
	} else if angle <= -math.Pi {
		angle += 2 * math.Pi
	}
	return angle
}

// --------------------------------------------------------------------------------
// - Rect
// --------------------------------------------------------------------------------
type Rect struct {
	Min, Max Vec
}

// Returns a zero'd rect
// func ZR(width, height float64) Rect {
// 	return Rect{minX, minY, maxX, maxY}
// }

// Returns a rect with specified dimensions
func R(minX, minY, maxX, maxY float64) Rect {
	return Rect{
		Vec{minX, minY},
		Vec{maxX, maxY},
	}
}

// Returns a centered Rect
func CR(center Vec2, radius Vec2) Rect {
	return R(
		center.X - float64(radius.X), center.Y - float64(radius.Y),
		center.X + float64(radius.Y), center.Y + float64(radius.Y),
	)
}

func (r Rect) WithCenter(v Vec2) Rect {
	// zRect := r.Moved(r.Center().Scaled(-1))
	// return zRect.Moved(v)
	w := r.W()/2
	h := r.H()/2
	return R(v.X - w, v.Y - h, v.X + w, v.Y + h)
}

func (r Rect) W() float64 {
	return r.Max.X - r.Min.X
}

func (r Rect) H() float64 {
	return r.Max.Y - r.Min.Y
}

func (r Rect) Center() Vec2 {
	return Vec2{r.Min.X + (r.W()/2), r.Min.Y + (r.H()/2)}
}

func (r Rect) Moved(v Vec2) Rect {
	return Rect{
		Min: r.Min.Add(v),
		Max: r.Max.Add(v),
	}
}

func (r Rect) Contains(pos Vec) bool {
	return pos.X > r.Min.X && pos.X < r.Max.X && pos.Y > r.Min.Y && pos.Y < r.Max.Y
}

func (r Rect) Intersects(r2 Rect) bool {
	return (
		r.Min.X <= r2.Max.X &&
			r.Max.X >= r2.Min.X &&
			r.Min.Y <= r2.Max.Y &&
			r.Max.Y >= r2.Min.Y)
}

func (r Rect) Pad(pad Rect) Rect {
	return R(r.Min.X - pad.Min.X, r.Min.Y - pad.Min.Y, r.Max.X + pad.Max.X, r.Max.Y + pad.Max.Y)
}

func (r Rect) Unpad(pad Rect) Rect {
	return R(r.Min.X + pad.Min.X, r.Min.Y + pad.Min.Y, r.Max.X - pad.Max.X, r.Max.Y - pad.Max.Y)
}


// --------------------------------------------------------------------------------
// - Circle
// --------------------------------------------------------------------------------
type Circle struct {
	Center Vec2
	Radius float64
}
func NewCircle(center Vec2, radius float64) Circle {
	return Circle{
		Center: center,
		Radius: radius,
	}
}

// // TODO: Not tested
// func (c Circle) Intersects(c2 Circle) bool {
// 	dx := c.Center.X - c2.Center.X
// 	dy := c.Center.Y - c2.Center.Y
// 	distSq := (dx * dx) + (dy * dy)
// 	totalRadiusSq := c.Radius + c2.Radius
// 	totalRadiusSq = totalRadiusSq * totalRadiusSq
// 	return dist <= totalRadiusSq
// }

// // --------------------------------------------------------------------------------
// // - Vec3
// // --------------------------------------------------------------------------------
// type Vec3 struct {
// 	X, Y, Z float64
// }

// func V3(x, y, z float64) Vec3 {
// 	return Vec3{x, y, z}
// }

// func (v Vec3) Add(v2 Vec3) Vec3 {
// 	return V3(v.X + v2.X, v.Y + v2.Y, v.Z + v2.Z)
// }

// func (v Vec3) Sub(v2 Vec3) Vec3 {
// 	return V3(v.X - v2.X, v.Y - v2.Y, v.Z - v2.Z)
// }

// func (v Vec3) Norm() Vec3 {
// 	len := v.Len()
// 	if len == 0 {
// 		return Vec3{}
// 	}
// 	return V3(v.X/len, v.Y/len, v.Z/len)
// }

// func (v Vec3) Len() float64 {
// 	return math.Sqrt((v.X * v.X) + (v.Y * v.Y) + (v.Z * v.Z))
// }

// func (v Vec3) Scaled(s float64) Vec3 {
// 	return V3(s * v.X, s * v.Y, s * v.Z)
// }

// // TODO - Need to rotate around something
// // func (v Vec3) Rotated(radians float64) Vec3 {
// // 	sin := math.Sin(radians)
// // 	cos := math.Cos(radians)
// // 	return V3(
// // 		v.X * cos - v.Y * sin,
// // 		v.X * sin + v.Y * cos,
// // 	)
// // }
