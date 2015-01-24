package raytracer

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

type Light interface {
	GetColor() Color64
	Direction(point mgl64.Vec3) mgl64.Vec3
	DistanceAttenuation(point mgl64.Vec3) float64
	ShadowAttenuation(point mgl64.Vec3) mgl64.Vec3
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
	r := p.Position.Sub(point).Len()
	return math.Min(1.0, 1.0 / (p.ConstCoeff + p.LinearCoeff * r + p.QuadCoeff * r * r))
}

func (p PointLight) ShadowAttenuation(point mgl64.Vec3) mgl64.Vec3 {
	// TODO: add attuation
	return mgl64.Vec3{1.0, 1.0, 1.0}
}
