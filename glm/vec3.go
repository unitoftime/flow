package glm

import "math"

func (v Vec3) Add(u Vec3) Vec3 {
	return Vec3{v.X + u.X, v.Y + u.Y, v.Z + u.Z}
}

func (v Vec3) Sub(u Vec3) Vec3 {
	return Vec3{v.X - u.X, v.Y - u.Y, v.Z - u.Z}
}

// Finds the dot product of two vectors
func (v Vec3) Dot(u Vec3) float64 {
	return (v.X * u.X) + (v.Y * u.Y) + (v.Z * u.Z)
}

// Finds the angle between two vectors
// TODO: Is this correct?
func (v Vec3) Angle(u Vec3) float64 {
	return math.Acos(v.Dot(u) / (v.Len() * u.Len()))
}

func (v Vec3) Theta() float64 {
	return math.Atan2(v.Y, v.X)
}

// Rotates the vector by theta on the XY 2d plane
func (v Vec3) Rotate2D(theta float64) Vec3 {
	t := theta
	x := v.X
	y := v.Y
	x1 := x*math.Cos(t) - y*math.Sin(t)
	y1 := x*math.Sin(t) + y*math.Cos(t)
	return Vec3{x1, y1, v.Z}
}

func (v Vec3) Len() float64 {
	// return float32(math.Hypot(float64(v.X), float64(v.Y)))
	a := v.X
	b := v.Y
	c := v.Z
	return math.Sqrt((a * a) + (b * b) + (c * c))
}

func (v Vec3) Vec2() Vec2 {
	return Vec2{v.X, v.Y}
}

func (v Vec3) Unit() Vec3 {
	len := v.Len()
	return Vec3{v.X / len, v.Y / len, v.Z / len}
}

func (v Vec3) Scaled(x, y, z float64) Vec3 {
	v.X *= x
	v.Y *= y
	v.Z *= z

	return v
}

func (v Vec3) Mult(s Vec3) Vec3 {
	return Vec3{
		X: v.X * s.X,
		Y: v.Y * s.Y,
		Z: v.Z * s.Z,
	}

	return v
}
