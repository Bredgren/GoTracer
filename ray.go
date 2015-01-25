package raytracer

import (
	"github.com/go-gl/mathgl/mgl64"
)

const (
	Rayε = 0.00001
)

type Ray struct {
	Origin mgl64.Vec3
	Direction mgl64.Vec3
}

func NewRay(origin, direction mgl64.Vec3) Ray {
	return Ray{origin, direction.Normalize()}
}

// At returns the point marked by the ray at t
func (r Ray) At(t float64) mgl64.Vec3 {
	return r.Origin.Add(r.Direction.Mul(t))
}

func TransformVec3(m mgl64.Mat4, v mgl64.Vec3) mgl64.Vec3 {
	return m.Mul4x1(v.Vec4(1)).Vec3()
}

func (r Ray) Transform(transform mgl64.Mat4) (newR Ray, len float64) {
	newOrigin := TransformVec3(transform, r.Origin)
	newDir := TransformVec3(transform, r.Origin.Add(r.Direction)).Sub(newOrigin)
	len = newDir.Len()
	return NewRay(newOrigin, newDir), len
}

func (r Ray) Reflect(isect Intersection) (reflRay Ray) {
	minusD := r.Direction.Mul(-1)
	nMD2 := isect.Normal.Dot(minusD) * 2
	reflDir := isect.Normal.Mul(nMD2).Add(r.Direction)
	return NewRay(r.At(isect.T), reflDir)
}

type Intersection struct {
	Object SceneObject
	Normal mgl64.Vec3
	T float64
}
