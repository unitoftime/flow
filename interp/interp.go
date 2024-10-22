package interp

import (
	"math"
	"time"

	"github.com/ungerik/go3d/float64/bezier2"
	"github.com/ungerik/go3d/float64/vec2"
	"github.com/unitoftime/flow/glm"
)

// TODO: use https://easings.net/
// Note: https://cubic-bezier.com

// This will calculate
func DynamicValue(val float64, fixedTime, dt time.Duration) float64 {
	// interpVal := val * dt.Seconds() / (16 * time.Millisecond).Seconds()
	interpVal := (val / fixedTime.Seconds()) * dt.Seconds()
	if interpVal > 1.0 {
		return 1.0
	} else if interpVal < 0 {
		return 0.0
	}
	return interpVal
}

var Linear *Lerp = &Lerp{}

var EaseOut Bezier = Bezier{
	bezier2.T{
		vec2.T{0.0, 0.0},
		vec2.T{1.0, 0.0},
		vec2.T{1.0, 0.0},
		vec2.T{1.0, 1.0},
	},
}
var EaseIn Bezier = Bezier{
	bezier2.T{
		vec2.T{0.0, 0.0},
		vec2.T{0.0, 1.0},
		vec2.T{0.0, 1.0},
		vec2.T{1.0, 1.0},
	},
}
var EaseInOut Bezier = Bezier{
	bezier2.T{
		vec2.T{0.0, 0.0},
		vec2.T{1.0, 0.0},
		vec2.T{0.0, 1.0},
		vec2.T{1.0, 1.0},
	},
}
var BezLerp Bezier = Bezier{
	bezier2.T{
		vec2.T{0.0, 0.0},
		vec2.T{0.0, 0.0},
		vec2.T{1.0, 1.0},
		vec2.T{1.0, 1.0},
	},
}

var BezFlash Bezier = Bezier{
	bezier2.T{
		vec2.T{0.0, 0.0},
		vec2.T{0.1, 1.0},
		vec2.T{0.2, 0.5},
		vec2.T{0.0, 0.0},
	},
}

func NewBezier(a, b, c, d glm.Vec2) Bezier {
	return Bezier{
		bezier2.T{
			vec2.T{a.X, a.Y},
			vec2.T{b.X, b.Y},
			vec2.T{c.X, c.Y},
			vec2.T{d.X, d.Y},
		},
	}
}

// Note: https://www.w3schools.com/cssref/func_cubic-bezier.php#:~:text=P0%20is%20(0%2C%200),transition%2Dtiming%2Dfunction%20property.
// Essentially: First point is (0, 0), last point is (1, 1) and you can define the two points in the middle
func NewCubicBezier(b, c glm.Vec2) Bezier {
	return Bezier{
		bezier2.T{
			vec2.T{0, 0},
			vec2.T{b.X, b.Y},
			vec2.T{c.X, c.Y},
			vec2.T{1, 1},
		},
	}
}

func Step(divRatio float64, a, b Interp) StepF {
	return StepF{
		a:        a,
		b:        b,
		divRatio: divRatio,
	}
}

type StepF struct {
	a, b     Interp
	divRatio float64
}

func (i StepF) get(t float64) (Interp, float64) {
	if t < i.divRatio {
		newT := t / (i.divRatio - 0)
		return i.a, newT
	}

	newT := t / (1 - i.divRatio)
	return i.b, newT
}
func (i StepF) Float64(a, b float64, t float64) float64 {
	itrp, val := i.get(t)
	return itrp.Float64(a, b, val)
}
func (i StepF) Float32(a, b float32, t float64) float32 {
	itrp, val := i.get(t)
	return itrp.Float32(a, b, val)
}
func (i StepF) Uint8(a, b uint8, t float64) uint8 {
	itrp, val := i.get(t)
	return itrp.Uint8(a, b, val)
}
func (i StepF) Vec2(a, b glm.Vec2, t float64) glm.Vec2 {
	itrp, val := i.get(t)
	return itrp.Vec2(a, b, val)
}

func Const(val float64) Bezier {
	a := vec2.T{0, val}
	return Bezier{
		bezier2.T{a, a, a, a},
	}
}

// var Sinusoid *Equation = &Equation{
// 	Func: SinFunc{},
// }

// https://cubic-bezier.com/#.22,1,.36,1
var EaseOutQuint Bezier = Bezier{
	bezier2.T{
		vec2.T{0.0, 0.0},
		vec2.T{0.22, 1.0},
		vec2.T{0.36, 1.0},
		vec2.T{1.0, 1.0},
	},
}

var EaseTest Bezier = Bezier{
	bezier2.T{
		vec2.T{0.0, 0.0},
		vec2.T{0.0, 1.0},
		vec2.T{0.0, 2.0},
		vec2.T{1.0, 1.0},
	},
}

type Interp interface {
	Uint8(uint8, uint8, float64) uint8
	Float32(float32, float32, float64) float32
	Float64(float64, float64, float64) float64
	Vec2(glm.Vec2, glm.Vec2, float64) glm.Vec2
}

type Lerp struct {
}

func (i Lerp) Float64(a, b float64, t float64) float64 {
	m := b - a // Slope = Rise over run | Note: Run = (1 - 0)
	y := (m * t) + a
	return y
}
func (i Lerp) Float32(a, b float32, t float64) float32 {
	m := b - a // Slope = Rise over run | Note: Run = (1 - 0)
	y := (m * float32(t)) + a
	return y
}
func (i Lerp) Uint8(a, b uint8, t float64) uint8 {
	return uint8(i.Float64(float64(a), float64(b), t))
}

func (i Lerp) Vec2(a, b glm.Vec2, t float64) glm.Vec2 {
	ret := glm.Vec2{
		i.Float64(a.X, b.X, t),
		i.Float64(a.Y, b.Y, t),
	}
	return ret
}

type Bezier struct {
	bezier2.T
}

func (i Bezier) Float64(a, b float64, t float64) float64 {
	iValue := i.T.Point(t)
	return Linear.Float64(a, b, iValue[1])
}
func (i Bezier) Float32(a, b float32, t float64) float32 {
	iValue := i.T.Point(t)
	return Linear.Float32(a, b, iValue[1])
}
func (i Bezier) Uint8(a, b uint8, t float64) uint8 {
	iValue := i.T.Point(t)
	return Linear.Uint8(a, b, iValue[1])
}
func (i Bezier) Vec2(a, b glm.Vec2, t float64) glm.Vec2 {
	iValue := i.T.Point(t)
	return Linear.Vec2(a, b, iValue[1])
}

type Sine struct {}

func (i Sine) Float64(a, b float64, t float64) float64 {
	iValue := math.Sin(t * math.Pi)
	return Linear.Float64(a, b, iValue)
}
func (i Sine) Float32(a, b float32, t float64) float32 {
	iValue := math.Sin(t * math.Pi)
	return Linear.Float32(a, b, iValue)
}
func (i Sine) Uint8(a, b uint8, t float64) uint8 {
	iValue := math.Sin(t * math.Pi)
	return Linear.Uint8(a, b, iValue)
}
func (i Sine) Vec2(a, b glm.Vec2, t float64) glm.Vec2 {
	iValue := math.Sin(t * math.Pi)
	return Linear.Vec2(a, b, iValue)
}

type Equation struct {
	Func Function
}

func (i *Equation) Float64(a, b float64, t float64) float64 {
	iValue := i.Func.Interp(t)
	return Linear.Float64(a, b, iValue)
}
func (i *Equation) Float32(a, b float32, t float64) float32 {
	iValue := i.Func.Interp(t)
	return Linear.Float32(a, b, iValue)
}
func (i *Equation) Uint8(a, b uint8, t float64) uint8 {
	iValue := i.Func.Interp(t)
	return Linear.Uint8(a, b, iValue)
}
func (i *Equation) Vec2(a, b glm.Vec2, t float64) glm.Vec2 {
	iValue := i.Func.Interp(t)
	return Linear.Vec2(a, b, iValue)
}

type Function interface {
	Interp(t float64) float64
}

type SinFunc struct {
	Radius float64
	Freq   float64
	ShiftY float64
}

func (s SinFunc) Interp(t float64) float64 {
	return s.Radius * (s.ShiftY + math.Sin(t*s.Freq))
}

type CosFunc struct {
	Radius float64
	Freq   float64
	ShiftY float64
}

func (s CosFunc) Interp(t float64) float64 {
	return s.Radius * (s.ShiftY + math.Cos(t*s.Freq))
}

type BezFunc struct {
	Radius float64
	Dur    float64
	Bezier Bezier
}

func (f BezFunc) Interp(t float64) float64 {
	if f.Dur == 0 {
		return f.Radius * f.Bezier.Float64(0, 1, t)
	} else {
		// else normalize by the dur
		if t >= f.Dur {
			return f.Radius
		}
		norm := t / f.Dur
		return f.Radius * f.Bezier.Float64(0, 1, norm)
	}
}

type LineFunc struct {
	Slope     float64
	Intercept float64 // The Y intercept
}

func (f LineFunc) Interp(t float64) float64 {
	return f.Slope*t + f.Intercept
}

type AddFunc struct {
	A, B Function
}

func (f AddFunc) Interp(t float64) float64 {
	return f.A.Interp(t) + f.B.Interp(t)
}

type MultFunc struct {
	A, B Function
}

func (f MultFunc) Interp(t float64) float64 {
	return f.A.Interp(t) * f.B.Interp(t)
}
