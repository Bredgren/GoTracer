package gotracer

import (
	"encoding/json"
	"testing"
	"io/ioutil"

	// "github.com/go-gl/mathgl/mgl64"
	"github.com/Bredgren/misc"
)

func TestParser(t *testing.T) {
	settings := make(SceneSettings)
	// settings["AmbientLight"] = [3]float64{1, 2, 3}
	// settings["Fake"] = [3]float64{1, 2, 3}
	// settings["Lights"] = []struct{
	// 	Type string
	// 	Color [3]float64
	// 	Orientation [3]float64
	// }{
	// 	{"light1", [3]float64{1, 2, 3}, [3]float64{.1, .2, .3}},
	// 	{"light2", [3]float64{4, 5, 6}, [3]float64{.4, .5, .6}},
	// }

	file, err := ioutil.ReadFile("scene/example.json")
	misc.Check(err)

	t.Log(settings)
	err = json.Unmarshal(file, &settings)
	misc.Check(err)
	t.Log(settings)

	scene := ParseSettings(settings)
	t.Log(scene)

	// settings := make(sceneSettings)
	// settings["Materials"] = make([]sceneSettings, 0)

	// matSet.Parameters["Emissive"] = mgl64.Vec3{1 ,2, 3}
	// settings.Materials = append(settings.Materials, matSet)

	// res, err := json.MarshalIndent(settings, "", "  ")
	// if err != nil {
	// 	t.Error(err)
	// }
	// t.Log(string(res))

	// settings2 := make(map[string]interface{})
	// err = json.Unmarshal(res, &settings2)
	// if err != nil {
	// 	t.Error("Decoding JSON:", err)
	// }
	// t.Log(settings2)
}
