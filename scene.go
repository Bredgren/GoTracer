package raytracer

import (
	"image/color"
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

type Scene struct {
	Camera Camera
	MaxReflection int
	AdaptiveThreshold float64
	AmbientLight Color64
	Lights []Light
	Objects []SceneObject
	Material map[string]Material
}

func (scene *Scene) TracePixel(x, y int) color.NRGBA {
	nx := float64(x) / float64(scene.Camera.ImageWidth)
	ny := float64(y) / float64(scene.Camera.ImageHeight)
	ray := scene.Camera.RayThrough(nx, ny)
	return scene.TraceRay(ray, 0, 1.0).NRGBA()
}

func (scene *Scene) TraceRay(ray Ray, depth int, contribution float64) Color64 {
	if depth <= scene.MaxReflection && contribution >= scene.AdaptiveThreshold {
		if isect, found := scene.Intersect(ray); found {
			material := scene.Material[isect.Object.GetMaterialName()]

			// Direct illumination
			illum := material.ShadeBlinnPhong(scene, ray, isect)

			// Reflection
			reflect := Color64{}
			if mgl64.Vec3(material.Reflective).Len() > Rayε {
				reflRay := ray.Reflect(isect)
				contrib := math.Max(material.Reflective[0], math.Max(material.Reflective[1], material.Reflective[2]))
				reflColor := scene.TraceRay(reflRay, depth + 1, contrib)
				reflect = material.Reflective.Product(reflColor)
			}

			return Color64(mgl64.Vec3(illum).Add(mgl64.Vec3(reflect))).Clamp()
		}
		// For fun color wheel:
		// r := uint8((ray.Direction.X() + 1) / 2 * 255)
		// g := uint8((ray.Direction.Y() + 1) / 2 * 255)
		// b := uint8((ray.Direction.Z() + 1) / 2 * 255)

		// No intersection, use background color
		return scene.Camera.Background
	}
	return Color64{}
}

// Intersect finds the first object that the given Ray intersects. Found will be
// false if no intersection was found.
func (scene *Scene) Intersect(ray Ray) (isect Intersection, found bool) {
	for _, object := range scene.Objects {
		inv := object.GetInvTransform()
		localRay, length := ray.Transform(inv)
		if i, hit := object.Intersect(localRay); hit {
			i.T /= length
			if !found || i.T < isect.T {
				found = true
				normInverse := object.GetTransform().Mat3().Inv().Transpose()
				i.Normal = normInverse.Mul3x1(i.Normal).Normalize()
				isect = i
			}
		}
	}
	return
}
