package glm

import "math"

type Vec2 struct {
	X, Y float64
}

func V2(x, y float64) Vec2 {
	return Vec2{x, y}
}

func (v Vec2) Vec3() Vec3 {
	return Vec3{v.X, v.Y, 0}
}

func (v Vec2) Add(v2 Vec2) Vec2 {
	return V2(v.X+v2.X, v.Y+v2.Y)
}

func (v Vec2) Sub(v2 Vec2) Vec2 {
	return V2(v.X-v2.X, v.Y-v2.Y)
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
	return V2(s*v.X, s*v.Y)
}

func (v Vec2) ScaledXY(s Vec2) Vec2 {
	return Vec2{v.X * s.X, v.Y * s.Y}
}

func (v Vec2) Rotated(radians float64) Vec2 {
	sin := math.Sin(radians)
	cos := math.Cos(radians)
	return V2(
		v.X*cos-v.Y*sin,
		v.X*sin+v.Y*cos,
	)
}

func (v Vec2) Angle() float64 {
	return math.Atan2(v.Y, v.X)
}

// Finds the angle between two vectors
func Angle(a, b Vec2) float64 {
	angle := a.Angle() - b.Angle()
	if angle > math.Pi {
		angle -= 2 * math.Pi
	} else if angle <= -math.Pi {
		angle += 2 * math.Pi
	}
	return angle
}

func (v Vec2) Snap() Vec2 {
	return Vec2{
		math.Round(v.X),
		math.Round(v.Y),
	}
}
