package raytracer

import (
	"image/color"
)

type Scene struct {
	Camera Camera
	Lights []Light
	Objects []SceneObject
	Material map[string]Material
}

func (scene *Scene) Trace(x, y int) color.NRGBA {

	return color.NRGBA{uint8(x), uint8(y), uint8(x + y), 255}
}

type Material struct {
	Name string
	Emissive Color64
	Ambient Color64
	Specular Color64
	Reflective Color64
	Diffuse Color64
	Transmissive Color64
	Shininess float64
	Index float64
}
