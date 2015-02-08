package raytracer

import (
	"image"
	_ "image/png"
	_ "image/jpeg"
	"log"
	"math"
	"os"

	"github.com/go-gl/mathgl/mgl64"
)

const (
	AirIndex = 1.0003
)

var textures map[string]*image.Image = make(map[string]*image.Image)

type Texture struct {
	tex *image.Image
	Width float64
	Height float64
}

func NewTexture(fileName string) *Texture {
	log.Printf("NewTexture(%s)", fileName)
	if fileName == "" {
		return nil
	}
	fileName = "texture/" + fileName

	t := Texture{}
	if textures[fileName] != nil {
		t.tex = textures[fileName]
	} else {
		file, err := os.Open(fileName)
		if err != nil {
			log.Fatal(err)
		}
		img, _, err := image.Decode(file)
		if err != nil {
			log.Fatal(err)
		}
		t.tex = &img
		textures[fileName] = &img
	}
	bounds := (*t.tex).Bounds()
	t.Width = float64(bounds.Max.X - bounds.Min.X)
	t.Height = float64(bounds.Max.Y - bounds.Min.Y)
	return &t
}

func (t *Texture) ColorAt(coord mgl64.Vec2) (color Color64) {
	x := coord.X() * t.Width
	y := coord.Y() * t.Height
	fx := math.Floor(x)
	fy := math.Floor(y)
	i := int(fx)
	j := int(fy)
	dx := x - fx
	dy := y - fy

	img := *(t.tex)

	ulr, ulg, ulb, _ := img.At(i, j).RGBA()
	urr, urg, urb, _ := img.At(i + 1, j).RGBA()
	llr, llg, llb, _ := img.At(i, j + 1).RGBA()
	lrr, lrg, lrb, _ := img.At(i + 1, j + 1).RGBA()

	r := (1 - dx) * (1 - dy) * float64(ulr) +
		dx * (1 - dy) * float64(urr) +
		(1 - dx) * dy * float64(llr) +
		dx * dy * float64(lrr)
	g := (1 - dx) * (1 - dy) * float64(ulg) +
		dx * (1 - dy) * float64(urg) +
		(1 - dx) * dy * float64(llg) +
		dx * dy * float64(lrg)
	b := (1 - dx) * (1 - dy) * float64(ulb) +
		dx * (1 - dy) * float64(urb) +
		(1 - dx) * dy * float64(llb) +
		dx * dy * float64(lrb)

	return Color64{r / 65535, g / 65535, b / 65535}
}

type Material struct {
	Name string
	Emissive Color64
	Ambient Color64
	Specular Color64
	Reflective Color64

	Diffuse Color64
	DiffuseTextureFile string
	DiffuseTexture *Texture

	Transmissive Color64
	TransmissiveTextureFile string
	TransmissiveTexture *Texture

	Shininess float64
	Index float64

	LogTransmissive Color64
}

func InitMaterial(m *Material) {
	m.LogTransmissive = Color64{
		math.Log(2 - m.Transmissive.R()),
		math.Log(2 - m.Transmissive.G()),
		math.Log(2 - m.Transmissive.B()),
	}
	m.DiffuseTexture = NewTexture(m.DiffuseTextureFile)
	m.TransmissiveTexture = NewTexture(m.TransmissiveTextureFile)
	if mgl64.FloatEqual(m.Index, 0) {
		m.Index = AirIndex
	}
}

func (m *Material) GetDiffuseColor(isect Intersection) Color64 {
	if m.DiffuseTexture != nil {
		return m.DiffuseTexture.ColorAt(isect.UVCoords)
	}
	return m.Diffuse
}

func (m *Material) GetTransmissiveColor(isect Intersection) Color64 {
	if m.TransmissiveTexture != nil {
		return m.TransmissiveTexture.ColorAt(isect.UVCoords)
	}
	return m.Transmissive
}

func (m *Material) GetLogTransmissiveColor(isect Intersection) Color64 {
	if m.TransmissiveTexture != nil {
		t := m.TransmissiveTexture.ColorAt(isect.UVCoords)
		return Color64{
			math.Log(2 - t.R()),
			math.Log(2 - t.G()),
			math.Log(2 - t.B()),
		}
	}
	return m.LogTransmissive
}

func (m *Material) BeersTrans(isect Intersection) Color64 {
	dist := isect.T
	return Color64{
		math.Exp(m.GetLogTransmissiveColor(isect).R() * -dist),
		math.Exp(m.GetLogTransmissiveColor(isect).G() * -dist),
		math.Exp(m.GetLogTransmissiveColor(isect).B() * -dist),
	}
}

func (m *Material) ShadeBlinnPhong(scene *Scene, ray Ray, isect Intersection) (color Color64) {
	point := ray.At(isect.T)
	colorVec := mgl64.Vec3(m.Emissive).Add(mgl64.Vec3(m.Ambient.Product(scene.AmbientLight)))
	for _, light := range scene.Lights {
		attenuation := light.ShadowAttenuation(point).Mul(light.DistanceAttenuation(point))
		lightDir := light.Direction(point)
		shade := isect.Normal.Dot(lightDir)
		if shade > 0 {
			h := lightDir.Sub(ray.Direction).Normalize()
			s := mgl64.Vec3(m.Specular).Mul(math.Pow(isect.Normal.Dot(h), m.Shininess))
			d := mgl64.Vec3(m.GetDiffuseColor(isect)).Mul(shade).Add(s)
			a := Color64(attenuation).Product(Color64(d))
			contribution := mgl64.Vec3(light.GetColor().Product(a))
			colorVec = colorVec.Add(contribution)
		}
	}

	return Color64(colorVec)
}
