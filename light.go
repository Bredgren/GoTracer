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
	Color Color64
	Position mgl64.Vec3
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

const (
	DirectionalLightDist = 1e10
)
type DirectionalLight struct {
	Scene *Scene
	Color Color64
	Orientation mgl64.Vec3
}

func NewDirectionalLight(scene *Scene, color Color64, orient mgl64.Vec3) DirectionalLight {
	return DirectionalLight{scene, color, orient.Normalize()}
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
	return ShadowAttenuation(d.Scene, d.Direction(point).Mul(DirectionalLightDist), point)
}

type SpotLight struct {
	Scene *Scene
	Color Color64
	Position mgl64.Vec3
	Orientation mgl64.Vec3
	Angle float64
	DropOff float64
	FadeAngle float64

	angleDegrees float64
	minAngle float64
	maxAngle float64
	edgeIntensity float64
}

func NewSpotLight(scene *Scene, color Color64, pos, orient mgl64.Vec3, angle, dropOff, fadeAngle float64) SpotLight {
	angleDegrees := angle
	angle = mgl64.DegToRad(angle)
	fadeAngle = mgl64.DegToRad(fadeAngle)
	minAngle := angle - fadeAngle
	maxAngle := angle + fadeAngle
	edgeIntensity := math.Pow(math.Cos(minAngle), dropOff)
	return SpotLight{scene, color, pos, orient.Normalize(), angle, dropOff, fadeAngle, angleDegrees, minAngle, maxAngle, edgeIntensity}
}

func (s SpotLight) GetColor() Color64 {
	return s.Color
}

func (s SpotLight) Direction(point mgl64.Vec3) mgl64.Vec3 {
	return s.Orientation.Mul(-1)
}

func (s SpotLight) DistanceAttenuation(point mgl64.Vec3) float64 {
	atten := 0.0
	dir := point.Sub(s.Position).Normalize()
	cosθ := dir.Dot(s.Orientation)
	θ := math.Acos(cosθ)
	if θ > 0 && θ < s.minAngle {
		atten = math.Pow(cosθ, s.DropOff)
	} else if θ >= s.minAngle && θ < s.maxAngle {
		diff := (s.maxAngle - θ) / (s.maxAngle - s.minAngle)
		if s.DropOff / s.angleDegrees < 0.5 {
			if diff < 0.5 {
				diff *= 2.0
				diff *= diff
				diff /= 2.0
			} else {
				diff = 2 * (1 - diff)
				diff *= diff
				diff = 1 - (diff / 2)
			}
		}
		atten = s.edgeIntensity * diff
	}
	return atten
}

func (s SpotLight) ShadowAttenuation(point mgl64.Vec3) mgl64.Vec3 {
	return ShadowAttenuation(s.Scene, s.Position.Sub(point), point)
}
