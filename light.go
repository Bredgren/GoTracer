package raytracer

import (
	"github.com/go-gl/mathgl/mgl64"
)

type Light interface {
	GetColor() Color64
	Direction(point mgl64.Vec3) mgl64.Vec3
	DistanceAttenuation(point mgl64.Vec3) float64
	// TODO: ShadowAttenuation(point mgl64.Vec3) mgl64.Vec3
}

type PointLight struct {
	Position mgl64.Vec3
	Color Color64
	ConstCoeff float64
	LinearCoeff float64
	QuadCoeff float64
}

func (p PointLight) GetColor() Color64 {
	return p.Color
}

func (p PointLight) Direction(point mgl64.Vec3) mgl64.Vec3 {
	return p.Position.Sub(point).Normalize()
}

func (p PointLight) DistanceAttenuation(point mgl64.Vec3) float64 {
	// TODO: add attuation
	return 1.0
}
