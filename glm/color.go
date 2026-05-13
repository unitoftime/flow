package glm

import (
	"image/color"
	"math"
)

var (
	White = RGBA{1, 1, 1, 1}
	Black = RGBA{0, 0, 0, 1}
	Transparent = RGBA{0, 0, 0, 0}
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

// Converts a non premultiplied hex uint64 and alpha into premultiplied
func HexColor(col uint64, alpha uint8) RGBA {
	return FromColor(color.RGBA{
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
	return FromColor(c)
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

func (c RGBA) Clamp() RGBA {
	return RGBA{
		R: Clamp(0, 1, c.R),
		G: Clamp(0, 1, c.G),
		B: Clamp(0, 1, c.B),
		A: Clamp(0, 1, c.A),
	}
}

func (c1 RGBA) Avg(c2 RGBA) RGBA {
	return c1.Add(c2).Mult(RGBA{0.5, 0.5, 0.5, 0.5})
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

// Assuming your package has an HSV struct. If not, here is a standard one:
type HSVA struct {
	H, S, V, A float64
}

// RGBAToHSV converts an RGBA color to HSVA (Hue, Saturation, Value, Alpha).
// All input and output values are in the range [0.0, 1.0].
func RGBAToHSV(c RGBA) HSVA {
	out := HSVA{A: c.A}

	min := math.Min(math.Min(c.R, c.G), c.B)
	max := math.Max(math.Max(c.R, c.G), c.B)
	delta := max - min

	out.V = max

	if max > 0.0 {
		out.S = delta / max
	} else {
		// If max is 0, the color is black. Saturation and Hue don't matter.
		out.S = 0.0
		out.H = 0.0 // Mathematically undefined, but we default to 0
		return out
	}

	if delta == 0.0 {
		// If delta is 0, the color is a shade of grey. Hue doesn't matter.
		out.H = 0.0
		return out
	}

	// Calculate Hue
	if c.R == max {
		out.H = (c.G - c.B) / delta
	} else if c.G == max {
		out.H = 2.0 + (c.B - c.R)/delta
	} else {
		out.H = 4.0 + (c.R - c.G)/delta
	}

	// Convert hue to [0, 1] range (from the standard 6-sector calculation)
	out.H /= 6.0
	if out.H < 0.0 {
		out.H += 1.0
	}

	return out
}

// HSVToRGBA converts an HSVA color back to RGBA.
// All input and output values are in the range [0.0, 1.0].
func HSVToRGBA(c HSVA) RGBA {
	out := RGBA{A: c.A}

	if c.S <= 0.0 {
		// If saturation is 0, the color is a shade of grey
		out.R = c.V
		out.G = c.V
		out.B = c.V
		return out
	}

	// Wrap hue to ensure it's in the [0, 1] bounds, then scale to 6 sectors
	h := c.H
	for h < 0.0 {
		h += 1.0
	}
	for h >= 1.0 {
		h -= 1.0
	}
	h *= 6.0

	sector := math.Floor(h)
	fract := h - sector

	p := c.V * (1.0 - c.S)
	q := c.V * (1.0 - (c.S * fract))
	t := c.V * (1.0 - (c.S * (1.0 - fract)))

	switch int(sector) {
	case 0:
		out.R, out.G, out.B = c.V, t, p
	case 1:
		out.R, out.G, out.B = q, c.V, p
	case 2:
		out.R, out.G, out.B = p, c.V, t
	case 3:
		out.R, out.G, out.B = p, q, c.V
	case 4:
		out.R, out.G, out.B = t, p, c.V
	default: // case 5
		out.R, out.G, out.B = c.V, p, q
	}

	return out
}

