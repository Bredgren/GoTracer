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

// ShadowAttenuation takes the unnormalized direction to the light so we know when
// we've "hit" it.
func ShadowAttenuation(scene *Scene, dir mgl64.Vec3, point mgl64.Vec3) mgl64.Vec3 {
	dist := dir.Len()
	dir = dir.Normalize()

	atten := Color64{1, 1, 1}
	distTraveled := 0.0

	shadowRay := NewRay(point, dir)
	isect, found := scene.Intersect(shadowRay)
	for found && distTraveled + isect.T < dist {
		distTraveled += isect.T
		material := scene.Material[isect.Object.GetMaterialName()]
		atten = atten.Product(material.Transmissive)

		shadowRay = NewRay(shadowRay.At(isect.T), dir)
		isect, found = scene.Intersect(shadowRay)
	}

	return mgl64.Vec3(atten)
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
	return ShadowAttenuation(p.Scene, p.Position.Sub(point), point)
}

type DirectionalLight struct {
	Scene *Scene
	Orientation mgl64.Vec3
	Color Color64
}

func NewDirectionalLight(scene *Scene, orient mgl64.Vec3, color Color64) DirectionalLight {
	return DirectionalLight{scene, orient.Normalize(), color}
}

func (d DirectionalLight) GetColor() Color64 {
	return d.Color
}

func (d DirectionalLight) Direction(point mgl64.Vec3) mgl64.Vec3 {
	return d.Orientation.Mul(-1)
}

func (d DirectionalLight) DistanceAttenuation(point mgl64.Vec3) float64 {
	return 1.0
}

func (d DirectionalLight) ShadowAttenuation(point mgl64.Vec3) mgl64.Vec3 {
	return ShadowAttenuation(d.Scene, d.Direction(point).Mul(100), point)
}
