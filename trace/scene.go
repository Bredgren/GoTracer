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
	// TODO: handle debug rendering
	return s.colorAtSub(centerX, centerY, pixelW, pixelH, 0).Color.NRGBA()
}

// Result contains the results of tracing a ray into a scene. It includes the color and the number
// of each type of ray that was produced.
type Result struct {
	Color    Color64
	RayCount [ray.NumTypes]int
}

func (s *Scene) colorAtSub(centerX, centerY, width, height float64, depth int) *Result {
	if depth >= s.AntiAlias.MaxDivisions || s.Global.FastRender {
		r := s.TraceDof(centerX, centerY)
		r.RayCount[ray.Camera]++
		return r
	}
	x1 := centerX - 0.25*width
	x2 := centerX + 0.25*width
	y1 := centerY - 0.25*height
	y2 := centerY + 0.25*height
	r1 := s.TraceDof(x1, y1)
	r2 := s.TraceDof(x1, y2)
	r3 := s.TraceDof(x2, y1)
	r4 := s.TraceDof(x2, y2)
	thresh := s.AntiAlias.Threshold
	if ColorsDifferent(r1.Color, r2.Color, thresh) || ColorsDifferent(r1.Color, r3.Color, thresh) ||
		ColorsDifferent(r1.Color, r4.Color, thresh) || ColorsDifferent(r2.Color, r3.Color, thresh) ||
		ColorsDifferent(r2.Color, r4.Color, thresh) || ColorsDifferent(r3.Color, r4.Color, thresh) {
		d := depth + 1
		r1 = s.colorAtSub(x1, y1, width/2, height/2, d)
		r2 = s.colorAtSub(x1, y2, width/2, height/2, d)
		r3 = s.colorAtSub(x2, y1, width/2, height/2, d)
		r4 = s.colorAtSub(x2, y2, width/2, height/2, d)
	}
	// Average the 4 subpixels, reuse r1
	r1.Color = Color64(mgl64.Vec3(r1.Color).Add(mgl64.Vec3(r2.Color)).Add(mgl64.Vec3(r3.Color)).Add(mgl64.Vec3(r4.Color)).Mul(0.25))
	r1.RayCount[ray.Camera] += r2.RayCount[ray.Camera] + r3.RayCount[ray.Camera] + r4.RayCount[ray.Camera]
	return r1
}

// TraceDof handles depth of field logic if the camera is configured to use it, otherwise
// it traces a normal ray though the camera and the given normalized window coordinates.
func (s *Scene) TraceDof(nx, ny float64) *Result {
	var r ray.Ray
	if !s.Camera.UseDof || s.Global.FastRender {
		s.Camera.RayThrough(nx, ny, &r)
		return s.TraceRay(&r, 0, 1.0)
	}
	// TODO handle DOF case
	s.Camera.RayThrough(nx, ny, &r)
	return s.TraceRay(&r, 0, 1.0)
}

// TraceRay sends a ray into the scene and returns the color it finds.
func (s *Scene) TraceRay(ray *ray.Ray, depth int, contribution float64) *Result {
	return &Result{
		Color: Color64{rand.Float64(), rand.Float64(), rand.Float64()},
	}
}

// NewScene creates and returns a new Scene from the given options.
func NewScene(options *Options) *Scene {
	return &Scene{
		Options: options,
		Camera:  NewCamera(&options.Camera, float64(options.Resolution.W)/float64(options.Resolution.H)),
	}
}
