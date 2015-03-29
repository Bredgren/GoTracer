package gotracer

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/Bredgren/misc"
	"github.com/go-gl/mathgl/mgl64"
)

type Parser func(scene *Scene, value interface{})

var SettingParsers map[string]Parser = make(map[string]Parser)

func ReadSettingsFile(fileName string) interface{} {
	log.Println("ReadSettingsFile", fileName)
	contents, err := ioutil.ReadFile(fileName)
	misc.Check(err)

	var settings interface{}
	err = json.Unmarshal(contents, &settings)
	misc.Check(err)

	return settings
}

func ParseSetting(scene *Scene, setting string, value interface{}) {
	if setting[0] == '_' {
		return
	}
	if fn := SettingParsers[setting]; fn == nil {
		log.Printf("Warning: unknown setting '%s'", setting)
	} else {
		fn(scene, value)
	}
}

// ParseColor64 takes an interface{} which it assumes is actually a [3]float64
// and converts it to a Color64.
func ParseColor64(floatArray interface{}) Color64 {
	a := floatArray.([]interface{})
	return Color64{a[0].(float64), a[1].(float64), a[2].(float64)}
}

// ParseVector takes an interface{} which it assumes is actually a [3]float64
// and converts it to a mgl64.Vec3.
func ParseVector(floatArray interface{}) mgl64.Vec3 {
	a := floatArray.([]interface{})
	return mgl64.Vec3{a[0].(float64), a[1].(float64), a[2].(float64)}
}
