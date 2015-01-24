package raytracer

import (
	"github.com/go-gl/mathgl/mgl64"
)

type SceneObject interface {
	Material() string
	Intersect(Ray) (Intersection, bool)
}

type SphereObject struct {
	MaterialName string
}

func (s SphereObject) Material() string {
	return s.MaterialName
}

func (s SphereObject) Intersect(r Ray) (Intersection, bool) {
	return Intersection{}, false
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
