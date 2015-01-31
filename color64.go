package raytracer

import (
	"image/color"
	"math"

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
		uint8(Clamp01(c.R()) * 255),
		uint8(Clamp01(c.G()) * 255),
		uint8(Clamp01(c.B()) * 255),
		255,
	}
}

func (c Color64) Product(other Color64) Color64 {
	return Color64{c[0] * other[0], c[1] * other[1], c[2] * other[2]}
}

func clamp(low, high float64) func(float64) float64 {
	return func(v float64) float64 {
		return math.Min(high, math.Max(low, v))
	}
}

var Clamp01 = clamp(0.0, 1.0)

func (c Color64) Clamp() Color64 {
	return Color64{Clamp01(c[0]), Clamp01(c[1]), Clamp01(c[2])}
}

func ColorsDifferent(color1, color2 Color64, thresh float64) bool {
	return mgl64.Vec3(color1).Sub(mgl64.Vec3(color2)).Len() > thresh
}
