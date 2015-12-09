package trace

import (
	"image/color"

	"github.com/go-gl/mathgl/mgl64"
)

// Color64 represents an RGB color with each compenent a flaot64.
type Color64 [3]float64

// R returns the red component.
func (c Color64) R() float64 {
	return c[0]
}

// G returns the green component.
func (c Color64) G() float64 {
	return c[1]
}

// B returns the blue component.
func (c Color64) B() float64 {
	return c[2]
}

// NRGBA returns the color.NRGBA equivalent of the Color64. The components are clamped between
// 0.0 and 1.0 and treated as a percent of 255.
func (c Color64) NRGBA() color.NRGBA {
	return color.NRGBA{
		uint8(mgl64.Clamp(c.R(), 0, 1) * 255),
		uint8(mgl64.Clamp(c.G(), 0, 1) * 255),
		uint8(mgl64.Clamp(c.B(), 0, 1) * 255),
		255,
	}
}

// ColorsDifferent returns true if the distance between the given colors is greater than thresh.
func ColorsDifferent(c1, c2 Color64, thresh float64) bool {
	return mgl64.Vec3(c1).Sub(mgl64.Vec3(c2)).Len() > thresh
}
