package texture

import (
	"image"
	"image/draw"
	"io"
	"log"
	"math"
	"net/http"
	"os"

	"github.com/Bredgren/gotracer/trace/color64"
	"github.com/asaskevich/govalidator"
	"github.com/go-gl/mathgl/mgl64"
)

// texture cache. All textures are loaded once and refered to by the same *Texture
var textures = make(map[string]*Texture)

// Texture represents a 2D texture image.
type Texture struct {
	img    *image.NRGBA
	offset mgl64.Vec2
	scale  mgl64.Vec2
}

// New Creates and inializes a new texture from the given source, which may be either
// a local file path or URL to an image. If the path is empty a blank texture is returned.
// If the path is invalid an error is returned.
func New(srcPath string, offset, scale mgl64.Vec2) (*Texture, error) {
	if srcPath == "" {
		return &Texture{img: image.NewNRGBA(image.Rectangle{})}, nil
	}
	if t, ok := textures[srcPath]; ok {
		return t, nil
	}
	// if srcPath is url, attempt to download image else just attempt to open image
	if govalidator.IsURL(srcPath) {
		resp, e := http.Get(srcPath)
		if e != nil {
			log.Fatalf("fetching image from url: %s: %v", srcPath, e)
		}
		defer resp.Body.Close()

		file, e := os.Create("/tmp/bgimage")
		if e != nil {
			log.Fatalf("creating temporary image: /tmp/bgimage: %v", e)
		}

		_, e = io.Copy(file, resp.Body)
		if e != nil {
			log.Fatalf("copying downloaded image to /tmp/bgimage: %v", e)
		}
		file.Close()
		srcPath = "/tmp/bgimage"
	}

	file, e := os.Open(srcPath)
	if e != nil {
		return nil, e
	}
	defer file.Close()
	img, _, e := image.Decode(file)
	if e != nil {
		return nil, e
	}

	b := img.Bounds()
	nrgba := image.NewNRGBA(b)
	draw.Draw(nrgba, b, img, b.Min, draw.Src)
	tex := &Texture{img: nrgba, offset: offset, scale: scale}
	textures[srcPath] = tex
	return tex, nil
}

// ColorAt takes normalized coordinates and returns the color linearly interpolated.
func (t *Texture) ColorAt(uv mgl64.Vec2) color64.Color64 {
	bounds := t.img.Bounds()
	w := float64(bounds.Max.X - 1)
	h := float64(bounds.Max.Y - 1)
	x, y := uv.X()*w, uv.Y()*h
	fx, fy := math.Floor(x), math.Floor(y)
	i, j := int(fx), int(fy)

	ulr, ulg, ulb, _ := t.img.At(i, j).RGBA()
	urr, urg, urb, _ := t.img.At(i+1, j).RGBA()
	llr, llg, llb, _ := t.img.At(i, j+1).RGBA()
	lrr, lrg, lrb, _ := t.img.At(i+1, j+1).RGBA()

	dx, dy := x-fx, y-fy

	ul := (1 - dx) * (1 - dy)
	ur := dx * (1 - dy)
	ll := (1 - dx) * dy
	lr := dx * dy
	r := ul*float64(ulr) + ur*float64(urr) + ll*float64(llr) + lr*float64(lrr)
	g := ul*float64(ulg) + ur*float64(urg) + ll*float64(llg) + lr*float64(lrg)
	b := ul*float64(ulb) + ur*float64(urb) + ll*float64(llb) + lr*float64(lrb)

	return color64.Color64{r / 65535, g / 65535, b / 65535}
}

// SetColorAt takes normalized coordinates and sets the color there to c, spreading it
// out reversed lineary interpolated.
func (t *Texture) SetColorAt(uv mgl64.Vec2, c color64.Color64) {
}
