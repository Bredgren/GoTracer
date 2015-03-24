package gotracer

import (
	"math/rand"
	"testing"

	"github.com/go-gl/mathgl/mgl64"
)

type ExpectedIsect struct {
	isect Intersection
	hit   bool
}

func isectEqual(i1, i2 Intersection) bool {
	return i1.Object == i2.Object &&
		i1.Normal.ApproxEqualThreshold(i2.Normal, Rayε) &&
		mgl64.FloatEqualThreshold(i1.T, i2.T, Rayε) &&
		i1.UVCoords.ApproxEqualThreshold(i2.UVCoords, Rayε)
}

// exp.isect.Object = nil means ignore isect.
func testIsect(t *testing.T, desc string, itr Intersecter, ray Ray, exp ExpectedIsect) {
	isect, hit := itr.Intersect(ray)
	if !hit && hit == exp.hit {
		return // isect is undefined if we don't hit
	}
	if (exp.isect.Object != nil && !isectEqual(isect, exp.isect)) || hit != exp.hit {
		actual := ExpectedIsect{isect, hit}
		t.Errorf("%v: Incorrect intersection. Expected %v, got %v", desc, exp, actual)
	}
}

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
		sphere.Intersect(ray)
	}
}

func BenchmarkSphereIntersectRandom(b *testing.B) {
	sphere := Sphere{NewObject(mgl64.Ident4(), nil)}

	min := -1.0
	max := 1.0
	rays := make([]Ray, 0)
	for i := 0; i < b.N; i++ {
		x := rand.Float64()*(max-min) - max
		y := rand.Float64()*(max-min) - max
		ray := NewRay(PrimaryRay, mgl64.Vec3{x, y, 5}, mgl64.Vec3{0, 0, -1})
		rays = append(rays, ray)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sphere.Intersect(rays[i])
	}
}

// var box = BoxObject{}

// func TestBoxIntersect(t *testing.T) {
// 	InitBoxObject(&box)

// 	// Origin ray
// 	ray := NewRay(PrimaryRay, mgl64.Vec3{0, 0, 5}, mgl64.Vec3{0, 0, -1})
// 	testIsect(t, "Origin", box, ray, mgl64.Vec3{0, 0, 1}, 4.5)

// 	// Origin ray behind
// 	ray.Origin = mgl64.Vec3{0, 0, -5}
// 	testIsectNoHit(t, "Behind", box, ray)

// 	// Left graze
// 	ray.Origin = mgl64.Vec3{-0.5 + Rayε, 0, 5}
// 	testIsect(t, "Left graze hit", box, ray, mgl64.Vec3{0, 0, 1}, 4.5)
// 	ray.Origin = mgl64.Vec3{-0.5 - Rayε, 0, 5}
// 	testIsectNoHit(t, "Left graze no hit", box, ray)

// 	// Right graze
// 	ray.Origin = mgl64.Vec3{0.5 - Rayε, 0, 5}
// 	testIsect(t, "Right graze hit", box, ray, mgl64.Vec3{0, 0, 1}, 4.5)
// 	ray.Origin = mgl64.Vec3{0.5 + Rayε, 0, 0}
// 	testIsectNoHit(t, "Right graze no hit", box, ray)

// 	// From inside
// 	ray.Origin = mgl64.Vec3{0, 0, 0}
// 	testIsect(t, "Inside", box, ray, mgl64.Vec3{0, 0, -1}, 0.5)

// 	// Really close
// 	ray.Origin = mgl64.Vec3{0, 0, 0.5}
// 	testIsect(t, "Touching front", box, ray, mgl64.Vec3{0, 0, -1}, 1)
// 	ray.Origin = mgl64.Vec3{0, 0, 0.5 + 2 * Rayε}
// 	testIsect(t, "Close to front", box, ray, mgl64.Vec3{0, 0, 1}, 2 * Rayε)
// 	ray.Origin = mgl64.Vec3{0, 0, 0.5 - Rayε}
// 	testIsect(t, "Insde front", box, ray, mgl64.Vec3{0, 0, -1}, 1 - Rayε)
// 	ray.Origin = mgl64.Vec3{0, 0, -0.5}
// 	testIsectNoHit(t, "Touching back", box, ray)
// }

// func BenchmarkBoxIntersect(b *testing.B) {
// 	InitBoxObject(&box)

// 	ray := NewRay(PrimaryRay, mgl64.Vec3{0, 0, 5}, mgl64.Vec3{0, 0, -1})
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		box.Intersect(ray)
// 	}
// }

// var square = SquareObject{}

// func TestSquareIntersect(t *testing.T) {
// 	InitSquareObject(&square)

// 	// Origin ray
// 	ray := NewRay(PrimaryRay, mgl64.Vec3{0, 0, 5}, mgl64.Vec3{0, 0, -1})
// 	testIsect(t, "Origin", square, ray, mgl64.Vec3{0, 0, 1}, 5)

// 	// Origin ray behind
// 	ray.Origin = mgl64.Vec3{0, 0, -5}
// 	testIsectNoHit(t, "Behind", square, ray)

// 	// Left graze
// 	ray.Origin = mgl64.Vec3{-0.5 + Rayε, 0, 5}
// 	testIsect(t, "Left graze hit", square, ray, mgl64.Vec3{0, 0, 1}, 5)
// 	ray.Origin = mgl64.Vec3{-0.5 - Rayε, 0, 5}
// 	testIsectNoHit(t, "Left graze no hit", square, ray)

// 	// Right graze
// 	ray.Origin = mgl64.Vec3{0.5 - Rayε, 0, 5}
// 	testIsect(t, "Right graze hit", square, ray, mgl64.Vec3{0, 0, 1}, 5)
// 	ray.Origin = mgl64.Vec3{0.5 + Rayε, 0, 5}
// 	testIsectNoHit(t, "Right graze no hit", square, ray)

// 	// Really close
// 	ray.Origin = mgl64.Vec3{0, 0, 0}
// 	testIsectNoHit(t, "Touching", square, ray)
// 	ray.Origin = mgl64.Vec3{0, 0, 2 * Rayε}
// 	testIsect(t, "Close to front", square, ray, mgl64.Vec3{0, 0, 1}, 2 * Rayε)
// 	ray.Origin = mgl64.Vec3{0, 0, -Rayε}
// 	testIsectNoHit(t, "Insde front", square, ray)
// 	ray.Origin = mgl64.Vec3{0, 0, -2 * Rayε}
// 	ray.Direction = mgl64.Vec3{0, 0, 1}
// 	testIsect(t, "Close to back", square, ray, mgl64.Vec3{0, 0, -1}, 2 * Rayε)
// }

// func BenchmarkSquareIntersect(b *testing.B) {
// 	InitSquareObject(&square)

// 	ray := NewRay(PrimaryRay, mgl64.Vec3{0, 0, 5}, mgl64.Vec3{0, 0, -1})
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		square.Intersect(ray)
// 	}
// }

// var cylinder = CylinderObject{}

// func TestCylinderIntersect(t *testing.T) {
// 	InitCylinderObject(&cylinder)

// 	// Origin ray
// 	ray := NewRay(PrimaryRay, mgl64.Vec3{0, 0, 5}, mgl64.Vec3{0, 0, -1})
// 	testIsectNoHit(t, "Origin", cylinder, ray)

// 	cylinder.Capped = true
// 	testIsect(t, "Origin", cylinder, ray, mgl64.Vec3{0, 0, 1}, 4)
// }

// func BenchmarkCylinderIntersect(b *testing.B) {
// 	cylinder.Capped = true
// 	InitCylinderObject(&cylinder)

// 	ray := NewRay(PrimaryRay, mgl64.Vec3{0, 0, 5}, mgl64.Vec3{0, 0, -1})
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		cylinder.Intersect(ray)
// 	}
// }

// var cone = ConeObject{}

// func BenchmarkConeIntersect(b *testing.B) {
// 	cone.Capped = true
// 	cone.BaseRadius = 0.1
// 	cone.TopRadius = 1.0
// 	InitConeObject(&cone)

// 	ray := NewRay(PrimaryRay, mgl64.Vec3{0, 0, 5}, mgl64.Vec3{0, 0, -1})
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		cone.Intersect(ray)
// 	}
// }
