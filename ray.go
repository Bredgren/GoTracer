package gotracer

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

const (
	Rayε = 0.00001
)

const (
	PrimaryRay = iota
	CollisionRay
	ReflectionRay
	RefractionRay
	ShadowRay
	NumRayTypes
)

type Ray struct {
	Type      int
	Origin    mgl64.Vec3
	Direction mgl64.Vec3
	InvDir    mgl64.Vec3
}

func NewRay(t int, origin, direction mgl64.Vec3) Ray {
	dir := direction.Normalize()
	return Ray{t, origin, dir, mgl64.Vec3{1.0 / dir.X(), 1.0 / dir.Y(), 1.0 / dir.Z()}}
}

// At returns the point marked by the ray at t
func (r Ray) At(t float64) mgl64.Vec3 {
	return r.Origin.Add(r.Direction.Mul(t))
}

// TransformVec3 is used to transform a vector by a matrix.
func TransformVec3(m mgl64.Mat4, v mgl64.Vec3) mgl64.Vec3 {
	return m.Mul4x1(v.Vec4(1)).Vec3()
}

// Transform transforms the ray by the given matrix. The length of the ray after the
// tranformation has uses so it is returned as well.
func (r Ray) Transform(transform mgl64.Mat4) (newRay Ray, len float64) {
	newOrigin := TransformVec3(transform, r.Origin)
	newDir := TransformVec3(transform, r.Origin.Add(r.Direction)).Sub(newOrigin)
	len = newDir.Len()
	return NewRay(r.Type, newOrigin, newDir), len
}

// Reflect returns the reflected ray at the given Intersection.
func (r Ray) Reflect(isect Intersection) (reflRay Ray) {
	mDir := r.Direction.Mul(-1)
	nMD2 := isect.Normal.Dot(mDir) * 2
	reflDir := isect.Normal.Mul(nMD2).Add(r.Direction)
	return NewRay(ReflectionRay, r.At(isect.T), reflDir)
}

// Refract returns the refracted ray at the given intersection.
func (r Ray) Refract(isect Intersection, n1, n2 float64) (refrRay Ray) {
	nn := n1 / n2
	cosθ1 := isect.Normal.Dot(r.Direction.Mul(-1))
	cosθ2 := math.Sqrt(1 - nn*nn*(1-cosθ1*cosθ1))
	refractDir := r.Direction.Mul(nn).Add(isect.Normal.Mul(nn*cosθ1 - cosθ2))
	return NewRay(RefractionRay, r.At(isect.T), refractDir)
}
