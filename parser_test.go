package gotracer

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	// "github.com/go-gl/mathgl/mgl64"
	"github.com/Bredgren/misc"
	"github.com/davecgh/go-spew/spew"
)

func TestParser(t *testing.T) {
	settings := make(SceneSettings)

	file, err := ioutil.ReadFile("scene/example.json")
	misc.Check(err)

	t.Log(settings)
	err = json.Unmarshal(file, &settings)
	misc.Check(err)
	t.Log(settings)

	scene := ParseSettings(settings)
	t.Log("Parsed scene:", spew.Sdump(scene))
}
