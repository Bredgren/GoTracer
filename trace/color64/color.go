package color64

import (
	"image/color"
	"math"

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

// Add does a component-wise addition.
func (c Color64) Add(other Color64) Color64 {
	return Color64{c[0] + other[0], c[1] + other[1], c[2] + other[2]}
}

// Mul multiplies each component by val.
func (c Color64) Mul(val float64) Color64 {
	return Color64{c[0] * val, c[1] * val, c[2] * val}
}

// Different returns true if the distance between the given colors is greater than thresh.
func Different(c1, c2 Color64, thresh float64) bool {
	r, g, b := c1[0]-c2[0], c1[1]-c2[1], c1[2]-c2[2]
	return r*r+g*g+b*b > thresh*thresh
}

// Different implementations for benchmarking
func different1(c1, c2 Color64, thresh float64) bool {
	return mgl64.Vec3(c1).Sub(mgl64.Vec3(c2)).Len() > thresh
}

func different2(c1, c2 Color64, thresh float64) bool {
	r, g, b := c1[0]-c2[0], c1[1]-c2[1], c1[2]-c2[2]
	return math.Hypot(math.Hypot(r, g), b) > thresh
}

func different3(c1, c2 Color64, thresh float64) bool {
	r, g, b := c1[0]-c2[0], c1[1]-c2[1], c1[2]-c2[2]
	return math.Sqrt(r*r+g*g+b*b) > thresh
}

func different4(c1, c2 Color64, thresh float64) bool {
	r, g, b := c1[0]-c2[0], c1[1]-c2[1], c1[2]-c2[2]
	return r*r+g*g+b*b > thresh*thresh
}
