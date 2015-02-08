package raytracer

import (
	"testing"

	"github.com/go-gl/mathgl/mgl64"
)

func TestTexture(t *testing.T) {
	material := &Material{DiffuseTextureFile: "test.png"}
	InitMaterial(material)
	t.Log((*material.DiffuseTexture.tex).At(0, 0))
	t.Log((*material.DiffuseTexture.tex).At(1, 0))
	t.Log((*material.DiffuseTexture.tex).At(0, 1))
	t.Log((*material.DiffuseTexture.tex).At(1, 1))
	isect := Intersection{UVCoords: mgl64.Vec2{0, 0.25}}
	t.Log(material.GetDiffuseColor(isect))
}
