package glm

import (
	"image/color"
	"math"
)

var (
	White = RGBA{1, 1, 1, 1}
	Black = RGBA{0, 0, 0, 1}
)

// Premultipled RGBA value scaled from [0, 1.0]
type RGBA struct {
	R, G, B, A float64
}

// TODO - conversion from golang colors
func FromUint8(r, g, b, a uint8) RGBA {
	return RGBA{
		float64(r) / float64(math.MaxUint8),
		float64(g) / float64(math.MaxUint8),
		float64(b) / float64(math.MaxUint8),
		float64(a) / float64(math.MaxUint8),
	}
}

func (c RGBA) ToUint8() color.NRGBA {
	return color.NRGBA{
		R: uint8(math.Round(c.R * float64(math.MaxUint8))),
		G: uint8(math.Round(c.G * float64(math.MaxUint8))),
		B: uint8(math.Round(c.B * float64(math.MaxUint8))),
		A: uint8(math.Round(c.A * float64(math.MaxUint8))),
	}
}

func HexColor(col uint64, alpha uint8) RGBA {
	return FromNRGBA(color.NRGBA{
		R: uint8((col >> 16) & 0xff),
		G: uint8((col >> 8) & 0xff),
		B: uint8(col & 0xff),
		A: alpha,
	})
}

func Alpha(a float64) RGBA {
	return RGBA{a, a, a, a}
}

func Greyscale(g float64) RGBA {
	return RGBA{g, g, g, 1.0}
}

func FromStraightRGBA(r, g, b float64, a float64) RGBA {
	return RGBA{r * a, g * a, b * a, a}
}

func FromNRGBA(c color.NRGBA) RGBA {
	r, g, b, a := c.RGBA()

	return RGBA{
		float64(r) / float64(math.MaxUint16),
		float64(g) / float64(math.MaxUint16),
		float64(b) / float64(math.MaxUint16),
		float64(a) / float64(math.MaxUint16),
	}
}

func FromRGBA(c color.RGBA) RGBA {
	return FromUint8(c.R, c.G, c.B, c.A)
}

func FromColor(c color.Color) RGBA {
	r, g, b, a := c.RGBA()

	return RGBA{
		float64(r) / float64(math.MaxUint16),
		float64(g) / float64(math.MaxUint16),
		float64(b) / float64(math.MaxUint16),
		float64(a) / float64(math.MaxUint16),
	}
}
func (c1 RGBA) Mult(c2 RGBA) RGBA {
	return RGBA{
		c1.R * c2.R,
		c1.G * c2.G,
		c1.B * c2.B,
		c1.A * c2.A,
	}
}

func (c1 RGBA) Add(c2 RGBA) RGBA {
	return RGBA{
		c1.R + c2.R,
		c1.G + c2.G,
		c1.B + c2.B,
		c1.A + c2.A,
	}
}

func (c RGBA) Desaturate(val float64) RGBA {
	// https://stackoverflow.com/questions/70966873/algorithm-to-desaturate-rgb-color
	i := (c.R + c.G + c.B) / 3

	dr := i - c.R
	dg := i - c.G
	db := i - c.B

	return RGBA{
		c.R + (dr * val),
		c.G + (dg * val),
		c.B + (db * val),
		c.A,
	}
}
