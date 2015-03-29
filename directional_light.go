package gotracer

import (
	"log"

	"github.com/go-gl/mathgl/mgl64"
)

const (
	DirectionalLightDist = 1e10
)

var (
	DirectionalLightDefaultColor       = Color64{1, 1, 1}
	DirectionalLightDefaultOrientation = mgl64.Vec3{0, -1, 0}
)

// DirectionalLight simulates a point light at an infinate distance. The opposite
// direction to it's orientation is used because more useful and saves work
type DirectionalLight struct {
	Color          Color64
	OrientationInv mgl64.Vec3
	Dir            mgl64.Vec3
}

func (l *DirectionalLight) Attenuation(scene *Scene, point mgl64.Vec3) Color64 {
	return ShadowAttenuation(scene, l.Dir, point)
}

func (l *DirectionalLight) Direction(from mgl64.Vec3) mgl64.Vec3 {
	return l.OrientationInv
}

func directionalLightParser(scene *Scene, value interface{}) {
	log.Println("directionalLightParser", value)
	valueMap := value.(map[string]interface{})
	light := DirectionalLight{
		Color:          DirectionalLightDefaultColor,
		OrientationInv: DirectionalLightDefaultOrientation.Mul(-1),
	}
	for setting, v := range valueMap {
		ParseSetting(scene, "Lights:Directional:"+setting, []interface{}{&light, v})
	}
	light.Dir = light.OrientationInv.Mul(DirectionalLightDist)
	scene.Lights = append(scene.Lights, &light)
}

func directionalLightTypeParser(scene *Scene, value interface{}) {
	log.Println("directionalLightTypeParser", value)
	v := value.([]interface{})
	if v[1].(string) != "Directional" {
		log.Fatal("Parsing a directional light but Type is not 'Directional'")
	}
}

func directionalLightColorParser(scene *Scene, value interface{}) {
	log.Println("directionalLightColorParser", value)
	v := value.([]interface{})
	light := v[0].(*DirectionalLight)
	light.Color = ParseColor64(v[1])
}

func directionalLightOrientationParser(scene *Scene, value interface{}) {
	log.Println("directionalLightOrientationParser", value)
	v := value.([]interface{})
	light := v[0].(*DirectionalLight)
	orient := ParseVector(v[1])
	if mgl64.FloatEqual(orient.Len(), 0) {
		log.Println("Bad orientation, setting to default")
		orient = DirectionalLightDefaultOrientation
	}
	light.OrientationInv = orient.Mul(-1).Normalize()
}

func init() {
	SettingParsers["Lights:Directional"] = directionalLightParser
	SettingParsers["Lights:Directional:Type"] = directionalLightTypeParser
	SettingParsers["Lights:Directional:Color"] = directionalLightColorParser
	SettingParsers["Lights:Directional:Orientation"] = directionalLightOrientationParser
}
