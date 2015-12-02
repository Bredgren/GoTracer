package trace

import (
	"image/color"
	"math/rand"

	"github.com/Bredgren/gotracer/trace/ray"
	"github.com/go-gl/mathgl/mgl64"
)

// Scene represents the scene and all data needed to render it.
type Scene struct {
	*Options
	Camera *Camera
}

// ColorAt returns the color of the pixel at (x, y).
func (s *Scene) ColorAt(x, y int) color.NRGBA {
	pixelW := 1 / float64(s.Resolution.W)
	pixelH := 1 / float64(s.Resolution.H)
	centerX := float64(x) * pixelW
	centerY := float64(y) * pixelH
	halfW := pixelW / 2
	halfH := pixelH / 2
	// TODO: handle debug rendering
	return s.colorAtSub(centerX-halfW, centerY-halfH, centerX+halfW, centerY+halfH, 0).Color.NRGBA()
}

// Result contains the results of tracing a ray into a scene. It includes the color and the number
// of each type of ray that was produced.
type Result struct {
	Color    Color64
	RayCount [ray.NumTypes]int
}

func (s *Scene) colorAtSub(xMin, yMin, xMax, yMax float64, depth int) *Result {
	// TODO: handle depth of field
	width := xMax - xMin
	height := yMax - yMin
	if depth >= s.AntiAlias.MaxDivisions || s.Global.FastRender {
		// Render center of pixel area
		r := s.TraceRay(s.Camera.RayThrough(xMin+width/2, yMin+height/2), 0, 1.0)
		r.RayCount[ray.Camera]++
		return r

	}
	x1 := xMin + 0.25*width
	x2 := xMin + 0.75*width
	y1 := yMin + 0.25*height
	y2 := yMin + 0.75*height
	r1 := s.TraceRay(s.Camera.RayThrough(x1, y1), 0, 1.0)
	r2 := s.TraceRay(s.Camera.RayThrough(x1, y2), 0, 1.0)
	r3 := s.TraceRay(s.Camera.RayThrough(x2, y1), 0, 1.0)
	r4 := s.TraceRay(s.Camera.RayThrough(x2, y2), 0, 1.0)
	thresh := s.AntiAlias.Threshold
	if ColorsDifferent(r1.Color, r2.Color, thresh) || ColorsDifferent(r1.Color, r3.Color, thresh) ||
		ColorsDifferent(r1.Color, r4.Color, thresh) || ColorsDifferent(r2.Color, r3.Color, thresh) ||
		ColorsDifferent(r2.Color, r4.Color, thresh) || ColorsDifferent(r3.Color, r4.Color, thresh) {
		halfW := width / 2
		halfH := height / 2
		d := depth + 1
		r1 = s.colorAtSub(xMin, yMin, xMin+halfW, yMin+halfH, d)
		r2 = s.colorAtSub(xMin+halfW, yMin, xMax, yMin+halfH, d)
		r3 = s.colorAtSub(xMin, yMin+halfH, xMin+halfW, yMax, d)
		r4 = s.colorAtSub(xMin+halfW, yMin+halfH, xMax, yMax, d)
	}
	// Average the 4 subpixels, reuse r1
	r1.Color = Color64(mgl64.Vec3(r1.Color).Add(mgl64.Vec3(r2.Color)).Add(mgl64.Vec3(r3.Color)).Add(mgl64.Vec3(r4.Color)).Mul(0.25))
	r1.RayCount[ray.Camera] += r2.RayCount[ray.Camera] + r3.RayCount[ray.Camera] + r4.RayCount[ray.Camera]
	return r1
}

// TraceRay sends a ray into the scene and returns the color it finds.
func (s *Scene) TraceRay(ray ray.Ray, depth int, contribution float64) *Result {
	return &Result{
		Color: Color64{rand.Float64(), rand.Float64(), rand.Float64()},
	}
}

// NewScene creates and returns a new Scene from the given options.
func NewScene(options *Options) *Scene {
	return &Scene{
		Options: options,
		Camera:  NewCamera(options.Camera),
	}
}
