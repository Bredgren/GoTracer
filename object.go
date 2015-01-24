package raytracer

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

type SceneObject interface {
	GetTransform() mgl64.Mat4
	GetInvTransform() mgl64.Mat4
	GetMaterialName() string
	// Intersect takes a ray in local coordinates and returns an intersection and true
	// if the intersects the object, false otherwise.
	Intersect(r Ray) (isect Intersection, hit bool)
}

type SphereObject struct {
	Transform mgl64.Mat4
	InvTransform mgl64.Mat4
	MaterialName string
}

func (s SphereObject) GetTransform() mgl64.Mat4 {
	return s.Transform
}

func (s SphereObject) GetInvTransform() mgl64.Mat4 {
	return s.InvTransform
}

func (s SphereObject) GetMaterialName() string {
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

	t1 := (-dp - discriminant) / dd
	if t1 > Rayε {
		isect.T = t1
		// Normalize because sphere is at origin
		isect.Normal = r.At(t1).Normalize()
		// TODO: set UV coordinates
		return isect, true
	}

	if t2 > Rayε {
		isect.T = t2
		// Normalize because sphere is at origin
		isect.Normal = r.At(t2).Normalize()
		// TODO: set UV coordinates
		return isect, true
	}

	return isect, false
}

type TriangleObject struct {
	Transform mgl64.Mat4
	InvTransform mgl64.Mat4
	MaterialName string
	PointA mgl64.Vec3
	PointB mgl64.Vec3
	PointC mgl64.Vec3
}

func (t TriangleObject) GetTransform() mgl64.Mat4 {
	return t.Transform
}

func (t TriangleObject) GetInvTransform() mgl64.Mat4 {
	return t.InvTransform
}

func (t TriangleObject) GetMaterialName() string {
	return t.MaterialName
}

func (t TriangleObject) Intersect(r Ray) (Intersection, bool) {
	return Intersection{}, false
}
