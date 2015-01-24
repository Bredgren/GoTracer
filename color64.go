package raytracer

import (
	"image/color"

	"github.com/go-gl/mathgl/mgl64"
)

type Color64 mgl64.Vec3

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

func (c Color64) Product(other Color64) Color64 {
	return Color64{c[0] * other[0], c[1] * other[1], c[2] * other[2]}
}
