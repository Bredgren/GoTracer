package raytracer

import (
	"testing"

	"github.com/go-gl/mathgl/mgl64"
)

var sphere = SphereObject{}

func checkIsect(t *testing.T, isect Intersection, expNormal mgl64.Vec3, expT float64) {
	if isect.Normal != expNormal {
		t.Errorf("Incorrect normal %v, expected %v", isect.Normal, expNormal)
	}
	if !mgl64.FloatEqualThreshold(isect.T, expT, Rayε) {
		t.Errorf("Incorrect T %v, expected %v", isect.T, expT)
	}
}

func checkIsectSanity(t *testing.T, isect Intersection) {
	if !mgl64.FloatEqualThreshold(isect.Normal.Len(), 1.0, Rayε) {
		t.Errorf("Incorrect normal length %v %v", isect.Normal.Len(), isect.Normal)
	}
	if isect.T < Rayε {
		t.Errorf("Incorrect T %v, too small", isect.T)
	}
}

func testIsect(t *testing.T, ray Ray, expNormal mgl64.Vec3, expT float64) {
	if isect, hit := sphere.Intersect(ray); hit {
		checkIsect(t, isect, expNormal, expT)
	} else {
		t.Errorf("Ray %v did not hit sphere", ray)
	}
}

func testIsectHit(t *testing.T, ray Ray) {
	if isect, hit := sphere.Intersect(ray); hit {
		checkIsectSanity(t, isect)
	} else {
		t.Errorf("Ray %v did not hit sphere", ray)
	}
}

func testIsectNoHit(t *testing.T, ray Ray) {
	if _, hit := sphere.Intersect(ray); hit {
		t.Errorf("Ray %v hit sphere", ray)
	}
}

func TestSphereIntersect(t *testing.T) {
	// Origin ray
	ray := NewRay(mgl64.Vec3{0, 0, 5}, mgl64.Vec3{0, 0, -1})
	testIsect(t, ray, mgl64.Vec3{0, 0, 1}, 4)

	// Origin ray behind
	ray.Origin = mgl64.Vec3{0, 0, -5}
	testIsectNoHit(t, ray)

	// Left graze
	ray.Origin = mgl64.Vec3{-(1 - Rayε), 0, 0}
	testIsectHit(t, ray)
	ray.Origin = mgl64.Vec3{-1, 0, 0}
	testIsectNoHit(t, ray)

	// Right graze
	ray.Origin = mgl64.Vec3{1 - Rayε, 0, 0}
	testIsectHit(t, ray)
	ray.Origin = mgl64.Vec3{1, 0, 0}
	testIsectNoHit(t, ray)

	// From inside
	ray.Origin = mgl64.Vec3{0, 0, 0}
	testIsect(t, ray, mgl64.Vec3{0, 0, -1}, 1)

	// Really close
	ray.Origin = mgl64.Vec3{0, 0, 1}
	testIsect(t, ray, mgl64.Vec3{0, 0, -1}, 2)
	ray.Origin = mgl64.Vec3{0, 0, 1 + Rayε}
	testIsect(t, ray, mgl64.Vec3{0, 0, 1}, Rayε)
	ray.Origin = mgl64.Vec3{0, 0, 1 - Rayε}
	testIsect(t, ray, mgl64.Vec3{0, 0, -1}, 2 - Rayε)
	ray.Origin = mgl64.Vec3{0, 0, -1}
	testIsectNoHit(t, ray)
}

func BenchmarkSphereIntersect(b *testing.B) {
	ray := Ray{mgl64.Vec3{0, 0, 5}, mgl64.Vec3{0, 0, -1}}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sphere.Intersect(ray)
	}
}
