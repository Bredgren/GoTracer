package gotracer

import (
	"log"

	"github.com/go-gl/mathgl/mgl64"
)

const (
	DirectionalLightDist = 1e10
)

var (
	DirectionalLightDefaultColor = Color64{1, 1, 1}
	DirectionalLightDefaultOrientation = mgl64.Vec3{0, -1, 0}
)

type DirectionalLight struct {
	Color          Color64
	// The opposite direction to it's orientation is more useful
	OrientationInv mgl64.Vec3
}

func (l DirectionalLight) Attenuation(scene *Scene, point mgl64.Vec3) Color64 {
	return ShadowAttenuation(scene, l.OrientationInv.Mul(DirectionalLightDist), point)
}

func (l DirectionalLight) Direction(from mgl64.Vec3) mgl64.Vec3 {
	return l.OrientationInv
}

func directionalLightParser(scene *Scene, value interface{}) {
	log.Println("directionalLightParser", value)
	v := value.(map[string]interface{})
	light := DirectionalLight{
		DirectionalLightDefaultColor,
		DirectionalLightDefaultOrientation,
	}
	for attribute, value := range v {
		switch attribute {
		case "Type":
			if value.(string) != "Directional" {
				log.Fatal("Parsing a directional light but Type is not Directional")
			}
		case "Color":
			light.Color = ParseColor64(value.([]interface{}))
		case "Orientation":
			orient := ParseVector(value.([]interface{}))
			if mgl64.FloatEqual(orient.Len(), 0) {
				orient = DirectionalLightDefaultOrientation
			}
			light.OrientationInv = orient.Mul(-1).Normalize()
		default:
			if attribute[0] != '_' {
				log.Printf("Waning: unknown Directional light attribute '%s'", attribute)
			}
		}
	}
	scene.Lights = append(scene.Lights, light)
}

func init() {
	SettingParsers["Lights/Directional"] = directionalLightParser
}
