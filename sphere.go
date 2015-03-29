package gotracer

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

// Sphere represents a sphere of radius 1 located at the origin.
type Sphere struct {
	*Object
}

func (s *Sphere) GetObject() *Object {
	return s.Object
}

// Intersect calculates the intersection between a ray and a unit sphere at
// the origin.
func (s *Sphere) Intersect(r *Ray, isect *Intersection) (hit bool) {
	// -(d . o) +- sqrt((d . o)^2 - (d . d)((o . o) - 1)) / (d . d)
	do := r.Dir.Dot(r.Origin)
	dd := r.Dir.Dot(r.Dir)
	oo := r.Origin.Dot(r.Origin)

	discriminant := do*do - dd*(oo-1)
	if discriminant < 0 {
		return false
	}

	discriminant = math.Sqrt(discriminant)

	t2 := (-do + discriminant) / dd
	if t2 <= Rayε {
		return false
	}

	t1 := (-do - discriminant) / dd
	if t1 > Rayε {
		isect.T = t1
		// No need to normalize because it's a unit sphere at the origin
		isect.Normal = r.At(t1)
		isect.Material = s.Object.Material
		// TODO: possible optimization would be to check if we even need uv coordinates
		u := 0.5 + (math.Atan2(isect.Normal.Y(), isect.Normal.X()) / (2 * math.Pi))
		v := 0.5 - (math.Asin(isect.Normal.Z()) / math.Pi)
		isect.UVCoords = mgl64.Vec2{u, v}
		return true
	}

	if t2 > Rayε {
		isect.T = t2
		isect.Normal = r.At(t2)
		isect.Material = s.Object.Material
		u := 0.5 + (math.Atan2(isect.Normal.Y(), isect.Normal.X()) / (2 * math.Pi))
		v := 0.5 - (math.Asin(isect.Normal.Z()) / math.Pi)
		isect.UVCoords = mgl64.Vec2{u, v}
		return true
	}

	return false
}
