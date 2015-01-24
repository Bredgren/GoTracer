package raytracer

import (
	"github.com/go-gl/mathgl/mgl64"
)

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

func (m Material) ShadeBlinnPhong(scene *Scene, ray Ray, isect Intersection) (color Color64) {
	point := ray.At(isect.T)
	for _, light := range scene.Lights {
		lightDir := light.Direction(point)
		if isect.Normal.Dot(lightDir) <= 0.0 {
			// Normal points away from light
			continue
		}
		atten := light.DistanceAttenuation(point)
		l := mgl64.Vec3(m.Diffuse).Mul(lightDir.Dot(isect.Normal))
		lambert := Color64(l.Mul(atten))
		color = Color64(mgl64.Vec3(color).Add(mgl64.Vec3(light.GetColor().Product(lambert))))
	}
	return color
}
