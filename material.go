package raytracer

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

const (
	AirIndex = 1.0003
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

	LogTransmissive Color64
}

func InitMaterial(m *Material) {
	m.LogTransmissive = Color64{
		math.Log(2 - m.Transmissive[0]),
		math.Log(2 - m.Transmissive[1]),
		math.Log(2 - m.Transmissive[2]),
	}
}

func (m *Material) BeersTrans(dist float64) Color64 {
	return Color64{
		math.Exp(m.LogTransmissive[0] * -dist),
		math.Exp(m.LogTransmissive[1] * -dist),
		math.Exp(m.LogTransmissive[2] * -dist),
	}
}

func (m *Material) ShadeBlinnPhong(scene *Scene, ray Ray, isect Intersection) (color Color64) {
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
