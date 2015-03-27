package gotracer

import (
	"testing"

	"github.com/go-gl/mathgl/mgl64"
)

func TestSceneIntersect(t *testing.T) {
	scene := Scene{}
	sphere := Sphere{NewObject(mgl64.Ident4(), nil)}
	scene.Objects = append(scene.Objects, &sphere)

	isect := Intersection{}
	ray := Ray{PrimaryRay, mgl64.Vec3{0, 0, 5}, mgl64.Vec3{0, 0, -1}}
	expIsect := Intersection{mgl64.Vec3{0, 0, 1}, 4, nil, mgl64.Vec2{0.5, 0}}

	hit := scene.Intersect(&ray, &isect)
	if !hit {
		t.Errorf("Expected intersection")
	}
	if !isectEqual(isect, expIsect) {
		t.Errorf("Incorrect intersection. Expected %v, got %v", expIsect, isect)
	}

	sphere.Object = NewObject(mgl64.Translate3D(2, 0, 0), nil)
	hit = scene.Intersect(&ray, &isect)
	if hit {
		t.Errorf("Expected no intersection")
	}

	ray.Origin = mgl64.Vec3{2, 0, 5}

	hit = scene.Intersect(&ray, &isect)
	if !hit {
		t.Errorf("Expected intersection")
	}
	if !isectEqual(isect, expIsect) {
		t.Errorf("Incorrect intersection. Expected %v, got %v", expIsect, isect)
	}
}

// func checkSceneIsect(t *testing.T, isect Intersection, expObj SceneObject, expNormal mgl64.Vec3, expT float64) {
// 	if isect.Object != expObj {
// 		t.Errorf("Incorrect object %v, expected %v", isect.Object, expObj)
// 	}
// 	if isect.Normal != expNormal {
// 		t.Errorf("Incorrect normal %v, expected %v", isect.Normal, expNormal)
// 	}
// 	if !mgl64.FloatEqualThreshold(isect.T, expT, Rayε) {
// 		t.Errorf("Incorrect T %v, expected %v", isect.T, expT)
// 	}
// }

// func TestIntersectScaled(t *testing.T) {
// 	scene := &Scene{}

// 	transform := mgl64.Scale3D(0.5, 0.5, 0.5)
// 	sphere1 := SphereObject{Transform: transform}
// 	InitSphereObject(&sphere1)
// 	scene.Objects = append(scene.Objects, sphere1)

// 	ray := NewRay(PrimaryRay, mgl64.Vec3{0, 0, 5}, mgl64.Vec3{0, 0, -1})
// 	if _, found := scene.Intersect(ray); !found {
// 		t.Errorf("Ray %v didn't intersect", ray)
// 	}

// 	ray.Origin = mgl64.Vec3{0, 0, -5}
// 	if _, found := scene.Intersect(ray); found {
// 		t.Errorf("Ray %v intersected", ray)
// 	}

// 	ray.Origin = mgl64.Vec3{-(0.5 - Rayε), 0, 5}
// 	if _, found := scene.Intersect(ray); !found {
// 		t.Errorf("Ray %v didn't intersect", ray)
// 	}
// 	ray.Origin = mgl64.Vec3{-(0.5 + Rayε), 0, 5}
// 	if _, found := scene.Intersect(ray); found {
// 		t.Errorf("Ray %v intersected", ray)
// 	}
// }

// func TestIntersectTranslated(t *testing.T) {
// 	scene := &Scene{}

// 	transform := mgl64.Translate3D(1, 0, 0)
// 	sphere1 := SphereObject{Transform: transform}
// 	InitSphereObject(&sphere1)
// 	scene.Objects = append(scene.Objects, sphere1)

// 	ray := NewRay(PrimaryRay, mgl64.Vec3{1, 0, 5}, mgl64.Vec3{0, 0, -1})
// 	if _, found := scene.Intersect(ray); !found {
// 		t.Errorf("Ray %v didn't intersect", ray)
// 	}

// 	ray.Origin = mgl64.Vec3{1, 0, -5}
// 	if _, found := scene.Intersect(ray); found {
// 		t.Errorf("Ray %v intersected", ray)
// 	}

// 	ray.Origin = mgl64.Vec3{Rayε, 0, 5}
// 	if _, found := scene.Intersect(ray); !found {
// 		t.Errorf("Ray %v didn't intersect", ray)
// 	}
// 	ray.Origin = mgl64.Vec3{-Rayε, 0, 5}
// 	if _, found := scene.Intersect(ray); found {
// 		t.Errorf("Ray %v intersected", ray)
// 	}

// 	ray.Origin = mgl64.Vec3{2 - Rayε, 0, 5}
// 	if _, found := scene.Intersect(ray); !found {
// 		t.Errorf("Ray %v didn't intersect", ray)
// 	}
// 	ray.Origin = mgl64.Vec3{2 + Rayε, 0, 5}
// 	if _, found := scene.Intersect(ray); found {
// 		t.Errorf("Ray %v intersected", ray)
// 	}
// }

// func TestIntersectScaledTranslated(t *testing.T) {
// 	scene := &Scene{}

// 	transform := mgl64.Translate3D(1, 0, 0.5).Mul4(mgl64.Scale3D(0.5, 0.5, 0.5))
// 	sphere1 := SphereObject{Transform: transform}
// 	InitSphereObject(&sphere1)
// 	scene.Objects = append(scene.Objects, sphere1)

// 	ray := NewRay(PrimaryRay, mgl64.Vec3{1, 0, 5}, mgl64.Vec3{0, 0, -1})
// 	if _, found := scene.Intersect(ray); !found {
// 		t.Errorf("Ray %v didn't intersect", ray)
// 	}

// 	ray.Origin = mgl64.Vec3{1, 0, -5}
// 	if _, found := scene.Intersect(ray); found {
// 		t.Errorf("Ray %v intersected", ray)
// 	}

// 	ray.Origin = mgl64.Vec3{0.5 + Rayε, 0, 5}
// 	if _, found := scene.Intersect(ray); !found {
// 		t.Errorf("Ray %v didn't intersect", ray)
// 	}
// 	ray.Origin = mgl64.Vec3{0.5 - Rayε, 0, 5}
// 	if _, found := scene.Intersect(ray); found {
// 		t.Errorf("Ray %v intersected", ray)
// 	}

// 	ray.Origin = mgl64.Vec3{1.5 - Rayε, 0, 5}
// 	if _, found := scene.Intersect(ray); !found {
// 		t.Errorf("Ray %v didn't intersect", ray)
// 	}
// 	ray.Origin = mgl64.Vec3{1.5 + Rayε, 0, 5}
// 	if _, found := scene.Intersect(ray); found {
// 		t.Errorf("Ray %v intersected", ray)
// 	}
// }

// func TestIntersect(t *testing.T) {
// 	t.Skip()
// 	scene := &Scene{}

// 	transform := mgl64.Scale3D(0.5, 0.5, 0.5)
// 	sphere1 := SphereObject{Transform: transform, MaterialName: "1"}
// 	InitSphereObject(&sphere1)
// 	// scene.Objects = append(scene.Objects, sphere1)

// 	transform = mgl64.Translate3D(1.8, 0, -0.5).Mul4(mgl64.Scale3D(0.5, 0.5, 0.5))
// 	sphere2 := SphereObject{Transform: transform, MaterialName: "2"}
// 	InitSphereObject(&sphere2)
// 	_ = sphere2
// 	// scene.Objects = append(scene.Objects, sphere2)

// 	transform = mgl64.Translate3D(-0.8, 0, 0.5).Mul4(mgl64.Scale3D(0.5, 0.5, 0.5))
// 	sphere3 := SphereObject{Transform: transform, MaterialName: "3"}
// 	InitSphereObject(&sphere3)
// 	scene.Objects = append(scene.Objects, sphere3)

// 	ray := NewRay(PrimaryRay, mgl64.Vec3{-0.2 - Rayε, 0, 5}, mgl64.Vec3{0, 0, -1})
// 	if isect, found := scene.Intersect(ray); found {
// 		if isect.Object != sphere3 {
// 			t.Errorf("Incorrect object %v, expected %v", isect.Object.GetMaterialName(), sphere3.GetMaterialName())
// 		}
// 	} else {
// 		t.Errorf("Ray %v didn't intersect", ray)
// 	}

// 	ray = NewRay(PrimaryRay, mgl64.Vec3{-0.2, 0, 5}, mgl64.Vec3{0, 0, -1})
// 	if isect, found := scene.Intersect(ray); found {
// 		if isect.Object != sphere1 {
// 			t.Errorf("Incorrect object %v, expected %v", isect.Object.GetMaterialName(), sphere1.GetMaterialName())
// 		}
// 	} else {
// 		t.Errorf("Ray %v didn't intersect", ray)
// 	}

// 	ray = NewRay(PrimaryRay, mgl64.Vec3{-0.2 + Rayε, 0, 5}, mgl64.Vec3{0, 0, -1})
// 	if isect, found := scene.Intersect(ray); found {
// 		if isect.Object != sphere1 {
// 			t.Errorf("Incorrect object %v, expected %v", isect.Object.GetMaterialName(), sphere1.GetMaterialName())
// 		}
// 	} else {
// 		t.Errorf("Ray %v didn't intersect", ray)
// 	}

// 	// ray := NewRay(PrimaryRay, mgl64.Vec3{0, 0, 5}, mgl64.Vec3{0, 0, -1})
// 	// if isect, found := scene.Intersect(ray); found {
// 	// 	expNormal := mgl64.Vec3{0, 0, 1}
// 	// 	expT := 4.5
// 	// 	checkSceneIsect(t, isect, sphere1, expNormal, expT)
// 	// } else {
// 	// 	t.Errorf("Ray %v didn't intersect", ray)
// 	// }

// 	// if isect, found := scene.Intersect(ray); found {
// 	// 	expNormal := mgl64.Vec3{0, 0, 1}
// 	// 	expT := 4.5
// 	// 	checkSceneIsect(t, isect, sphere1, expNormal, expT)
// 	// } else {
// 	// 	t.Errorf("Ray %v didn't intersect", ray)
// 	// }

// 	// ray = NewRay(PrimaryRay, mgl64.Vec3{0, 0, -10}, mgl64.Vec3{0, 0, 1})
// 	// if isect, found := scene.Intersect(ray); found {
// 	// 	expNormal := mgl64.Vec3{0, 0, -1}
// 	// 	expT := 4.0
// 	// 	checkSceneIsect(t, isect, sphere2, expNormal, expT)
// 	// } else {
// 	// 	t.Errorf("Ray %v didn't intersect", ray)
// 	// }

// 	// transform = mgl64.Translate3D(0, 0, 5)
// 	// sphere3 := SphereObject{transform, transform.Inv(), ""}
// 	// scene.Objects = append(scene.Objects, sphere3)
// 	// ray = NewRay(PrimaryRay, mgl64.Vec3{0, 0, 10}, mgl64.Vec3{0, 0, -1})
// 	// if isect, found := scene.Intersect(ray); found {
// 	// 	expNormal := mgl64.Vec3{0, 0, 1}
// 	// 	expT := 4.0
// 	// 	checkSceneIsect(t, isect, sphere3, expNormal, expT)
// 	// } else {
// 	// 	t.Errorf("Ray %v didn't intersect", ray)
// 	// }
// }
