package raytracer

import (
	"image/color"
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

type Scene struct {
	Camera Camera
	MaxDepth int
	AdaptiveThreshold float64
  AAMaxDivisions int
	AAThreshold float64
	AmbientLight Color64
	Lights []Light
	Objects []SceneObject
	Material map[string]Material
}

func (scene *Scene) TracePixel(x, y int) color.NRGBA {
	pixelWidth := 1 / float64(scene.Camera.ImageWidth)
	pixelHeight := 1 / float64(scene.Camera.ImageHeight)
	centerX := float64(x) * pixelWidth
	centerY := float64(y) * pixelHeight
	if scene.AAThreshold == 0 {
		ray := scene.Camera.RayThrough(centerX, centerY)
		return scene.TraceRay(ray, 0, 1.0).NRGBA()
	}
	halfWidth := pixelWidth / 2
	halfHeight := pixelHeight / 2
	xMin := centerX - halfWidth
	yMin := centerY - halfHeight
	xMax := centerX + halfWidth
	yMax := centerY + halfHeight
	return scene.TraceSubPixel(xMin, yMin, xMax, yMax, 0).NRGBA()
}

func (scene *Scene) TraceSubPixel(xMin, yMin, xMax, yMax float64, depth int) Color64 {
	width := xMax - xMin
	height := yMax - yMin
	if depth >= scene.AAMaxDivisions {
		x := xMin + 0.5 * width
		y := yMin + 0.5 * height
		return scene.TraceRay(scene.Camera.RayThrough(x, y), 0, 1.0)
	}
	x1 := xMin + 0.25 * width
	x2 := xMin + 0.75 * width
	y1 := yMin + 0.25 * height
	y2 := yMin + 0.75 * height
	color1 := scene.TraceRay(scene.Camera.RayThrough(x1, y1), 0, 1.0)
	color2 := scene.TraceRay(scene.Camera.RayThrough(x2, y1), 0, 1.0)
	color3 := scene.TraceRay(scene.Camera.RayThrough(x1, y2), 0, 1.0)
	color4 := scene.TraceRay(scene.Camera.RayThrough(x2, y2), 0, 1.0)
	thresh := scene.AAThreshold
	if ColorsDifferent(color1, color2, thresh) || ColorsDifferent(color1, color3, thresh) ||
		ColorsDifferent(color1, color4, thresh) || ColorsDifferent(color2, color3, thresh) ||
		ColorsDifferent(color2, color4, thresh) || ColorsDifferent(color3, color4, thresh) {
		halfWidth := width / 2
		halfHeight := height / 2
		d := depth + 1
		color1 = scene.TraceSubPixel(xMin, yMin, xMin + halfWidth, yMin + halfHeight, d)
		color2 = scene.TraceSubPixel(xMin + halfWidth, yMin, xMax, yMin + halfHeight, d)
		color3 = scene.TraceSubPixel(xMin, yMin + halfHeight, xMin + halfWidth, yMax, d)
		color4 = scene.TraceSubPixel(xMin + halfWidth, yMin + halfHeight, xMax, yMax, d)
	}
	sum := mgl64.Vec3(color1).Add(mgl64.Vec3(color2)).Add(mgl64.Vec3(color3)).Add(mgl64.Vec3(color4))
	return Color64(sum.Mul(0.25))
}

func (scene *Scene) TraceRay(ray Ray, depth int, contribution float64) Color64 {
	if depth <= scene.MaxDepth && contribution >= scene.AdaptiveThreshold {
		if isect, found := scene.Intersect(ray); found {
			material := scene.Material[isect.Object.GetMaterialName()]

			// Direct illumination
			illum := material.ShadeBlinnPhong(scene, ray, isect)

			entering := true
			insideIndex := material.Index
			outsideIndex := AirIndex
			if (isect.Normal.Dot(ray.Direction.Mul(-1)) < 0) {
				// Exiting object
				insideIndex, outsideIndex = outsideIndex, insideIndex
				isect.Normal = isect.Normal.Mul(-1)
				entering = false
				illum = Color64{}
			}
			_ = entering

			// Reflection
			reflect := Color64{}
			if mgl64.Vec3(material.Reflective).Len() > Rayε {
				reflRay := ray.Reflect(isect)
				contrib := math.Max(material.Reflective[0], math.Max(material.Reflective[1], material.Reflective[2]))
				reflColor := scene.TraceRay(reflRay, depth + 1, contrib)
				reflect = material.Reflective.Product(reflColor)
			}

			// Refraction
			refract := Color64{}
			if mgl64.Vec3(material.Transmissive).Len() > Rayε {
				if !TotalInternalReflection(outsideIndex, insideIndex, isect.Normal, ray.Direction.Mul(-1)) {
					refrRay := ray.Refract(isect, outsideIndex, insideIndex)
					contrib := math.Max(material.Transmissive[0], math.Max(material.Transmissive[1], material.Transmissive[2]))
					refrColor := scene.TraceRay(refrRay, depth + 1, contrib)
					refract = material.Transmissive.Product(refrColor)
					// TODO: Fresnel
				}
			}

			return Color64(mgl64.Vec3(illum).Add(mgl64.Vec3(reflect)).Add(mgl64.Vec3(refract)))
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

func TotalInternalReflection(outsideIndex, insideIndex float64, normal, direction mgl64.Vec3) bool {
	criticalAngle := math.Asin(insideIndex / outsideIndex)
	angle := math.Acos(normal.Dot(direction))
	return angle > criticalAngle
}
