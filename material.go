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

func Clamp(low, high float64) func(float64) float64 {
	return func(v float64) float64 {
		return math.Min(high, math.Max(low, v))
	}
}

var Clamp01 = Clamp(0.0, 1.0)

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

	color = Color64{Clamp01(colorVec[0]), Clamp01(colorVec[1]), Clamp01(colorVec[2])}
	return color
}
