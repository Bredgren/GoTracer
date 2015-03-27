package gotracer

import (
	// "math"
	// "math/rand"

	"github.com/go-gl/mathgl/mgl64"
)

type Light interface {
	Attenuation(scene *Scene, point mgl64.Vec3) Color64
	Direction(from mgl64.Vec3) mgl64.Vec3
}

// ShadowAttenuation takes the unnormalized direction to the light so we know when
// we've "hit" it.
func ShadowAttenuation(scene *Scene, dir mgl64.Vec3, point mgl64.Vec3) Color64 {
	var dist float64 = dir.Len()
	dir = dir.Normalize()

	atten := Color64{1, 1, 1}
	distTraveled := 0.0

	shadowRay := Ray{ShadowRay, point, dir}
	isect := Intersection{}
	hit := scene.Intersect(&shadowRay, &isect)
	for hit && atten.Len2() > Rayε && distTraveled+isect.T < dist {
		distTraveled += isect.T
		material := isect.Material
		if isect.Normal.Dot(dir) > 0 {
			// Exiting object
			// atten = atten.Product(material.BeersTrans(isect))
		} else {
			// Entering object
			atten = atten.Product(material.Transmissive.ColorAt(isect.UVCoords))
		}

		shadowRay.Origin = shadowRay.At(isect.T)
		hit = scene.Intersect(&shadowRay, &isect)
	}

	return atten
}

// type Light interface {
// 	GetColor() Color64
// 	Direction(point mgl64.Vec3) mgl64.Vec3
// 	DistanceAttenuation(point mgl64.Vec3) float64
// 	ShadowAttenuation(point mgl64.Vec3) mgl64.Vec3
// 	Attenuation(point mgl64.Vec3) mgl64.Vec3
// }

// type PointLight struct {
// 	Scene *Scene
// 	Color Color64
// 	Position mgl64.Vec3
// 	ConstCoeff float64
// 	LinearCoeff float64
// 	QuadCoeff float64
// }

// func InitPointLight(p *PointLight) {
// }

// func (p PointLight) GetColor() Color64 {
// 	return p.Color
// }

// func (p PointLight) Direction(point mgl64.Vec3) mgl64.Vec3 {
// 	return p.Position.Sub(point).Normalize()
// }

// func PointDistAtten(r, constCoeff, linearCoeff, quadCoeff float64) float64 {
// 	return math.Min(1.0, 1.0 / (constCoeff + linearCoeff * r + quadCoeff * r * r))
// }

// func (p PointLight) DistanceAttenuation(point mgl64.Vec3) float64 {
// 	r := p.Position.Sub(point).Len()
// 	return PointDistAtten(r, p.ConstCoeff, p.LinearCoeff, p.QuadCoeff)
// }

// func (p PointLight) ShadowAttenuation(point mgl64.Vec3) mgl64.Vec3 {
// 	return ShadowAttenuation(p.Scene, p.Position.Sub(point), point)
// }

// const (
// 	DirectionalLightDist = 1e10
// )
// type DirectionalLight struct {
// 	Scene *Scene
// 	Color Color64
// 	Orientation mgl64.Vec3
// }

// func InitDirectionalLight(d *DirectionalLight) {
// 	d.Orientation = d.Orientation.Normalize()
// }

// func (d DirectionalLight) GetColor() Color64 {
// 	return d.Color
// }

// func (d DirectionalLight) Direction(point mgl64.Vec3) mgl64.Vec3 {
// 	return d.Orientation.Mul(-1)
// }

// func (d DirectionalLight) DistanceAttenuation(point mgl64.Vec3) float64 {
// 	return 1.0
// }

// func (d DirectionalLight) ShadowAttenuation(point mgl64.Vec3) mgl64.Vec3 {
// 	return ShadowAttenuation(d.Scene, d.Direction(point).Mul(DirectionalLightDist), point)
// }

// type SpotLight struct {
// 	Scene *Scene
// 	Color Color64
// 	Position mgl64.Vec3
// 	Orientation mgl64.Vec3
// 	Angle float64
// 	DropOff float64
// 	FadeAngle float64

// 	angleDegrees float64
// 	minAngle float64
// 	maxAngle float64
// 	edgeIntensity float64
// }

// func InitSpotLight(s *SpotLight) {
// 	s.angleDegrees = s.Angle
// 	s.Angle = mgl64.DegToRad(s.Angle)
// 	s.FadeAngle = mgl64.DegToRad(s.FadeAngle)
// 	s.minAngle = s.Angle - s.FadeAngle
// 	s.maxAngle = s.Angle + s.FadeAngle
// 	s.edgeIntensity = math.Pow(math.Cos(s.minAngle), s.DropOff)
// 	s.Orientation = s.Orientation.Normalize()
// }

// func (s SpotLight) GetColor() Color64 {
// 	return s.Color
// }

// func (s SpotLight) Direction(point mgl64.Vec3) mgl64.Vec3 {
// 	return s.Orientation.Mul(-1)
// }

// func (s SpotLight) DistanceAttenuation(point mgl64.Vec3) float64 {
// 	atten := 0.0
// 	dir := point.Sub(s.Position).Normalize()
// 	cosθ := dir.Dot(s.Orientation)
// 	θ := math.Acos(cosθ)
// 	if θ > 0 && θ < s.minAngle {
// 		atten = math.Pow(cosθ, s.DropOff)
// 	} else if θ >= s.minAngle && θ < s.maxAngle {
// 		diff := (s.maxAngle - θ) / (s.maxAngle - s.minAngle)
// 		if s.DropOff / s.angleDegrees < 0.5 {
// 			if diff < 0.5 {
// 				diff *= 2.0
// 				diff *= diff
// 				diff /= 2.0
// 			} else {
// 				diff = 2 * (1 - diff)
// 				diff *= diff
// 				diff = 1 - (diff / 2)
// 			}
// 		}
// 		atten = s.edgeIntensity * diff
// 	}
// 	return atten
// }

// func (s SpotLight) ShadowAttenuation(point mgl64.Vec3) mgl64.Vec3 {
// 	return ShadowAttenuation(s.Scene, s.Position.Sub(point), point)
// }

// type AreaLight struct {
// 	Scene *Scene
// 	Color Color64
// 	Position mgl64.Vec3
// 	Orientation mgl64.Vec3
// 	UpDir mgl64.Vec3
// 	Size float64
// 	Samples int
// 	ConstCoeff float64
// 	LinearCoeff float64
// 	QuadCoeff float64
// 	Accelerated bool

// 	u mgl64.Vec3
// 	v mgl64.Vec3
// }

// func InitAreaLight(a *AreaLight) {
// 	a.Orientation = a.Orientation.Normalize()
// 	a.UpDir = a.UpDir.Normalize()

// 	z := a.Orientation.Normalize()
// 	y := a.UpDir
// 	x := z.Cross(y).Normalize()
// 	y = z.Cross(x).Normalize()
// 	m := mgl64.Mat3FromCols(x, y, z)
// 	a.u = m.Mul3x1(mgl64.Vec3{1, 0, 0}.Mul(a.Size))
// 	a.v = m.Mul3x1(mgl64.Vec3{0, -1, 0}.Mul(a.Size))
// }

// func (a AreaLight) GetColor() Color64 {
// 	return a.Color
// }

// func (a AreaLight) Direction(point mgl64.Vec3) mgl64.Vec3 {
// 	return a.Position.Sub(point).Normalize()
// }

// func (a AreaLight) DistanceAttenuation(point mgl64.Vec3) float64 {
// 	r := a.Position.Sub(point).Len()
// 	return PointDistAtten(r, a.ConstCoeff, a.LinearCoeff, a.QuadCoeff)
// }

// func (a AreaLight) GridAttenuation(point mgl64.Vec3, samples int) mgl64.Vec3 {
// 	atten := mgl64.Vec3{}
// 	s := float64(samples)
// 	sizePerSample := 1.0 / s
// 	for y := 0.0; y < s; y++ {
// 		for x := 0.0; x < s; x++ {
// 			xx := x * sizePerSample
// 			yy := y * sizePerSample
// 			rx, ry := rand.Float64() * sizePerSample, rand.Float64() * sizePerSample
// 			pos := a.GetPosition(xx + rx, yy + ry)
// 			a := ShadowAttenuation(a.Scene, pos.Sub(point), point)
// 			atten = atten.Add(a)
// 		}
// 	}
// 	return atten.Mul(1.0 / (s * s))
// }

// func (a AreaLight) ShadowAttenuation(point mgl64.Vec3) mgl64.Vec3 {
// 	if a.Accelerated {
// 		// Check the corners to see there is a significant difference before going all out
// 		center := ShadowAttenuation(a.Scene, a.Position.Sub(point), point)
// 		atten1 := ShadowAttenuation(a.Scene, a.GetPosition(0.0, 0.0).Sub(point), point)
// 		atten2 := ShadowAttenuation(a.Scene, a.GetPosition(0.0, 1.0).Sub(point), point)
// 		atten3 := ShadowAttenuation(a.Scene, a.GetPosition(1.0, 0.0).Sub(point), point)
// 		atten4 := ShadowAttenuation(a.Scene, a.GetPosition(1.0, 1.0).Sub(point), point)
// 		min := math.Min(center.Len(), math.Min(atten1.Len(), math.Min(atten2.Len(), math.Min(atten3.Len(), atten4.Len()))))
// 		max := math.Max(center.Len(), math.Max(atten1.Len(), math.Max(atten2.Len(), math.Max(atten3.Len(), atten4.Len()))))
// 		if max - min < Rayε {
// 			return atten1.Add(atten2).Add(atten3).Add(atten4).Mul(0.25)
// 		}
// 	}
// 	return a.GridAttenuation(point, a.Samples)
// }

// // GetPosition takes normalized "light" coordinates and returns the point in space
// // that cooresponds to that point. This works the same way as the Camera.
// func (a AreaLight) GetPosition(nx, ny float64) mgl64.Vec3 {
// 	return a.Position.Add(a.u.Mul(nx - 0.5)).Add(a.v.Mul(ny - 0.5))
// }
