package gotracer

import (
	"image"
	_ "image/png"
	_ "image/jpeg"
	"log"
	"math"
	"os"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/Bredgren/misc"
)

// textures is a cache of the textures loaded already.
var textures map[string]*image.Image = make(map[string]*image.Image)

// Texture is a type that satisifies the MaterialAttribute interface. It enables the
// use of a png or jpeg image to specify an attribute of a material.
type Texture struct {
	tex *image.Image
	Width float64
	Height float64
}

// NewTexture creates a new Texture from the specified png/jpeg image. Textures are
// cached so calling this with the fileName more than once returns the same Texture.
func NewTexture(fileName string) *Texture {
	log.Printf("NewTexture(%s)", fileName)
	if fileName == "" {
		return nil
	}

	t := Texture{}
	if textures[fileName] != nil {
		t.tex = textures[fileName]
	} else {
		file, err := os.Open(fileName)
		misc.Check(err)
		img, _, err := image.Decode(file)
		misc.Check(err)
		t.tex = &img
		textures[fileName] = &img
	}
	bounds := (*t.tex).Bounds()
	t.Width = float64(bounds.Max.X - bounds.Min.X)
	t.Height = float64(bounds.Max.Y - bounds.Min.Y)
	return &t
}

// Performs bilinear interpolation at points between pixels. See MaterialAttribute.ColorAt.
func (t *Texture) ColorAt(coord mgl64.Vec2) (color Color64) {
	var x float64 = coord.X() * (t.Width - 1)
	var y float64 = coord.Y() * (t.Height - 1)
	var fx float64 = math.Floor(x)
	var fy float64 = math.Floor(y)
	var i int = int(fx)
	var j int = int(fy)
	var img image.Image = *(t.tex)

	ulr, ulg, ulb, _ := img.At(i, j).RGBA()
	urr, urg, urb, _ := img.At(i + 1, j).RGBA()
	llr, llg, llb, _ := img.At(i, j + 1).RGBA()
	lrr, lrg, lrb, _ := img.At(i + 1, j + 1).RGBA()

	var dx float64 = x - fx
	var dy float64 = y - fy
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
