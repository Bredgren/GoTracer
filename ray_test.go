package gotracer

import (
	"math"
	"testing"

	"github.com/go-gl/mathgl/mgl64"
)

func checkAt(test *testing.T, ray Ray, t float64, exp mgl64.Vec3) {
	at := ray.At(t)
	if !at.ApproxEqualThreshold(exp, Rayε) {
		test.Errorf("Ray.At(%v) = %v, expected %v", t, at, exp)
	}
}

func TestRayAt(t *testing.T) {
	origin := mgl64.Vec3{0, 0, 1}
	direction := mgl64.Vec3{0, 0, 1}
	ray := Ray{PrimaryRay, origin, direction}
	checkAt(t, ray, 0, origin)
	checkAt(t, ray, 1, mgl64.Vec3{0, 0, 2})
}

func checkRay(t *testing.T, ray Ray, expOrigin, expDirection mgl64.Vec3) {
	if !ray.Origin.ApproxEqualThreshold(expOrigin, Rayε) {
		t.Errorf("Incorrect origin %v, expected %v", ray.Origin, expOrigin)
	}
	if !ray.Dir.ApproxEqualThreshold(expDirection, Rayε) {
		t.Errorf("Incorrect direction %v, expected %v", ray.Dir, expDirection)
	}
	if !mgl64.FloatEqualThreshold(ray.Dir.Len(), 1.0, Rayε) {
		t.Errorf("Direction not normal %v", ray.Dir)
	}
}

func TestRayTransform(t *testing.T) {
	ray := Ray{PrimaryRay, mgl64.Vec3{0, 0, 0}, mgl64.Vec3{0, 0, 1}}

	translate := mgl64.Translate3D(1, 2, 3)
	newRay, _ := ray.Transform(&translate)
	expOrigin := mgl64.Vec3{1, 2, 3}
	expDirection := mgl64.Vec3{0, 0, 1}
	checkRay(t, newRay, expOrigin, expDirection)

	rad := mgl64.DegToRad(45)
	rotate := mgl64.HomogRotate3D(rad, mgl64.Vec3{0, 1, 0})
	newRay, _ = ray.Transform(&rotate)
	expOrigin = mgl64.Vec3{0, 0, 0}
	expDirection = mgl64.Vec3{math.Sin(rad), 0, math.Cos(rad)}
	checkRay(t, newRay, expOrigin, expDirection)

	scale := mgl64.Scale3D(2, 2, 2)
	newRay, _ = ray.Transform(&scale)
	expOrigin = mgl64.Vec3{0, 0, 0}
	expDirection = mgl64.Vec3{0, 0, 1}
	checkRay(t, newRay, expOrigin, expDirection)

	ray = Ray{PrimaryRay, mgl64.Vec3{1, 0, 0}, mgl64.Vec3{1, 0, 1}.Normalize()}

	translate = mgl64.Translate3D(1, 2, 3)
	newRay, _ = ray.Transform(&translate)
	expOrigin = mgl64.Vec3{2, 2, 3}
	expDirection = mgl64.Vec3{1, 0, 1}.Normalize()
	checkRay(t, newRay, expOrigin, expDirection)

	rad = mgl64.DegToRad(45)
	rotate = mgl64.HomogRotate3D(rad, mgl64.Vec3{0, 1, 0})
	newRay, _ = ray.Transform(&rotate)
	expOrigin = mgl64.Vec3{math.Sin(rad), 0, -math.Cos(rad)}
	expDirection = mgl64.Vec3{1, 0, 0}
	checkRay(t, newRay, expOrigin, expDirection)

	scale = mgl64.Scale3D(2, 2, 2)
	newRay, _ = ray.Transform(&scale)
	expOrigin = mgl64.Vec3{2, 0, 0}
	expDirection = mgl64.Vec3{1, 0, 1}.Normalize()
	checkRay(t, newRay, expOrigin, expDirection)

	x := -(0.5 - Rayε)
	ray = Ray{PrimaryRay, mgl64.Vec3{x, 0, -5}, mgl64.Vec3{0, 0, -1}}
	scale = mgl64.Scale3D(0.5, 0.5, 0.5).Inv()
	newRay, _ = ray.Transform(&scale)
	expOrigin = mgl64.Vec3{2 * x, 0, -10}
	expDirection = mgl64.Vec3{0, 0, -1}
	checkRay(t, newRay, expOrigin, expDirection)
}

func checkReflectAngle(t *testing.T, angle float64) {
	rad := mgl64.DegToRad(angle)
	x := math.Sin(rad)
	y := math.Cos(rad)
	ray := Ray{PrimaryRay, mgl64.Vec3{x, y, 0}, mgl64.Vec3{-x, -y, 0}}
	isect := Intersection{mgl64.Vec3{0, 1, 0}, 1, nil, mgl64.Vec2{}}
	reflRay := ray.Reflect(&isect)
	expOrigin := mgl64.Vec3{0, 0, 0}
	expDirection := mgl64.Vec3{-x, y, 0}.Normalize()
	checkRay(t, reflRay, expOrigin, expDirection)
}

func TestRayReflect(t *testing.T) {
	for angle := 0.0; angle <= 360; angle++ {
		checkReflectAngle(t, angle)
	}
}

func TestRayRefract(t *testing.T) {
	ray := Ray{PrimaryRay, mgl64.Vec3{0, 0, 0}, mgl64.Vec3{0.707107, -0.707107, 0}}
	isect := Intersection{mgl64.Vec3{0, 1, 0}, 1, nil, mgl64.Vec2{}}
	refrRay := ray.Refract(&isect, 0.9, 1)
	expOrigin := mgl64.Vec3{0.707107, -0.707107, 0}
	expDirection := mgl64.Vec3{0.636396, -0.771352, 0}
	checkRay(t, refrRay, expOrigin, expDirection)
}

func BenchmarkRayTransform(b *testing.B) {
	transform := mgl64.Ident4()
	transform = transform.Mul4(mgl64.Translate3D(1, 2, 3))
	rad := mgl64.DegToRad(45)
	transform = transform.Mul4(mgl64.HomogRotate3D(rad, mgl64.Vec3{0.5, 1, 1.5}))
	transform = transform.Mul4(mgl64.Scale3D(1.2, 2.3, 3.4))
	ray := Ray{PrimaryRay, mgl64.Vec3{1, 0, 0}, mgl64.Vec3{1, 0, 1}.Normalize()}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ray.Transform(&transform)
	}
}

func BenchmarkRayReflect(b *testing.B) {
	ray := Ray{PrimaryRay, mgl64.Vec3{1, 1, 0}, mgl64.Vec3{-1, -1, 0}}
	isect := Intersection{mgl64.Vec3{0, 1, 0}, 1, nil, mgl64.Vec2{}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ray.Reflect(&isect)
	}
}

func BenchmarkRayRefract(b *testing.B) {
	ray := Ray{PrimaryRay, mgl64.Vec3{0, 0, 0}, mgl64.Vec3{0.707107, -0.707107, 0}}
	isect := Intersection{mgl64.Vec3{0, 1, 0}, 1, nil, mgl64.Vec2{}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ray.Refract(&isect, 0.9, 1)
	}
}
