package raytracer

import (
	"testing"

	"github.com/go-gl/mathgl/mgl64"
)

var sphere = SphereObject{}

func checkIsect(t *testing.T, test string, isect Intersection, expNormal mgl64.Vec3, expT float64) {
	if isect.Normal != expNormal {
		t.Errorf("%v: Incorrect normal %v, expected %v", test, isect.Normal, expNormal)
	}
	if !mgl64.FloatEqualThreshold(isect.T, expT, Rayε) {
		t.Errorf("%v: Incorrect T %v, expected %v", test, isect.T, expT)
	}
}

func checkIsectSanity(t *testing.T, test string, isect Intersection) {
	if !mgl64.FloatEqualThreshold(isect.Normal.Len(), 1.0, Rayε) {
		t.Errorf("%v: Incorrect normal length %v %v", test, isect.Normal.Len(), isect.Normal)
	}
	if isect.T < Rayε {
		t.Errorf("%v: Incorrect T %v, too small", test, isect.T)
	}
}

func testIsect(t *testing.T, test string, object SceneObject, ray Ray, expNormal mgl64.Vec3, expT float64) {
	if isect, hit := object.Intersect(ray); hit {
		checkIsect(t, test, isect, expNormal, expT)
	} else {
		t.Errorf("%v: Ray %v did not hit object", test, ray)
	}
}

func testIsectHit(t *testing.T, test string, object SceneObject, ray Ray) {
	if isect, hit := object.Intersect(ray); hit {
		checkIsectSanity(t, test, isect)
	} else {
		t.Errorf("%v: Ray %v did not hit object", test, ray)
	}
}

func testIsectNoHit(t *testing.T, test string, object SceneObject, ray Ray) {
	if _, hit := object.Intersect(ray); hit {
		t.Errorf("%v: Ray %v hit object", test, ray)
	}
}

func TestSphereIntersect(t *testing.T) {
	InitSphereObject(&sphere)

	// Origin ray
	ray := NewRay(PrimaryRay, mgl64.Vec3{0, 0, 5}, mgl64.Vec3{0, 0, -1})
	testIsect(t, "Origin", sphere, ray, mgl64.Vec3{0, 0, 1}, 4)

	// Origin ray behind
	ray.Origin = mgl64.Vec3{0, 0, -5}
	testIsectNoHit(t, "Behind", sphere, ray)

	// Left graze
	ray.Origin = mgl64.Vec3{-(1 - Rayε), 0, 5}
	testIsectHit(t, "Left graze hit", sphere, ray)
	ray.Origin = mgl64.Vec3{-1, 0, 0}
	testIsectNoHit(t, "Left graze no hit", sphere, ray)

	// Right graze
	ray.Origin = mgl64.Vec3{1 - Rayε, 0, 5}
	testIsectHit(t, "Right graze hit", sphere, ray)
	ray.Origin = mgl64.Vec3{1, 0, 0}
	testIsectNoHit(t, "Right graze hit", sphere, ray)

	// From inside
	ray.Origin = mgl64.Vec3{0, 0, 0}
	testIsect(t, "Inside", sphere, ray, mgl64.Vec3{0, 0, -1}, 1)

	// Really close
	ray.Origin = mgl64.Vec3{0, 0, 1}
	testIsect(t, "Touching front", sphere, ray, mgl64.Vec3{0, 0, -1}, 2)
	ray.Origin = mgl64.Vec3{0, 0, 1 + Rayε}
	testIsect(t, "Close to front", sphere, ray, mgl64.Vec3{0, 0, 1}, Rayε)
	ray.Origin = mgl64.Vec3{0, 0, 1 - Rayε}
	testIsect(t, "Close inside front", sphere, ray, mgl64.Vec3{0, 0, -1}, 2 - Rayε)
	ray.Origin = mgl64.Vec3{0, 0, -1}
	testIsectNoHit(t, "Touching back", sphere, ray)
}

func BenchmarkSphereIntersect(b *testing.B) {
	InitSphereObject(&sphere)

	ray := NewRay(PrimaryRay, mgl64.Vec3{0, 0, 5}, mgl64.Vec3{0, 0, -1})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sphere.Intersect(ray)
	}
}

var box = BoxObject{}

func TestBoxIntersect(t *testing.T) {
	InitBoxObject(&box)

	// Origin ray
	ray := NewRay(PrimaryRay, mgl64.Vec3{0, 0, 5}, mgl64.Vec3{0, 0, -1})
	testIsect(t, "Origin", box, ray, mgl64.Vec3{0, 0, 1}, 4.5)

	// Origin ray behind
	ray.Origin = mgl64.Vec3{0, 0, -5}
	testIsectNoHit(t, "Behind", box, ray)

	// Left graze
	ray.Origin = mgl64.Vec3{-0.5 + Rayε, 0, 5}
	testIsect(t, "Left graze hit", box, ray, mgl64.Vec3{0, 0, 1}, 4.5)
	ray.Origin = mgl64.Vec3{-0.5 - Rayε, 0, 5}
	testIsectNoHit(t, "Left graze no hit", box, ray)

	// Right graze
	ray.Origin = mgl64.Vec3{0.5 - Rayε, 0, 5}
	testIsect(t, "Right graze hit", box, ray, mgl64.Vec3{0, 0, 1}, 4.5)
	ray.Origin = mgl64.Vec3{0.5 + Rayε, 0, 0}
	testIsectNoHit(t, "Right graze no hit", box, ray)

	// From inside
	ray.Origin = mgl64.Vec3{0, 0, 0}
	testIsect(t, "Inside", box, ray, mgl64.Vec3{0, 0, -1}, 0.5)

	// Really close
	ray.Origin = mgl64.Vec3{0, 0, 0.5}
	testIsect(t, "Touching front", box, ray, mgl64.Vec3{0, 0, -1}, 1)
	ray.Origin = mgl64.Vec3{0, 0, 0.5 + 2 * Rayε}
	testIsect(t, "Close to front", box, ray, mgl64.Vec3{0, 0, 1}, 2 * Rayε)
	ray.Origin = mgl64.Vec3{0, 0, 0.5 - Rayε}
	testIsect(t, "Insde front", box, ray, mgl64.Vec3{0, 0, -1}, 1 - Rayε)
	ray.Origin = mgl64.Vec3{0, 0, -0.5}
	testIsectNoHit(t, "Touching back", box, ray)
}

func BenchmarkBoxIntersect(b *testing.B) {
	InitBoxObject(&box)

	ray := NewRay(PrimaryRay, mgl64.Vec3{0, 0, 5}, mgl64.Vec3{0, 0, -1})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		box.Intersect(ray)
	}
}

var square = SquareObject{}

func TestSquareIntersect(t *testing.T) {
	InitSquareObject(&square)

	// Origin ray
	ray := NewRay(PrimaryRay, mgl64.Vec3{0, 0, 5}, mgl64.Vec3{0, 0, -1})
	testIsect(t, "Origin", square, ray, mgl64.Vec3{0, 0, 1}, 5)

	// Origin ray behind
	ray.Origin = mgl64.Vec3{0, 0, -5}
	testIsectNoHit(t, "Behind", square, ray)

	// Left graze
	ray.Origin = mgl64.Vec3{-0.5 + Rayε, 0, 5}
	testIsect(t, "Left graze hit", square, ray, mgl64.Vec3{0, 0, 1}, 5)
	ray.Origin = mgl64.Vec3{-0.5 - Rayε, 0, 5}
	testIsectNoHit(t, "Left graze no hit", square, ray)

	// Right graze
	ray.Origin = mgl64.Vec3{0.5 - Rayε, 0, 5}
	testIsect(t, "Right graze hit", square, ray, mgl64.Vec3{0, 0, 1}, 5)
	ray.Origin = mgl64.Vec3{0.5 + Rayε, 0, 5}
	testIsectNoHit(t, "Right graze no hit", square, ray)

	// Really close
	ray.Origin = mgl64.Vec3{0, 0, 0}
	testIsectNoHit(t, "Touching", square, ray)
	ray.Origin = mgl64.Vec3{0, 0, 2 * Rayε}
	testIsect(t, "Close to front", square, ray, mgl64.Vec3{0, 0, 1}, 2 * Rayε)
	ray.Origin = mgl64.Vec3{0, 0, -Rayε}
	testIsectNoHit(t, "Insde front", square, ray)
	ray.Origin = mgl64.Vec3{0, 0, -2 * Rayε}
	ray.Direction = mgl64.Vec3{0, 0, 1}
	testIsect(t, "Close to back", square, ray, mgl64.Vec3{0, 0, -1}, 2 * Rayε)
}

func BenchmarkSquareIntersect(b *testing.B) {
	InitSquareObject(&square)

	ray := NewRay(PrimaryRay, mgl64.Vec3{0, 0, 5}, mgl64.Vec3{0, 0, -1})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		square.Intersect(ray)
	}
}


var cylinder = CylinderObject{}

func TestCylinderIntersect(t *testing.T) {
	InitCylinderObject(&cylinder)

	// Origin ray
	ray := NewRay(PrimaryRay, mgl64.Vec3{0, 0, 5}, mgl64.Vec3{0, 0, -1})
	testIsect(t, "Origin", cylinder, ray, mgl64.Vec3{0, 0, 1}, 5)

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
}

// func BenchmarkSquareIntersect(b *testing.B) {
// 	InitSquareObject(&square)

// 	ray := NewRay(PrimaryRay, mgl64.Vec3{0, 0, 5}, mgl64.Vec3{0, 0, -1})
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		square.Intersect(ray)
// 	}
// }
