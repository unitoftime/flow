package physics

import (
	"math"
)
// TODO - Maybe make a new package to hold all vector math?

// I feel like using generics here just complicate client code too much
// import (
// 	"constraints"
// )
// type Number interface {
// 	constraints.Integer | constraints.Float
// }

// type Vec2[T Number] struct {
// 	X, Y T
// }

// func V2[T Number](x, y T) Vec2[T] {
// 	return Vec2[T]{x, y}
// }

// func (v Vec2[T])Sub(v2 Vec2[T]) Vec2[T] {
// 	return V2(v.X - v2.X, v.Y - v2.Y)
// }

type Vec3 struct {
	X, Y, Z float64
}

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
