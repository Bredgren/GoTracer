package gotracer

import (
	"github.com/go-gl/mathgl/mgl64"
)

const (
	AirIndex = 1.0003
)

// NOTES:
// Reflective
//   - determines how much diffuse vs reflected light
//   - base for fresnel
// Roughness (gloss/smoothness)
//   - be sure to conserve energy. Light will appear dimmer when spread out
//   - fresnel less noticable with nore roughness

// MaterialAttribute is an interface for getting a color. Designed with texture
// mapping in mind.
type MaterialAttribute interface {
	// ColorAt takes normalized coordinates and return the color at those coordinates.
	ColorAt(mgl64.Vec2) Color64
}

type Material struct {
	Name         string
	Emissive     MaterialAttribute
	Ambient      MaterialAttribute
	Diffuse      MaterialAttribute
	Specular     MaterialAttribute
	Reflective   MaterialAttribute
	Smoothness   MaterialAttribute
	Transmissive MaterialAttribute
	Index        MaterialAttribute
	Normal       MaterialAttribute
	IsLiquid     MaterialAttribute
	BRDF         BRDF
}

// UniformColor represents a MaterialAttribute that is the same at every location.
type UniformColor struct {
	Color Color64
}

func (c UniformColor) ColorAt(_ mgl64.Vec2) Color64 {
	return c.Color
}
