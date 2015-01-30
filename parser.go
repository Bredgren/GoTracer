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
	MaxDepth int
	AdaptiveThreshold float64
 	DirectionalLights []DirectionalLight
	PointLights []PointLight
 	SpotLights []SpotLight
 	AreaLights []AreaLight
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
			if prop.RotateAxis.Len() == 0.0 {
				prop.RotateAxis = mgl64.Vec3{0, 1, 0}
			}
			rotate := mgl64.HomogRotate3D(mgl64.DegToRad(prop.RotateAngle), prop.RotateAxis.Normalize())
			if prop.Scale.Len() == 0.0 {
				prop.Scale = mgl64.Vec3{1, 1, 1}
			}
			scale := mgl64.Scale3D(prop.Scale[0], prop.Scale[1], prop.Scale[2])
			newTransform := transform.Mul4(translate.Mul4(rotate).Mul4(scale))
			parseSceneObject(subObject, scene, newTransform)
		}
	case "Sphere":
		sphere := SphereObject{}
		sphere.Transform = transform
		sphere.MaterialName = object.Properties.Material
		InitSphereObject(&sphere)
		scene.Objects = append(scene.Objects, &sphere)
	case "Box":
		box := BoxObject{}
		box.Transform = transform
		box.MaterialName = object.Properties.Material
		InitBoxObject(&box)
		scene.Objects = append(scene.Objects, &box)
	case "Square":
		square := SquareObject{}
		square.Transform = transform
		square.MaterialName = object.Properties.Material
		InitSquareObject(&square)
		scene.Objects = append(scene.Objects, &square)
	// case "Triangle":
	// 	tri := TriangleObject{
	// 		transform,
	// 		transform.Inv(),
	// 		object.Properties.Material,
	// 		object.Properties.PointA,
	// 		object.Properties.PointB,
	// 		object.Properties.PointC,
	// 	}
	// 	scene.Objects = append(scene.Objects, tri)
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

	scene := &Scene{}

	scene.Camera = NewCamera(settings.Render.Camera)

	scene.Lights = make([]Light, 0)
	for _, pLight := range settings.Render.PointLights {
		pLight.Scene = scene
		InitPointLight(&pLight)
		scene.Lights = append(scene.Lights, &pLight)
	}
	for _, dLight := range settings.Render.DirectionalLights {
		dLight.Scene = scene
		InitDirectionalLight(&dLight)
		scene.Lights = append(scene.Lights, &dLight)
	}
	for _, sLight := range settings.Render.SpotLights {
		sLight.Scene = scene
		InitSpotLight(&sLight)
		scene.Lights = append(scene.Lights, &sLight)
	}
	for _, aLight := range settings.Render.AreaLights {
		aLight.Scene = scene
		InitAreaLight(&aLight)
		scene.Lights = append(scene.Lights, &aLight)
	}

	scene.Material = make(map[string]Material)
	scene.AmbientLight = settings.Render.AmbientLight
	scene.MaxDepth = settings.Render.MaxDepth
	scene.AdaptiveThreshold = settings.Render.AdaptiveThreshold

	for _, material := range settings.Materials {
		scene.Material[material.Name] = material
	}

	for _, object := range settings.Scene {
		// init identity transform and pass in
		transform := mgl64.Ident4()
		parseSceneObject(object, scene, transform)
	}

	log.Printf(spew.Sdump(scene))
	return scene
}
