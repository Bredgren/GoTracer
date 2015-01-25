package raytracer

import (
	"math"

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
	colorVec := mgl64.Vec3(m.Emissive).Add(mgl64.Vec3(m.Ambient.Product(scene.AmbientLight)))
	for _, light := range scene.Lights {
		attenuation := light.ShadowAttenuation(point).Mul(light.DistanceAttenuation(point))
		lightDir := light.Direction(point)
		shade := isect.Normal.Dot(lightDir)
		if shade > 0 {
			h := lightDir.Sub(ray.Direction).Normalize()
			s := mgl64.Vec3(m.Specular).Mul(math.Pow(isect.Normal.Dot(h), m.Shininess))
			d := mgl64.Vec3(m.Diffuse).Mul(shade).Add(s)
			a := Color64(attenuation).Product(Color64(d))
			contribution := mgl64.Vec3(light.GetColor().Product(a))
			colorVec = colorVec.Add(contribution)
		}
	}

	return Color64(colorVec)
}
