package raytracer

import (
	"image/color"
)

type Color64 [3]float64

func (c Color64) R() float64 {
	return c[0]
}

func (c Color64) G() float64 {
	return c[1]
}

func (c Color64) B() float64 {
	return c[2]
}

func (c Color64) NRGBA() color.NRGBA {
	return color.NRGBA{
		uint8(c.R() * 255),
		uint8(c.G() * 255),
		uint8(c.B() * 255),
		255,
	}
}
