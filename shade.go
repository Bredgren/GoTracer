package gotracer

type ShaderParams struct {
	Lights []Light
	Ray    Ray
	Isect  Intersection
	Mat    *Material
}

type Shader interface {
	Shade(ShaderParams) Color64
}

type BlinnPhongShader struct {
}

func (s *BlinnPhongShader) Shade(params ShaderParams) (color Color64) {
	return Color64{0, 0, 0}
}

// func (m *Material) ShadeBlinnPhong(scene *Scene, ray Ray, isect Intersection) (color Color64) {
// 	point := ray.At(isect.T)
// 	colorVec := mgl64.Vec3(m.GetEmissiveColor(isect)).Add(mgl64.Vec3(m.GetAmbientColor(isect).Product(scene.AmbientLight)))
// 	for _, light := range scene.Lights {
// 		attenuation := light.ShadowAttenuation(point).Mul(light.DistanceAttenuation(point))
// 		lightDir := light.Direction(point)
// 		shade := isect.Normal.Dot(lightDir)
// 		if shade > 0 {
// 			h := lightDir.Sub(ray.Direction).Normalize()
// 			s := mgl64.Vec3(m.GetSpecularColor(isect)).Mul(math.Pow(isect.Normal.Dot(h), m.GetShininessValue(isect)))
// 			d := mgl64.Vec3(m.GetDiffuseColor(isect)).Mul(shade).Add(s)
// 			a := Color64(attenuation).Product(Color64(d))
// 			contribution := mgl64.Vec3(light.GetColor().Product(a))
// 			colorVec = colorVec.Add(contribution)
// 		}
// 	}

// 	return Color64(colorVec)
// }
