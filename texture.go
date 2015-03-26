package gotracer

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"os"

	"github.com/Bredgren/misc"
	"github.com/go-gl/mathgl/mgl64"
)

// textures is a cache of the textures loaded already.
var textures map[string]image.Image = make(map[string]image.Image)

// Texture is a type that satisifies the MaterialAttribute interface. It enables the
// use of a png or jpeg image to specify an attribute of a material.
type Texture struct {
	tex    image.Image
	Width  float64
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
		t.tex = img
		textures[fileName] = img
	}
	bounds := t.tex.Bounds()
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

	ulr, ulg, ulb, _ := t.tex.At(i, j).RGBA()
	urr, urg, urb, _ := t.tex.At(i+1, j).RGBA()
	llr, llg, llb, _ := t.tex.At(i, j+1).RGBA()
	lrr, lrg, lrb, _ := t.tex.At(i+1, j+1).RGBA()

	var dx float64 = x - fx
	var dy float64 = y - fy

	ul := (1 - dx) * (1 - dy)
	ur := dx * (1 - dy)
	ll := (1 - dx) * dy
	lr := dx * dy
	r := ul*float64(ulr) + ur*float64(urr) + ll*float64(llr) + lr*float64(lrr)
	g := ul*float64(ulg) + ur*float64(urg) + ll*float64(llg) + lr*float64(lrg)
	b := ul*float64(ulb) + ur*float64(urb) + ll*float64(llb) + lr*float64(lrb)

	return Color64{r / 65535, g / 65535, b / 65535}
}
