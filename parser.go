package gotracer

import (
	"log"

	"github.com/go-gl/mathgl/mgl64"
)

type SceneSettings map[string]interface{}

type Parser interface{
	Parse(scene *Scene, value interface{})
	GetDependencies() []string
}

var SettingParsers map[string]Parser = make(map[string]Parser)

func dependenciesSatisfied(deps []string, parsed map[string]bool) bool {
	for _, dep := range deps {
		if !parsed[dep] {
			return false
		}
	}
	return true
}

func ParseSubsetting(name string, scene *Scene, settings map[string]interface{}) {
	if parser := SettingParsers[name]; parser != nil {
		if len(parser.GetDependencies()) > 0 {
			log.Fatalf("Subsettings not allowed to have dependencies. Check '%v'", name)
		}
		parser.Parse(scene, settings)
	} else {
		log.Printf("Warning: unknown subsetting '%s'", name)
	}
}

func ParseSettings(settings SceneSettings) *Scene {
	log.Println("ParseSettings", settings)
	// Keep track of which top level settings we've parsed to allow dependencies. This
	// allows settings to ensure that other settings get handled fist (such as objects
	// makeing sure materials are parsed first).
	parsed := make(map[string]bool, len(settings))
	for k := range settings {
		parsed[k] = false
	}

	var scene *Scene = NewScene()

Outer:
	for attribute, beenParsed := range parsed {
		if beenParsed {
			continue
		}
		var parser Parser = SettingParsers[attribute]
		if parser == nil {
			log.Printf("Warning: unknown attribute '%s'", attribute)
			parsed[attribute] = true
			goto Outer
		}
		if dependenciesSatisfied(parser.GetDependencies(), parsed) {
			parser.Parse(scene, settings[attribute])
			parsed[attribute] = true
			goto Outer
		}
		// Dependencies for this attribute aren't satisfied yet, check the next one
	}

	// If we get here either all have been parsed and we're done or some have
	// impossible dependencies.
	for attribute, beenParsed := range parsed {
		if !beenParsed {
			var deps []string = SettingParsers[attribute].GetDependencies()
			log.Fatalf("Attribute %v has impossible dependencies %v", attribute, deps)
		}
	}
	log.Println("Parsing done")

	return scene
}

// ParseColor64 takes an []interface{} which it assumes is actually a [3]float64
// and converts it to a Color64.
func ParseColor64(floatArray []interface{}) Color64 {
	return Color64{
		floatArray[0].(float64),
		floatArray[1].(float64),
		floatArray[2].(float64),
	}
}

// ParseVector takes an []interface{} which it assumes is actually a [3]float64
// and converts it to a mgl64.Vec3.
func ParseVector(floatArray []interface{}) mgl64.Vec3 {
	return mgl64.Vec3{
		floatArray[0].(float64),
		floatArray[1].(float64),
		floatArray[2].(float64),
	}
}

// type renderSettings struct {
// 	ImageWidth int
// 	ImageHeight int
// 	Camera Camera
// 	AmbientLight Color64
// 	MaxDepth int
// 	AdaptiveThreshold float64
// 	AAMaxDivisions int
// 	AAThreshold float64
//  	DirectionalLights []DirectionalLight
// 	PointLights []PointLight
//  	SpotLights []SpotLight
//  	AreaLights []AreaLight
// }

// type properties struct {
// 	Translate mgl64.Vec3
// 	RotateAxis mgl64.Vec3
// 	RotateAngle float64
// 	Scale mgl64.Vec3

// 	Material string

// 	Capped bool
// 	TopRadius float64
// 	BaseRadius float64
// }

// type sceneObjectSettings struct {
// 	Type string
// 	Properties properties
// 	SubObjects []sceneObjectSettings
// }

// type parsedSettings struct {
// 	Render renderSettings
// 	Materials []*Material
// 	Scene []sceneObjectSettings
// }

// func parseSceneObject(object sceneObjectSettings, scene *Scene, transform mgl64.Mat4) {
// 	switch object.Type {
// 	case "Transform":
// 		for _, subObject := range object.SubObjects {
// 			prop := object.Properties
// 			translate := mgl64.Translate3D(prop.Translate[0], prop.Translate[1], prop.Translate[2])
// 			if prop.RotateAxis.Len() == 0.0 {
// 				prop.RotateAxis = mgl64.Vec3{0, 1, 0}
// 			}
// 			rotate := mgl64.HomogRotate3D(mgl64.DegToRad(prop.RotateAngle), prop.RotateAxis.Normalize())
// 			if prop.Scale.Len() == 0.0 {
// 				prop.Scale = mgl64.Vec3{1, 1, 1}
// 			}
// 			scale := mgl64.Scale3D(prop.Scale[0], prop.Scale[1], prop.Scale[2])
// 			newTransform := transform.Mul4(translate.Mul4(rotate).Mul4(scale))
// 			parseSceneObject(subObject, scene, newTransform)
// 		}
// 	case "Sphere":
// 		sphere := SphereObject{}
// 		sphere.Transform = transform
// 		sphere.MaterialName = object.Properties.Material
// 		InitSphereObject(&sphere)
// 		scene.Objects = append(scene.Objects, &sphere)
// 	case "Box":
// 		box := BoxObject{}
// 		box.Transform = transform
// 		box.MaterialName = object.Properties.Material
// 		InitBoxObject(&box)
// 		scene.Objects = append(scene.Objects, &box)
// 	case "Square":
// 		square := SquareObject{}
// 		square.Transform = transform
// 		square.MaterialName = object.Properties.Material
// 		InitSquareObject(&square)
// 		scene.Objects = append(scene.Objects, &square)
// 	case "Cylinder":
// 		cylinder := CylinderObject{}
// 		cylinder.Transform = transform
// 		cylinder.MaterialName = object.Properties.Material
// 		cylinder.Capped = object.Properties.Capped
// 		InitCylinderObject(&cylinder)
// 		scene.Objects = append(scene.Objects, &cylinder)
// 	case "Cone":
// 		cone := ConeObject{}
// 		cone.Transform = transform
// 		cone.MaterialName = object.Properties.Material
// 		cone.Capped = object.Properties.Capped
// 		cone.TopRadius = object.Properties.TopRadius
// 		cone.BaseRadius = object.Properties.BaseRadius
// 		InitConeObject(&cone)
// 		scene.Objects = append(scene.Objects, &cone)
// 	// case "Triangle":
// 	// 	tri := TriangleObject{
// 	// 		transform,
// 	// 		transform.Inv(),
// 	// 		object.Properties.Material,
// 	// 		object.Properties.PointA,
// 	// 		object.Properties.PointB,
// 	// 		object.Properties.PointC,
// 	// 	}
// 	// 	scene.Objects = append(scene.Objects, tri)
// 	}
// }

// func Parse(fileName string) *Scene {
// 	log.Println("Parse", fileName)

// 	file, err := ioutil.ReadFile(fileName)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	settings := parsedSettings{}
// 	err = json.Unmarshal(file, &settings)
// 	if err != nil {
// 		log.Fatalln("Decoding JSON:", err)
// 	}

// 	scene := &Scene{}

// 	scene.Camera = NewCamera(settings.Render.Camera)

// 	scene.Lights = make([]Light, 0)
// 	for _, pLight := range settings.Render.PointLights {
// 		pLight.Scene = scene
// 		InitPointLight(&pLight)
// 		scene.Lights = append(scene.Lights, &pLight)
// 	}
// 	for _, dLight := range settings.Render.DirectionalLights {
// 		dLight.Scene = scene
// 		InitDirectionalLight(&dLight)
// 		scene.Lights = append(scene.Lights, &dLight)
// 	}
// 	for _, sLight := range settings.Render.SpotLights {
// 		sLight.Scene = scene
// 		InitSpotLight(&sLight)
// 		scene.Lights = append(scene.Lights, &sLight)
// 	}
// 	for _, aLight := range settings.Render.AreaLights {
// 		aLight.Scene = scene
// 		InitAreaLight(&aLight)
// 		scene.Lights = append(scene.Lights, &aLight)
// 	}

// 	scene.Material = make(map[string]*Material)
// 	scene.AmbientLight = settings.Render.AmbientLight
// 	scene.MaxDepth = settings.Render.MaxDepth
// 	scene.AdaptiveThreshold = settings.Render.AdaptiveThreshold
// 	scene.AAMaxDivisions = settings.Render.AAMaxDivisions
// 	scene.AAThreshold = settings.Render.AAThreshold

// 	for _, material := range settings.Materials {
// 		InitMaterial(material)
// 		scene.Material[material.Name] = material
// 	}

// 	for _, object := range settings.Scene {
// 		// init identity transform and pass it in
// 		transform := mgl64.Ident4()
// 		parseSceneObject(object, scene, transform)
// 	}

// 	spew.Config.MaxDepth = 3
// 	log.Printf(spew.Sdump(scene))
// 	return scene
// }
