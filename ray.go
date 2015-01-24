package raytracer

import (
	"github.com/go-gl/mathgl/mgl64"
)

const (
	Rayε = 0.00001
)

type Ray struct {
	Origin mgl64.Vec3
	Direction mgl64.Vec3
}

// At returns the point marked by the ray at t
func (r Ray) At(t float64) mgl64.Vec3 {
	return r.Origin.Add(r.Direction.Mul(t))
}

type Intersection struct {
	Object SceneObject
	Normal mgl64.Vec3
	T float64
}
