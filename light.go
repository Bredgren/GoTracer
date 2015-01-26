package raytracer

import (
	"math"
	"math/rand"

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

func PointDistAtten(r, constCoeff, linearCoeff, quadCoeff float64) float64 {
	return math.Min(1.0, 1.0 / (constCoeff + linearCoeff * r + quadCoeff * r * r))
}

func (p PointLight) DistanceAttenuation(point mgl64.Vec3) float64 {
	r := p.Position.Sub(point).Len()
	return PointDistAtten(r, p.ConstCoeff, p.LinearCoeff, p.QuadCoeff)
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


type AreaLight struct {
	Scene *Scene
	Color Color64
	Position mgl64.Vec3
	Orientation mgl64.Vec3
	UpDir mgl64.Vec3
	Size float64
	Samples int
	ConstCoeff float64
	LinearCoeff float64
	QuadCoeff float64

	u mgl64.Vec3
	v mgl64.Vec3
}

// func NewAreaLight(scene *Scene, color Color64, pos, orient, upDir mgl64.Vec3, size float64, samples int) AreaLight {
func NewAreaLight(a AreaLight) (light AreaLight) {
	light = a
	light.Orientation = light.Orientation.Normalize()
	light.UpDir = light.UpDir.Normalize()

	z := a.Orientation.Normalize()
	y := light.UpDir
	x := y.Cross(z).Normalize()
	y = z.Cross(x).Normalize()
	m := mgl64.Mat3FromCols(x, y, z)
	light.u = m.Mul3x1(mgl64.Vec3{1, 0, 0}.Mul(light.Size))
	light.v = m.Mul3x1(mgl64.Vec3{0, -1, 0}.Mul(light.Size))

	return light
}

func (a AreaLight) GetColor() Color64 {
	return a.Color
}

func (a AreaLight) Direction(point mgl64.Vec3) mgl64.Vec3 {
	return a.Position.Sub(point).Normalize()
}

func (a AreaLight) DistanceAttenuation(point mgl64.Vec3) float64 {
	r := a.Position.Sub(point).Len()
	return PointDistAtten(r, a.ConstCoeff, a.LinearCoeff, a.QuadCoeff)
}

func (a AreaLight) GridAttenuation(point mgl64.Vec3, samples int) mgl64.Vec3 {
	atten := mgl64.Vec3{}
	sizePerSample := float64(a.Size) / float64(samples)
	for y := 0.0; y < 1.0; y += sizePerSample {
		for x := 0.0; x < 1.0; x += sizePerSample {
			rx, ry := rand.Float64() / sizePerSample, rand.Float64() / sizePerSample
			pos := a.GetPosition(x + rx, y + ry)
			atten = atten.Add(ShadowAttenuation(a.Scene, pos.Sub(point), point))
		}
	}
	return atten.Mul(1.0 / float64(a.Samples * a.Samples))
}

func (a AreaLight) ShadowAttenuation(point mgl64.Vec3) mgl64.Vec3 {
	// Check with a 2x2 grid to see there is a significant difference before going all out
	// center := ShadowAttenuation(a.Scene, a.Position.Sub(point), point)
	// atten4 := a.GridAttenuation(point, 2)
	// if center.ApproxEqualThreshold(atten4, Rayε / 2) {
	// 	return atten4
	// }
	return a.GridAttenuation(point, a.Samples)
}

// GetPosition takes normalized "light" coordinates and returns the point in space
// that cooresponds to that point. This works the same way as the Camera.
func (a AreaLight) GetPosition(nx, ny float64) mgl64.Vec3 {
	return a.Position.Add(a.u.Mul(nx - 0.5)).Add(a.v.Mul(ny - 0.5))
}
