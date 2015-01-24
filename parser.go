package raytracer

import (
	"encoding/json"
	"log"
	"io/ioutil"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-gl/mathgl/mgl64"
)

type renderSettings struct {
	ImageWidth int
	ImageHeight int
	Camera Camera
	AmbientLight Color64
}

type transformProperties struct {
	Translate mgl64.Vec3
	Rotate mgl64.Vec3
	Scale mgl64.Vec3
	SubObjects []sceneObjectSettings
}

type properties struct {
	Translate mgl64.Vec3
	Rotate mgl64.Vec3
	Scale mgl64.Vec3
	PointA mgl64.Vec3
	PointB mgl64.Vec3
	PointC mgl64.Vec3
	Radius float64
	Material string
}

type sceneObjectSettings struct {
	Type string
	Properties properties
	SubObjects []sceneObjectSettings
}

type parsedSettings struct {
	Render renderSettings
	Materials []Material
	Scene []sceneObjectSettings
}

func parseSceneObject(object sceneObjectSettings, scene *Scene) {
	switch object.Type {
	case "Transform":
		for _, subObject := range object.SubObjects {
			// pass in current transform modified by this one
			parseSceneObject(subObject, scene)
		}
	case "Sphere":
		sphere := SphereObject{}
		sphere.Center = object.Properties.PointA
		sphere.Radius = object.Properties.Radius
		sphere.MaterialName = object.Properties.Material
		// pass in current transform
		scene.Objects = append(scene.Objects, sphere)
	case "Triangle":
		tri := TriangleObject{}
		tri.PointA = object.Properties.PointA
		tri.PointB = object.Properties.PointB
		tri.PointC = object.Properties.PointC
		tri.MaterialName = object.Properties.Material
		// pass in current transform
		scene.Objects = append(scene.Objects, tri)
	}
}

func Parse(fileName string) *Scene {
	log.Println("Parse", fileName)

	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	settings := parsedSettings{}
	err = json.Unmarshal(file, &settings)
	if err != nil {
		log.Fatalln("Decoding JSON:", err)
	}
	log.Printf(spew.Sdump(settings))

	scene := Scene{
		Camera: settings.Render.Camera,
		Material: make(map[string]Material),
	}
	scene.Camera.Update()

	for _, material := range settings.Materials {
		scene.Material[material.Name] = material
	}

	for _, object := range settings.Scene {
		// init identity transform and pass in
		parseSceneObject(object, &scene)
	}

	log.Printf(spew.Sdump(scene))
	return &scene
}
