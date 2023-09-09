package flow

// // import (
// // 	"github.com/ungerik/go3dfloat64/quaternion"
// // )

// // // TODO: float64 vs float32?

// type Vec3 struct {
// 	X, Y, Z float64
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

// // type Quat struct {
// // 	quaternion.T
// // }

// type Transform struct {
// 	Position Vec3
// 	Scale Vec3
// 	// Rotation Quat
// }

// // func (t Transform) Matrix() {
	
// // 	// TODO: quat
// // }
// func NewTransform() Transform {
// 	return Transform{
// 		Position: Vec3{},
// 		Scale: Vec3{},
// 	}
// }
