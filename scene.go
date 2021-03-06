package gotracer

import (
	"log"
	// "image/color"
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

var (
	ImageHeightDefault       = 1
	ImageWidthDefault        = 1
	AdaptiveThresholdDefault = 0.0
	AAMaxDivisionsDefault    = 0
	AAThresholdDefault       = 0.0
	MaxRecursiveDepthDefault = 1
	AmbientLightDefault      = Color64{0, 0, 0}
	BackgroundDefault        = UniformBackground{Color64{0, 0, 0}}
)

type Scene struct {
	// Camera Camera
	ImageHeight       int
	ImageWidth        int
	AdaptiveThreshold float64
	AAMaxDivisions    int
	AAThreshold       float64
	MaxRecursiveDepth int
	AmbientLight      Color64
	Background        Background
	Lights            []Light
	Objects           []Intersecter
}

// NewScene returns an empty scene with default values.
func NewScene() *Scene {
	return &Scene{
		ImageHeight:       ImageHeightDefault,
		ImageWidth:        ImageWidthDefault,
		AdaptiveThreshold: AdaptiveThresholdDefault,
		AAMaxDivisions:    AAMaxDivisionsDefault,
		AAThreshold:       AAThresholdDefault,
		MaxRecursiveDepth: MaxRecursiveDepthDefault,
		AmbientLight:      AmbientLightDefault,
		Background:        BackgroundDefault,
		Lights:            make([]Light, 0),
		Objects:           make([]Intersecter, 0),
	}
}

// NewSceneFromFile returns a new Scene populated from the settings in the given json
// file. The format of the file should follow the specification of scene/format.txt.
func NewSceneFromFile(fileName string) *Scene {
	var scene *Scene = NewScene()
	var settings interface{} = ReadSettingsFile(fileName)
	SettingParsers["Scene"](scene, settings)
	return scene
}

// func (scene *Scene) TracePixel(x, y int) color.NRGBA {
// 	pixelWidth := 1 / float64(scene.Camera.ImageWidth)
// 	pixelHeight := 1 / float64(scene.Camera.ImageHeight)
// 	centerX := float64(x) * pixelWidth
// 	centerY := float64(y) * pixelHeight
// 	if scene.AAThreshold == 0 {
// 		ray := scene.Camera.RayThrough(centerX, centerY)
// 		return scene.TraceRay(ray, 0, 1.0).NRGBA()
// 	}
// 	halfWidth := pixelWidth / 2
// 	halfHeight := pixelHeight / 2
// 	xMin := centerX - halfWidth
// 	yMin := centerY - halfHeight
// 	xMax := centerX + halfWidth
// 	yMax := centerY + halfHeight
// 	return scene.TraceSubPixel(xMin, yMin, xMax, yMax, 0).NRGBA()
// }

// func (scene *Scene) TraceSubPixel(xMin, yMin, xMax, yMax float64, depth int) Color64 {
// 	width := xMax - xMin
// 	height := yMax - yMin
// 	if depth >= scene.AAMaxDivisions {
// 		x := xMin + 0.5 * width
// 		y := yMin + 0.5 * height
// 		return scene.TraceRay(scene.Camera.RayThrough(x, y), 0, 1.0)
// 	}
// 	x1 := xMin + 0.25 * width
// 	x2 := xMin + 0.75 * width
// 	y1 := yMin + 0.25 * height
// 	y2 := yMin + 0.75 * height
// 	color1 := scene.TraceRay(scene.Camera.RayThrough(x1, y1), 0, 1.0)
// 	color2 := scene.TraceRay(scene.Camera.RayThrough(x2, y1), 0, 1.0)
// 	color3 := scene.TraceRay(scene.Camera.RayThrough(x1, y2), 0, 1.0)
// 	color4 := scene.TraceRay(scene.Camera.RayThrough(x2, y2), 0, 1.0)
// 	thresh := scene.AAThreshold
// 	if ColorsDifferent(color1, color2, thresh) || ColorsDifferent(color1, color3, thresh) ||
// 		ColorsDifferent(color1, color4, thresh) || ColorsDifferent(color2, color3, thresh) ||
// 		ColorsDifferent(color2, color4, thresh) || ColorsDifferent(color3, color4, thresh) {
// 		halfWidth := width / 2
// 		halfHeight := height / 2
// 		d := depth + 1
// 		color1 = scene.TraceSubPixel(xMin, yMin, xMin + halfWidth, yMin + halfHeight, d)
// 		color2 = scene.TraceSubPixel(xMin + halfWidth, yMin, xMax, yMin + halfHeight, d)
// 		color3 = scene.TraceSubPixel(xMin, yMin + halfHeight, xMin + halfWidth, yMax, d)
// 		color4 = scene.TraceSubPixel(xMin + halfWidth, yMin + halfHeight, xMax, yMax, d)
// 	}
// 	sum := mgl64.Vec3(color1).Add(mgl64.Vec3(color2)).Add(mgl64.Vec3(color3)).Add(mgl64.Vec3(color4))
// 	return Color64(sum.Mul(0.25))
// }

// func (scene *Scene) TraceRay(ray Ray, depth int, contribution float64) Color64 {
// 	if depth <= scene.MaxDepth && contribution >= scene.AdaptiveThreshold {
// 		if isect, found := scene.Intersect(ray); found {
// 			material := scene.Material[isect.Object.GetMaterialName()]

// 			exiting := false
// 			insideIndex := material.GetIndexValue(isect)
// 			outsideIndex := AirIndex
// 			if (isect.Normal.Dot(ray.Direction) > 0) {
// 				// Exiting object
// 				insideIndex, outsideIndex = outsideIndex, insideIndex
// 				isect.Normal = isect.Normal.Mul(-1)
// 				exiting = true
// 			}

// 			// Bump mapping
// 			bump := material.GetNormal(isect).Mul(2).Sub(mgl64.Vec3{1, 1, 1})
// 			cosθ := isect.Normal.Dot(mgl64.Vec3{0, 0, 1})
// 			rotateVec := mgl64.Vec3{0, 0, 1}.Cross(isect.Normal)
// 			len2 := rotateVec.X() * rotateVec.X() + rotateVec.Y() * rotateVec.Y() + rotateVec.Z() * rotateVec.Z()
// 			if len2 > Rayε {
// 				rotMat := mgl64.HomogRotate3D(math.Acos(cosθ), rotateVec).Mat3()
// 				isect.Normal = isect.Normal.Add(rotMat.Mul3x1(bump))
// 			} else if isect.Normal.Z() < 0 {
// 				isect.Normal = isect.Normal.Sub(bump)
// 			} else {
// 				isect.Normal = isect.Normal.Add(bump)
// 			}
// 			isect.Normal = isect.Normal.Normalize()

// 			// Direct illumination
// 			illum := material.ShadeBlinnPhong(scene, ray, isect)

// 			// Reflection
// 			reflect := Color64{}
// 			kr := material.GetReflectiveColor(isect)
// 			if kr.Len2() > Rayε {
// 				reflRay := ray.Reflect(isect)
// 				contrib := math.Max(kr.R(), math.Max(kr.G(), kr.B()))
// 				reflColor := scene.TraceRay(reflRay, depth + 1, contrib)
// 				reflect = kr.Product(reflColor)
// 			}
// 			// c := 0.0

// 			// Refraction
// 			refract := Color64{}
// 			kt := material.GetTransmissiveColor(isect)
// 			if kt.Len2() > Rayε {
// 				if !TotalInternalReflection(outsideIndex, insideIndex, isect.Normal, ray.Direction.Mul(-1)) {
// 					refrRay := ray.Refract(isect, outsideIndex, insideIndex)
// 					contrib := math.Max(kt.R(), math.Max(kt.G(), kt.B()))
// 					refrColor := scene.TraceRay(refrRay, depth + 1, contrib)
// 					if exiting {
// 						refract = refrColor.Product(material.BeersTrans(isect))
// 					} else {
// 						refract = refrColor.Product(kt)
// 					}
// 					// c = isect.Normal.Mul(-1).Dot(refrRay.Direction)
// 				} else {
// 					// Total internal reflection
// 					return Color64(mgl64.Vec3(illum).Add(mgl64.Vec3(reflect)))
// 				}
// 			}

// 			// return Color64(mgl64.Vec3(illum).Add(mgl64.Vec3(reflect)).Add(mgl64.Vec3(refract)))

// 			// if !exiting {
// 			// 	c = isect.Normal.Dot(ray.Direction.Mul(-1))
// 			// }
// 			// R0 := math.Pow((insideIndex - AirIndex) / (insideIndex + AirIndex), 2)
// 			// R := R0 + (1 - R0) * math.Pow(1 - c, 5)
// 			// return Color64(mgl64.Vec3(illum).Add(mgl64.Vec3(reflect).Mul(R)).Add(mgl64.Vec3(refract).Mul(1 - R)))

// 			// R := math.Pow((insideIndex - outsideIndex) / (insideIndex + outsideIndex), 2)
// 			// T := (4 * insideIndex * outsideIndex) / math.Pow(insideIndex + outsideIndex, 2)
// 			R := 1.0
// 			T := 1.0
// 			return Color64(mgl64.Vec3(illum).Add(mgl64.Vec3(reflect).Mul(R)).Add(mgl64.Vec3(refract).Mul(T)))
// 		}
// 		// For fun color wheel:
// 		// r := uint8((ray.Direction.X() + 1) / 2 * 255)
// 		// g := uint8((ray.Direction.Y() + 1) / 2 * 255)
// 		// b := uint8((ray.Direction.Z() + 1) / 2 * 255)

// 		// No intersection, use background color
// 		return scene.Camera.Background
// 	}
// 	return Color64{}
// }

// Intersect finds the first object that the given Ray intersects. hit will be
// false if no intersection was found.
func (scene *Scene) Intersect(ray *Ray, isect *Intersection) (found bool) {
	localIsect := Intersection{}
	for _, object := range scene.Objects {
		var inv mgl64.Mat4 = object.GetObject().TransformInv
		localRay, length := ray.Transform(&inv)
		if hit := object.Intersect(&localRay, &localIsect); hit {
			localIsect.T /= length
			if !found || localIsect.T < isect.T {
				found = true
				var t3it mgl64.Mat3 = object.GetObject().T3it
				localIsect.Normal = t3it.Mul3x1(localIsect.Normal).Normalize()
				*isect = localIsect
			}
		}
	}
	return
}

// func TotalInternalReflection(outsideIndex, insideIndex float64, normal, direction mgl64.Vec3) bool {
// 	criticalAngle := math.Asin(insideIndex / outsideIndex)
// 	angle := math.Acos(normal.Dot(direction))
// 	return angle > criticalAngle
// }

func sceneParser(scene *Scene, value interface{}) {
	log.Println("sceneParser", value)
	sceneSettings := value.(map[string]interface{})
	for setting, v := range sceneSettings {
		if setting != "Objects" {
			ParseSetting(scene, setting, v)
		}
	}
	// Objects need to be parsed after materials
	ParseSetting(scene, "Objects", sceneSettings["Objects"])
}

func imageHeightParser(scene *Scene, value interface{}) {
	log.Println("imageHeightParser", value)
	v := int(math.Max(value.(float64), 1))
	scene.ImageHeight = v
}

func imageWidthParser(scene *Scene, value interface{}) {
	log.Println("imageWidthParser", value)
	v := int(math.Max(value.(float64), 1))
	scene.ImageWidth = v
}

func adaptiveThresholdParser(scene *Scene, value interface{}) {
	log.Println("adaptiveThresholdParser", value)
	v := math.Max(value.(float64), 0)
	scene.AdaptiveThreshold = v
}

func aaMaxDivisionsParser(scene *Scene, value interface{}) {
	log.Println("aaMaxDivisionsParser", value)
	v := int(math.Max(value.(float64), 0))
	scene.AAMaxDivisions = v
}

func aaThresholdParser(scene *Scene, value interface{}) {
	log.Println("aaThresholdParser", value)
	v := math.Max(value.(float64), 0)
	scene.AAThreshold = v
}

func maxRecursiveDepthParser(scene *Scene, value interface{}) {
	log.Println("maxRecursiveDepthParser", value)
	v := int(math.Max(value.(float64), 1))
	scene.MaxRecursiveDepth = v
}

func ambientLightParser(scene *Scene, value interface{}) {
	log.Println("ambientLightParser", value)
	scene.AmbientLight = ParseColor64(value)
}

func lightsParser(scene *Scene, value interface{}) {
	log.Println("lightsParser", value)
	lightsList := value.([]interface{})
	for _, lightIface := range lightsList {
		lightMap := lightIface.(map[string]interface{})
		lightType := lightMap["Type"]
		if lightType == nil {
			log.Fatal("No 'Type' specified for light")
		}

		typeName := lightMap["Type"].(string)
		_ = typeName
		ParseSetting(scene, "Lights:"+typeName, lightMap)
	}
}

func init() {
	SettingParsers["Scene"] = sceneParser
	SettingParsers["ImageHeight"] = imageHeightParser
	SettingParsers["ImageWidth"] = imageWidthParser
	SettingParsers["AdaptiveThreshold"] = adaptiveThresholdParser
	SettingParsers["AAMaxDivisions"] = aaMaxDivisionsParser
	SettingParsers["AAThreshold"] = aaThresholdParser
	SettingParsers["MaxRecursiveDepth"] = maxRecursiveDepthParser
	SettingParsers["AmbientLight"] = ambientLightParser
	// SettingParsers["Background"] =
	SettingParsers["Lights"] = lightsParser
	// SettingParsers["Objects"] =
}
