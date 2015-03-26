package gotracer

import (
	"testing"

	"github.com/go-gl/mathgl/mgl64"
)

func TestShadowAttenuation(t *testing.T) {
	lightDir := mgl64.Vec3{0, 0, 1}.Mul(DirectionalLightDist)
	scene := Scene{}
	point := mgl64.Vec3{0, 0, -2}
	expAtten := Color64{1, 1, 1}

	atten := ShadowAttenuation(&scene, lightDir, point)
	if atten != expAtten {
		t.Errorf("Incorrect attenuation. Expected %v, got %v", expAtten, atten)
	}

	material := Material{
		Transmissive: UniformColor{Color64{}},
	}
	sphere := Sphere{NewObject(mgl64.Ident4(), &material)}
	scene.Objects = append(scene.Objects, sphere)
	expAtten = Color64{0, 0, 0}
	atten = ShadowAttenuation(&scene, lightDir, point)
	if atten != expAtten {
		t.Errorf("Incorrect attenuation. Expected %v, got %v", expAtten, atten)
	}
}
