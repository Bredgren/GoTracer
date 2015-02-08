package raytracer

import (
	"testing"

	"github.com/go-gl/mathgl/mgl64"
)

func TestAreaLight(t *testing.T) {
	scene := Scene{}
	scene.Material = make(map[string]*Material)
	scene.Material["1"] = &Material{}
	light := AreaLight{
		Scene: &scene,
		Color: Color64{1, 1, 1},
		Position: mgl64.Vec3{0, 4, 0},
		Orientation: mgl64.Vec3{0, -1, 0},
		UpDir: mgl64.Vec3{0, 0, -1},
		Size: 1,
		Samples: 4,
	}
	InitAreaLight(&light)
	u := mgl64.Vec3{1, 0, 0}
	if !light.u.ApproxEqual(u) {
		t.Errorf("Bad u %v, expected %v", light.u, u)
	}
	v := mgl64.Vec3{0, 0, -1}
	if !light.v.ApproxEqual(v) {
		t.Errorf("Bad v %v, expected %v", light.v, v)
	}

	point := mgl64.Vec3{0, 0, 0}
	expAtten := mgl64.Vec3{1, 1, 1}
	for samples := 1; samples < 50; samples++ {
		atten := light.GridAttenuation(point, samples)
		if !atten.ApproxEqual(expAtten) {
			t.Errorf("(1.0) Bad attenuation %v, expected %v (%v samples)", atten, expAtten, samples)
		}
	}

	light.Size = 0.3
	InitAreaLight(&light)

	point = mgl64.Vec3{0, 0, 0}
	expAtten = mgl64.Vec3{1, 1, 1}
	for samples := 1; samples < 50; samples++ {
		atten := light.GridAttenuation(point, samples)
		if !atten.ApproxEqual(expAtten) {
			t.Errorf("(0.3) Bad attenuation %v, expected %v (%v samples)", atten, expAtten, samples)
		}
	}

	light.Size = 0.4
	InitAreaLight(&light)

	translate := mgl64.Translate3D(0.5, 4 - Rayε, 0)
	rotate := mgl64.HomogRotate3D(mgl64.DegToRad(90), mgl64.Vec3{1, 0, 0})
	scale := mgl64.Scale3D(1, 1, 1)
	transform := translate.Mul4(rotate).Mul4(scale)

	square := SquareObject{Transform: transform, MaterialName: "1"}
	InitSquareObject(&square)
	scene.Objects = append(scene.Objects, square)

	t.Log("start 0.4")
	samples := 2
	atten := light.GridAttenuation(point, samples)
	expAtten = mgl64.Vec3{0.5, 0.5, 0.5}
	if !atten.ApproxEqual(expAtten) {
		t.Errorf("(0.4) Bad attenuation %v, expected %v (%v samples)", atten, expAtten, samples)
	}
}
