package trace

import (
	"image/color"
	"math/rand"

	"github.com/go-gl/mathgl/mgl64"
)

// Scene represents the scene and all data needed to render it.
type Scene struct {
	Options
}

// ColorAt returns the color of the pixel at (x, y).
func (s *Scene) ColorAt(x, y int) color.NRGBA {
	pixelW := 1 / float64(s.Resolution.W)
	pixelH := 1 / float64(s.Resolution.H)
	centerX := float64(x) * pixelW
	centerY := float64(y) * pixelH
	halfW := pixelW / 2
	halfH := pixelH / 2
	return s.colorAtSub(centerX-halfW, centerY-halfH, centerX+halfW, centerY+halfH, 0).NRGBA()
}

func (s *Scene) colorAtSub(xMin, yMin, xMax, yMax float64, depth int) Color64 {
	width := xMax - xMin
	height := yMax - yMin
	if depth >= s.AntiAlias.MaxDivisions || s.Global.FastRender {
		// Render center of pixel area
		x := xMin + width/2
		y := yMin + height/2
		return s.TraceRay(x, y, 0, 1.0)
		// 		return scene.TraceRay(scene.Camera.RayThrough(x, y), 0, 1.0)

	}
	x1 := xMin + 0.25*width
	x2 := xMin + 0.75*width
	y1 := yMin + 0.25*height
	y2 := yMin + 0.75*height
	c1 := s.TraceRay(x1, y1, 0, 1.0)
	c2 := s.TraceRay(x2, y1, 0, 1.0)
	c3 := s.TraceRay(x1, y2, 0, 1.0)
	c4 := s.TraceRay(x2, y2, 0, 1.0)
	thresh := s.AntiAlias.Threshold
	if ColorsDifferent(c1, c2, thresh) || ColorsDifferent(c1, c3, thresh) ||
		ColorsDifferent(c1, c4, thresh) || ColorsDifferent(c2, c3, thresh) ||
		ColorsDifferent(c2, c4, thresh) || ColorsDifferent(c3, c4, thresh) {
		halfW := width / 2
		halfH := height / 2
		d := depth + 1
		c1 = s.colorAtSub(xMin, yMin, xMin+halfW, yMin+halfH, d)
		c2 = s.colorAtSub(xMin+halfW, yMin, xMax, yMin+halfH, d)
		c3 = s.colorAtSub(xMin, yMin+halfH, xMin+halfW, yMax, d)
		c4 = s.colorAtSub(xMin+halfW, yMin+halfH, xMax, yMax, d)
	}
	// Average the 4 subpixels
	return Color64(mgl64.Vec3(c1).Add(mgl64.Vec3(c2)).Add(mgl64.Vec3(c3)).Add(mgl64.Vec3(c4)).Mul(0.25))
}

// TraceRay sends a ray ino the scene and returns the color it finds.
func (s *Scene) TraceRay(x, y float64, depth int, contribution float64) Color64 {
	// Replace x/y with Ray or calculate here?
	// ray := s.Camera.RayThrough(x, y)
	return Color64{rand.Float64(), rand.Float64(), rand.Float64()}
}

// NewScene creates and returns a new Scene from the given options.
func NewScene(options *Options) *Scene {
	return &Scene{}
}
