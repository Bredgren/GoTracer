package gotracer

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

// func init() {
// 	parser.RegisterBRDF("lambertian", LambertianBRDF)
//  ...
// }

// A bidirectional reflectance distribution function.
// Returns the ratio of reflected radiance exiting along wr to the irradiance
// incident on the surface from direction wi.
// wi - normalized negative incomming light direction (toward light)
// wr - normalized outgoing direction (toward camera)
// type BRDF func(wi, wr mgl64.Vec3, isect Intersection) mgl64.Vec3
type BRDF func(scene *Scene, ray *Ray, isect *Intersection) (color Color64)

// func LambertianBRDF(wi, wr mgl64.Vec3, isect Intersection) mgl64.Vec3 {
// 	kd := isect.Object.Material.Diffuse.ColorAt(isect.UVCoords)
// 	return mgl64.Vec3(kd).Mul(wi.Dot(isect.Normal))
// }
func LambertianBRDF(scene *Scene, ray *Ray, isect *Intersection) (color Color64) {
	var point mgl64.Vec3 = ray.At(isect.T)
	var kd Color64 = isect.Object.Material.Diffuse.ColorAt(isect.UVCoords)
	for _, light := range scene.Lights {
		var cosTheta float64 = isect.Normal.Dot(light.Direction(point))
		if cosTheta > 0 {
			var diffuse Color64 = kd.Mul(cosTheta)
			var contribution Color64 = diffuse.Product(light.Attenuation(scene, point))
			color = color.Add(contribution)
		}
	}
	return
}

// func BlinnPhongBRDF(wi, wr mgl64.Vec3, isect Intersection) mgl64.Vec3 {
// 	// diffuse := LambertianBRDF(wi, wr, isect)
// 	// specular := ...
// 	return mgl64.Vec3{0, 0, 0}
// }
func BlinnPhongBRDF(scene *Scene, ray *Ray, isect *Intersection) (color Color64) {
	var point mgl64.Vec3 = ray.At(isect.T)
	var mat *Material = isect.Object.Material
	var ke Color64 = mat.Emissive.ColorAt(isect.UVCoords)
	var ka Color64 = mat.Ambient.ColorAt(isect.UVCoords)
	var kd Color64 = mat.Diffuse.ColorAt(isect.UVCoords)
	var ks Color64 = mat.Specular.ColorAt(isect.UVCoords)
	var shininess float64 = mat.Smoothness.ColorAt(isect.UVCoords).R() * 100
	color = ke.Add(ka.Product(scene.AmbientLight))
	for _, light := range scene.Lights {
		var lightDir mgl64.Vec3 = light.Direction(point)
		var cosTheta float64 = isect.Normal.Dot(lightDir)
		if cosTheta > 0 {
			var h mgl64.Vec3 = lightDir.Sub(ray.Dir).Normalize()
			var specular Color64 = ks.Mul(math.Pow(isect.Normal.Dot(h), shininess))
			var diffuse Color64 = kd.Mul(cosTheta)
			var contribution Color64 = diffuse.Add(specular)
			color = color.Add(contribution.Product(light.Attenuation(scene, point)))
		}
	}
	return
}

// func TorranceSparrowBRDF(wi, wr mgl64.Vec3, isect Intersection) mgl64.Vec3 {
// 	return mgl64.Vec3{0, 0, 0}
// }

// Calculate Lr(wr) = sum for each light i, Li(wi) * BRDF(wi, wr) * (wi . isect.Normal)
// where wr = -ray, wi = direction to light, Li(wi) = radiance from light source,
// func Shade(lights []*Light, ray Ray, isect Intersection) (color Color64) {
// var point mgl64.Vec3 = ray.At(params.isect.T)
// uv := isect.UVCoords
// colorVec := mgl64.Vec3(m.Emissive.ColorAt(uv)).Add(mgl64.Vec3(m.Ambient.ColorAt(uv).Product(scene.AmbientLight)))
// for _, light := range scene.Lights {
// 	attenuation := light.ShadowAttenuation(point).Mul(light.DistanceAttenuation(point))
// 	lightDir := light.Direction(point)
// 	shade := isect.Normal.Dot(lightDir)
// 	if shade > 0 {
// 		h := lightDir.Sub(ray.Direction).Normalize()
// 		s := mgl64.Vec3(m.Specular.ColorAt(uv)).Mul(math.Pow(isect.Normal.Dot(h), m.Gloss.ColorAt(uv)))
// 		d := mgl64.Vec3(m.Diffuse.ColorAt(uv)).Mul(shade).Add(s)
// 		a := Color64(attenuation).Product(Color64(d))
// 		contribution := mgl64.Vec3(light.GetColor().Product(a))
// 		colorVec = colorVec.Add(contribution)
// 	}
// }

// return Color64(colorVec)
// 	return Color64{0, 0, 0}
// }

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
