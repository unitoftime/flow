package interp

import (
	"math"
	"encoding/gob"

	"github.com/ungerik/go3d/float64/bezier2"
	"github.com/ungerik/go3d/float64/vec2"
)

// TODO: use https://easings.net/

func init() {
	gob.Register(Lerp{})
	gob.Register(Bezier{})
	gob.Register(Equation{})
	gob.Register(SinFunc{})
}

var Linear *Lerp = &Lerp{}

var EaseOut *Bezier = &Bezier{
	bezier2.T{
		vec2.T{0.0, 0.0},
		vec2.T{1.0, 0.0},
		vec2.T{1.0, 0.0},
		vec2.T{1.0, 1.0},
	},
}
var EaseIn *Bezier = &Bezier{
	bezier2.T{
		vec2.T{0.0, 0.0},
		vec2.T{0.0, 1.0},
		vec2.T{0.0, 1.0},
		vec2.T{1.0, 1.0},
	},
}
var EaseInOut *Bezier = &Bezier{
	bezier2.T{
		vec2.T{0.0, 0.0},
		vec2.T{1.0, 0.0},
		vec2.T{0.0, 1.0},
		vec2.T{1.0, 1.0},
	},
}
var Sinusoid *Equation = &Equation{
	Func: SinFunc{},
}

//https://cubic-bezier.com/#.22,1,.36,1
var EaseOutQuint *Bezier = &Bezier{
	bezier2.T{
		vec2.T{0.0, 0.0},
		vec2.T{0.22, 1.0},
		vec2.T{0.36, 1.0},
		vec2.T{1.0, 1.0},
	},
}

var EaseTest *Bezier = &Bezier{
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
	Vec2(vec2.T, vec2.T, float64) vec2.T
}

type Lerp struct {
}
func (i *Lerp) Float64(a, b float64, t float64) float64 {
	m := b - a // Slope = Rise over run | Note: Run = (1 - 0)
	y := (m * t) + a
	return y
}
func (i *Lerp) Float32(a, b float32, t float64) float32 {
	m := b - a // Slope = Rise over run | Note: Run = (1 - 0)
	y := (m * float32(t)) + a
	return y
}
func (i *Lerp) Uint8(a, b uint8, t float64) uint8 {
	return uint8(i.Float64(float64(a), float64(b), t))
}

func (i *Lerp) Vec2(a, b vec2.T, t float64) vec2.T {
	ret := vec2.T{
		i.Float64(a[0], b[0], t),
		i.Float64(a[1], b[1], t),
	}
	return ret
}

type Bezier struct {
	bezier2.T
}
func (i *Bezier) Float64(a, b float64, t float64) float64 {
	iValue := i.T.Point(t)
	return Linear.Float64(a, b, iValue[1])
}
func (i *Bezier) Float32(a, b float32, t float64) float32 {
	iValue := i.T.Point(t)
	return Linear.Float32(a, b, iValue[1])
}
func (i *Bezier) Uint8(a, b uint8, t float64) uint8 {
	iValue := i.T.Point(t)
	return Linear.Uint8(a, b, iValue[1])
}
func (i *Bezier) Vec2(a, b vec2.T, t float64) vec2.T {
	iValue := i.T.Point(t)
	return Linear.Vec2(a, b, iValue[1])
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
func (i *Equation) Vec2(a, b vec2.T, t float64) vec2.T {
	iValue := i.Func.Interp(t)
	return Linear.Vec2(a, b, iValue)
}

type Function interface {
	Interp(t float64) float64
}

type SinFunc struct {}
func (s SinFunc) Interp(t float64) float64 {
	return math.Sin(t * math.Pi)
}
