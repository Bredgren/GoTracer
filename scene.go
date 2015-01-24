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

func (scene *Scene) TracePixel(x, y int) color.NRGBA {
	nx := float64(x) / float64(scene.Camera.ImageWidth)
	ny := float64(y) / float64(scene.Camera.ImageHeight)
	ray := scene.Camera.RayThrough(nx, ny)
	return scene.TraceRay(ray)
}

func (scene *Scene) TraceRay(ray Ray) color.NRGBA {
	r := uint8((ray.Direction.X() + 1) / 2 * 255)
	g := uint8((ray.Direction.Y() + 1) / 2 * 255)
	b := uint8((ray.Direction.Z() + 1) / 2 * 255)
	return color.NRGBA{r, g, b, 255}

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
