package raytracer

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

type SceneObject interface {
	Material() string
	// Intersect takes a ray in local coordinates and returns an intersection and true
	// if the intersects the object, false otherwise.
	Intersect(r Ray) (isect Intersection, hit bool)
}

type SphereObject struct {
	MaterialName string
}

func (s SphereObject) Material() string {
	return s.MaterialName
}

func (s SphereObject) Intersect(r Ray) (isect Intersection, hit bool) {
	isect = Intersection{Object: s}

	// -(d . o) +- sqrt((d . o)^2 - (d . d)((o . o) - 1)) / (d . d)
	dp := r.Direction.Dot(r.Origin)
	dd := r.Direction.Dot(r.Direction)
	pp := r.Origin.Dot(r.Origin)

	discriminant := dp * dp - dd * (pp - 1)
	if discriminant < 0 {
		return isect, false
	}

	discriminant = math.Sqrt(discriminant)

	t2 := (-dp + discriminant) / dd
	if t2 <= Rayε {
		return isect, false
	}

	t1 := (-dp + discriminant) / dd
	if t1 > Rayε {
		p := r.At(t1)
		// Normalize because sphere is at origin
		p = p.Normalize()
		isect.T = t1
		isect.Normal = p
		// TODO: set UV coordinates
		return isect, true
	}

	if t2 > Rayε {
		p := r.At(t2)
		p = p.Normalize()
		isect.T = t2
		isect.Normal = p
		// TODO: set UV coordinates
		return isect, true
	}

	return isect, false
}

type TriangleObject struct {
	PointA mgl64.Vec3
	PointB mgl64.Vec3
	PointC mgl64.Vec3
	MaterialName string
}

func (t TriangleObject) Material() string {
	return t.MaterialName
}

func (t TriangleObject) Intersect(r Ray) (Intersection, bool) {
	return Intersection{}, false
}
