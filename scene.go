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
	isect, found := scene.Intersect(ray)
	if found {
		return scene.Material[isect.Object.Material()].ShadeBlinnPhong(scene, ray, isect).NRGBA()
	}

	// For fun color wheel:
	// r := uint8((ray.Direction.X() + 1) / 2 * 255)
	// g := uint8((ray.Direction.Y() + 1) / 2 * 255)
	// b := uint8((ray.Direction.Z() + 1) / 2 * 255)

	// No intersection, use background color
	return scene.Camera.Background.NRGBA()
}

// Intersect finds the first object that the given Ray intersects. Found will be
// false if no intersection was found.
func (scene *Scene) Intersect(ray Ray) (isect Intersection, found bool) {
	for _, object := range scene.Objects {
		if i, hit := object.Intersect(ray); hit && (!found || i.T < isect.T) {
			found = true
			isect = i
		}
	}
	return
}
