package phy2

import (
	"math"
)
// TODO - Maybe make a new package to hold all vector math?

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
