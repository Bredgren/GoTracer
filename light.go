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
	Scene *Scene
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
	dir := p.Position.Sub(point)
	dist := dir.Len()
	dir = dir.Normalize()

	atten := Color64{1, 1, 1}
	distTraveled := 0.0

	shadowRay := NewRay(point, dir)
	isect, found := p.Scene.Intersect(shadowRay)
	for found && distTraveled + isect.T < dist {
		distTraveled += isect.T
		material := p.Scene.Material[isect.Object.GetMaterialName()]
		atten = atten.Product(material.Transmissive)

		shadowRay = NewRay(shadowRay.At(isect.T), dir)
		isect, found = p.Scene.Intersect(shadowRay)
	}

	return mgl64.Vec3(atten)
}
