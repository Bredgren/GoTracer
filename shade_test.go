package gotracer

import (
	"math/rand"
	"testing"

	"github.com/go-gl/mathgl/mgl64"
)

func TestLambertianBRDF(t *testing.T) {
	lights := []Light{&DirectionalLight{Color64{1, 1, 1}, mgl64.Vec3{0, 0, 1}}}
	ray := NewRay(PrimaryRay, mgl64.Vec3{0, 0, 5}, mgl64.Vec3{0, 0, -1})
	expColor := Color64{1, 0, 0}
	material := Material{Diffuse: UniformColor{expColor}}
	obj := NewObject(mgl64.Ident4(), &material)
	sphere := Sphere{obj}
	isect, hit := sphere.Intersect(&ray)
	if !hit {
		t.Fatal("didn't hit")
	}
	color := LambertianBRDF(lights, &ray, &isect)
	if color != expColor {
		t.Errorf("Incorrect color. Expected %v, got %v", expColor, color)
	}
}

func BenchmarkLambertianBRDF(b *testing.B) {
	lights := []Light{&DirectionalLight{Color64{1, 1, 1}, mgl64.Vec3{0, 0, 1}}}
	expColor := Color64{1, 0, 0}
	material := Material{Diffuse: UniformColor{expColor}}
	obj := NewObject(mgl64.Ident4(), &material)
	sphere := Sphere{obj}

	ray := NewRay(PrimaryRay, mgl64.Vec3{0, 0, 5}, mgl64.Vec3{0, 0, -1})
	isect, _ := sphere.Intersect(&ray)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LambertianBRDF(lights, &ray, &isect)
	}
}

type bdrfStruct struct {
	ray   *Ray
	isect *Intersection
}

func BenchmarkLambertianBRDFRandom(b *testing.B) {
	lights := []Light{&DirectionalLight{Color64{1, 1, 1}, mgl64.Vec3{0, 0, 1}}}
	expColor := Color64{1, 0, 0}
	material := Material{Diffuse: UniformColor{expColor}}
	obj := NewObject(mgl64.Ident4(), &material)
	sphere := Sphere{obj}

	min := -1.0
	max := 1.0
	tests := make([]bdrfStruct, 0)
	for i := 0; i < b.N; i++ {
		x := rand.Float64()*(max-min) - max
		y := rand.Float64()*(max-min) - max
		z := rand.Float64()*(max-min) - max
		point := mgl64.Vec3{x, y, z}
		dir := point.Mul(-1).Normalize()
		ray := NewRay(PrimaryRay, point, dir)
		isect, _ := sphere.Intersect(&ray)
		tests = append(tests, bdrfStruct{&ray, &isect})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LambertianBRDF(lights, tests[i].ray, tests[i].isect)
	}
}
