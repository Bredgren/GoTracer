package trace

import (
	"image/color"
	"math/rand"
)

// Scene represents the scene and all data needed to render it.
type Scene struct {
}

// ColorAt returns the color of the pixel at (x, y).
func (s *Scene) ColorAt(x, y int) color.NRGBA {
	return color.NRGBA{
		R: uint8(rand.Intn(255)),
		G: uint8(rand.Intn(255)),
		B: uint8(rand.Intn(255)),
		A: 255,
	}
}

// NewScene creates and returns a new Scene from the given options.
func NewScene(options *Options) *Scene {
	return &Scene{}
}
