package gotracer

import (
	"math/rand"
	"testing"

	"github.com/go-gl/mathgl/mgl64"
)

func TestSphereIntersect(t *testing.T) {
	obj := NewObject(mgl64.Ident4(), nil)
	sphere := Sphere{obj}

	ray := NewRay(PrimaryRay, mgl64.Vec3{0, 0, 5}, mgl64.Vec3{0, 0, -1})
	exp := ExpectedIsect{Intersection{obj, mgl64.Vec3{0, 0, 1}, 4, mgl64.Vec2{0.5, 0}}, true}
	testIsect(t, "Origin", sphere, ray, exp)

	ray.Origin = mgl64.Vec3{0, 0, -5}
	exp.hit = false
	testIsect(t, "Behind", sphere, ray, exp)

	ray.Origin = mgl64.Vec3{-(1 - Rayε), 0, 5}
	exp.isect.Object = nil
	exp.hit = true
	testIsect(t, "Left graze hit", sphere, ray, exp)
	ray.Origin = mgl64.Vec3{-1, 0, 0}
	exp.hit = false
	testIsect(t, "Left graze no hit", sphere, ray, exp)

	ray.Origin = mgl64.Vec3{1 - Rayε, 0, 5}
	exp.hit = true
	testIsect(t, "Right graze hit", sphere, ray, exp)
	ray.Origin = mgl64.Vec3{1, 0, 0}
	exp.hit = false
	testIsect(t, "Right graze no hit", sphere, ray, exp)

	ray.Origin = mgl64.Vec3{0, 0, 0}
	exp.isect.Object = obj
	exp.isect.Normal = mgl64.Vec3{0, 0, -1}
	exp.isect.T = 1
	exp.isect.UVCoords = mgl64.Vec2{0.5, 1}
	exp.hit = true
	testIsect(t, "Inside", sphere, ray, exp)

	ray.Origin = mgl64.Vec3{0, 0, 1}
	exp.isect.Object = obj
	exp.isect.Normal = mgl64.Vec3{0, 0, -1}
	exp.isect.T = 2
	exp.isect.UVCoords = mgl64.Vec2{0.5, 1}
	testIsect(t, "Touching front", sphere, ray, exp)
	ray.Origin = mgl64.Vec3{0, 0, 1 + Rayε}
	exp.isect.Normal = mgl64.Vec3{0, 0, 1}
	exp.isect.T = Rayε
	exp.isect.UVCoords = mgl64.Vec2{0.5, 0}
	testIsect(t, "Close to front", sphere, ray, exp)
	ray.Origin = mgl64.Vec3{0, 0, 1 - Rayε}
	exp.isect.Normal = mgl64.Vec3{0, 0, -1}
	exp.isect.T = 2 - Rayε
	exp.isect.UVCoords = mgl64.Vec2{0.5, 1}
	testIsect(t, "Close inside front", sphere, ray, exp)
	ray.Origin = mgl64.Vec3{0, 0, -1}
	exp.hit = false
	testIsect(t, "Touching back", sphere, ray, exp)
}

func BenchmarkSphereIntersect(b *testing.B) {
	sphere := Sphere{NewObject(mgl64.Ident4(), nil)}

	ray := NewRay(PrimaryRay, mgl64.Vec3{0, 0, 5}, mgl64.Vec3{0, 0, -1})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sphere.Intersect(&ray)
	}
}

func BenchmarkSphereIntersectRandom(b *testing.B) {
	sphere := Sphere{NewObject(mgl64.Ident4(), nil)}

	min := -1.0
	max := 1.0
	rays := make([]*Ray, 0)
	for i := 0; i < b.N; i++ {
		x := rand.Float64()*(max-min) - max
		y := rand.Float64()*(max-min) - max
		z := rand.Float64()*(max-min) - max
		point := mgl64.Vec3{x, y, z}
		dir := point.Mul(-1).Normalize()
		ray := NewRay(PrimaryRay, point, dir)
		rays = append(rays, &ray)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sphere.Intersect(rays[i])
	}
}
