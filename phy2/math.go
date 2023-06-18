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

func (v Vec2) Dot(u Vec2) float64 {
	return (v.X * u.X) + (v.Y * u.Y)
}

func (v Vec2) Len() float64 {
	return math.Sqrt((v.X * v.X) + (v.Y * v.Y))
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

func (r Rect) WithCenter(v Vec2) Rect {
	zRect := r.Moved(r.Center().Scaled(-1))
	return zRect.Moved(v)
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

func (r Rect) Contains(pos Pos) bool {
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
