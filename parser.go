package raytracer

import (
	"encoding/json"
	"log"
	"io/ioutil"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-gl/mathgl/mgl64"
)

type cameraSettings struct {
	ImageWidth int
	ImageHeight int
	Position mgl64.Vec3
	LookAt mgl64.Vec3
	UpDir mgl64.Vec3
	FOV float64
	Background Color64
}

type renderSettings struct {
	ImageWidth int
	ImageHeight int
	Camera cameraSettings
	AmbientLight Color64
	MaxDepth int
	AdaptiveThreshold float64
	PointLights []PointLight
}

type properties struct {
	Translate mgl64.Vec3
	RotateAxis mgl64.Vec3
	RotateAngle float64
	Scale mgl64.Vec3

	Material string

	PointA mgl64.Vec3
	PointB mgl64.Vec3
	PointC mgl64.Vec3
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

func parseSceneObject(object sceneObjectSettings, scene *Scene, transform mgl64.Mat4) {
	switch object.Type {
	case "Transform":
		for _, subObject := range object.SubObjects {
			prop := object.Properties
			translate := mgl64.Translate3D(prop.Translate[0], prop.Translate[1], prop.Translate[2])
			rotate := mgl64.HomogRotate3D(mgl64.DegToRad(prop.RotateAngle), prop.RotateAxis.Normalize())
			if prop.Scale.Len() == 0.0 {
				prop.Scale = mgl64.Vec3{1, 1, 1}
			}
			scale := mgl64.Scale3D(prop.Scale[0], prop.Scale[1], prop.Scale[2])
			newTransform := transform.Mul4(translate.Mul4(rotate).Mul4(scale))
			parseSceneObject(subObject, scene, newTransform)
		}
	case "Sphere":
		sphere := SphereObject{transform, transform.Inv(), object.Properties.Material}
		scene.Objects = append(scene.Objects, sphere)
	case "Box":
		box := BoxObject{transform, transform.Inv(), object.Properties.Material}
		scene.Objects = append(scene.Objects, box)
	case "Triangle":
		tri := TriangleObject{
			transform,
			transform.Inv(),
			object.Properties.Material,
			object.Properties.PointA,
			object.Properties.PointB,
			object.Properties.PointC,
		}
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

	camSet := settings.Render.Camera
	cam := NewCamera(camSet.ImageWidth, camSet.ImageHeight, camSet.Position,
		camSet.LookAt, camSet.UpDir, camSet.FOV, camSet.Background)

	lights := make([]Light, 0)
	for _, light := range settings.Render.PointLights {
		lights = append(lights, light)
	}

	scene := Scene{
		Camera: cam,
		Material: make(map[string]Material),
		Lights: lights,
		AmbientLight: settings.Render.AmbientLight,
		MaxDepth: settings.Render.MaxDepth,
		AdaptiveThreshold: settings.Render.AdaptiveThreshold,
	}
	scene.Camera.Update()

	for _, material := range settings.Materials {
		scene.Material[material.Name] = material
	}

	for _, object := range settings.Scene {
		// init identity transform and pass in
		transform := mgl64.Ident4()
		parseSceneObject(object, &scene, transform)
	}

	log.Printf(spew.Sdump(scene))
	return &scene
}
