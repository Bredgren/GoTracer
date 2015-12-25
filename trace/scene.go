package trace

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/Bredgren/gotracer/trace/bvh"
	"github.com/Bredgren/gotracer/trace/color64"
	"github.com/Bredgren/gotracer/trace/object"
	"github.com/Bredgren/gotracer/trace/options"
	"github.com/Bredgren/gotracer/trace/ray"
	"github.com/Bredgren/gotracer/trace/texture"
	"github.com/go-gl/mathgl/mgl64"
)

// Scene represents the scene and all data needed to render it.
type Scene struct {
	*options.Options
	Camera  *Camera
	BgColor color64.Color64
	BgTex   *texture.Texture

	bvh *bvh.Node
}

// ColorAt returns the color of the pixel at (x, y).
func (s *Scene) ColorAt(x, y int) color.NRGBA {
	pixelW := 1 / float64(s.Resolution.W)
	pixelH := 1 / float64(s.Resolution.H)
	centerX := float64(x) * pixelW
	centerY := float64(y) * pixelH

	var rayCounts ray.Counts
	c := s.colorAtSub(centerX, centerY, pixelW, pixelH, 0, &rayCounts)

	if !s.Options.Debug.Enabled {
		return c.NRGBA()
	}

	switch s.Options.Debug.Type {
	case "Ray Count":
		scale := math.Max(float64(s.Options.Debug.Scale), 1)
		val := math.Min(float64(rayCounts[ray.Camera])/scale, 1)
		return color64.Color64{val, val, val}.NRGBA()
	default:
		log.Fatalf("unknown debug type: %s", s.Options.Debug.Type)
		return color.NRGBA{}
	}
}

// NewScene creates and returns a new Scene from the given options.
func NewScene(options *options.Options) *Scene {
	objs, e := object.MakeObjects(options)
	if e != nil {
		log.Fatalf("creating objects: %v", e)
	}
	objects := make([]bvh.Intersector, len(objs))
	for i, obj := range objs {
		objects[i] = obj
	}

	// TODO: calculate illumination maps

	var tex *texture.Texture
	if options.Background.Type == "Uniform" || options.Global.FastRender {
		options.Background.Image = ""
	}
	tex, e = texture.New(options.Background.Image, mgl64.Vec2{}, mgl64.Vec2{1, 1})
	if e != nil {
		log.Fatalf("creating background texture: %v", e)
	}

	return &Scene{
		Options: options,
		Camera:  NewCamera(&options.Camera, float64(options.Resolution.W)/float64(options.Resolution.H)),
		BgColor: color64.Color64{options.Background.Color.R, options.Background.Color.G, options.Background.Color.B},
		BgTex:   tex,
		bvh:     bvh.NewTree(objects),
	}
}

// Result contains the results of tracing a ray into a scene. It includes the color and the number
// of each type of ray that was produced.
type Result struct {
	Color    color64.Color64
	RayCount [ray.NumTypes]int
}

func (s *Scene) colorAtSub(centerX, centerY, width, height float64, depth int, rayCounts *ray.Counts) color64.Color64 {
	if depth >= s.AntiAlias.MaxDivisions || s.Global.FastRender {
		return s.TraceDof(centerX, centerY, rayCounts)
	}
	x1 := centerX - 0.25*width
	x2 := centerX + 0.25*width
	y1 := centerY - 0.25*height
	y2 := centerY + 0.25*height
	c1 := s.TraceDof(x1, y1, rayCounts)
	c2 := s.TraceDof(x1, y2, rayCounts)
	c3 := s.TraceDof(x2, y1, rayCounts)
	c4 := s.TraceDof(x2, y2, rayCounts)
	thresh := s.AntiAlias.Threshold
	if color64.Different(c1, c2, thresh) || color64.Different(c1, c3, thresh) ||
		color64.Different(c1, c4, thresh) || color64.Different(c2, c3, thresh) ||
		color64.Different(c2, c4, thresh) || color64.Different(c3, c4, thresh) {
		d := depth + 1
		c1 = s.colorAtSub(x1, y1, width/2, height/2, d, rayCounts)
		c2 = s.colorAtSub(x1, y2, width/2, height/2, d, rayCounts)
		c3 = s.colorAtSub(x2, y1, width/2, height/2, d, rayCounts)
		c4 = s.colorAtSub(x2, y2, width/2, height/2, d, rayCounts)
	}
	// Average the 4 subpixels
	return c1.Add(c2).Add(c3).Add(c4).Mul(0.25)
}

// TraceDof handles depth of field logic if the camera is configured to use it, otherwise
// it traces a normal ray though the camera and the given normalized window coordinates.
func (s *Scene) TraceDof(nx, ny float64, rayCounts *ray.Counts) color64.Color64 {
	// Initial center ray is always cast
	var centerRay ray.Ray
	s.Camera.RayThrough(nx, ny, &centerRay)
	rayCounts[ray.Camera]++
	color := s.TraceRay(&centerRay, 0, 1.0, rayCounts)
	if !s.Camera.UseDof || s.Global.FastRender {
		return color
	}

	// Always using at least one DOF ray
	var dofRay ray.Ray
	s.Camera.DofRayThrough(&centerRay, &dofRay)
	rayCounts[ray.Camera]++
	c := s.TraceRay(&dofRay, 0, 1.0, rayCounts)
	numRays := 2
	color = color.Add(c).Mul(1 / float64(numRays))

	maxRays := s.Options.Camera.Dof.MaxRays
	thresh := s.Options.Camera.Dof.AdaptiveThreshold
	for numRays < maxRays && color64.Different(color, c, thresh) {
		rayCounts[ray.Camera]++
		c = s.TraceRay(&dofRay, 0, 1.0, rayCounts)
		numRays++
		// Rolling average
		color = color.Mul(float64(numRays) - 1).Add(c).Mul(1 / float64(numRays))
	}
	return color
}

// TraceRay sends a ray into the scene and returns the color it finds.
func (s *Scene) TraceRay(r *ray.Ray, depth int, contribution float64, rayCounts *ray.Counts) color64.Color64 {
	isect := bvh.IntersectResult{}
	s.bvh.Intersect(r, &isect)
	if isect.Object == nil {
		return s.BackgroundColor(r)
	}
	// c := 0.0
	// if isect.T <= 2 {
	// 	c = 1.0
	// } else if isect.T >= 20 {
	// 	c = 0.9
	// } else {
	// 	c = 1 - (isect.T-2)/18.0
	// }
	// return Color64{c, c, c}
	at := r.At(isect.T)
	rd := (math.Sin(0.2*at.X()) + 1) / 2
	g := (math.Sin(4.0*at.Y()) + 1) / 2
	b := (math.Sin(0.2*at.Z()) + 1) / 2
	return color64.Color64{rd, g, b}
}

// BackgroundColor returns the color a ray returns when it hits no objects.
func (s *Scene) BackgroundColor(r *ray.Ray) color64.Color64 {
	switch s.Background.Type {
	case "Skybox":
		// Assumes a cube unfolded like so:
		//          top(+Y)
		// left(-X) front(-Z) right(+X) back(+Z)
		//          bottom(-Y)
		x, y, z := r.Dir.X(), r.Dir.Y(), r.Dir.Z()
		w, h := 0.25, 1.0/3.0
		var uv mgl64.Vec2
		switch {
		case z <= -math.Abs(x) && z <= -math.Abs(y): // front
			t := -1 / z
			du, dv := ((x*t)+1)/2, 1-((y*t)+1)/2
			uv = mgl64.Vec2{w + w*du, h + h*dv}
		case x <= -math.Abs(y) && x <= -math.Abs(z): // left
			t := -1 / x
			du, dv := 1-((z*t)+1)/2, 1-((y*t)+1)/2
			uv = mgl64.Vec2{w * du, h + h*dv}
		case x >= math.Abs(y) && x >= math.Abs(z): // right
			t := 1 / x
			du, dv := ((z*t)+1)/2, 1-((y*t)+1)/2
			uv = mgl64.Vec2{2*w + w*du, h + h*dv}
		case y <= -math.Abs(x) && y <= -math.Abs(z): // bottom
			t := -1 / y
			du, dv := ((x*t)+1)/2, ((z*t)+1)/2
			uv = mgl64.Vec2{w + w*du, 2*h + h*dv}
		case y > math.Abs(x) && y > math.Abs(z): // top
			t := 1 / y
			du, dv := ((x*t)+1)/2, 1-((z*t)+1)/2
			uv = mgl64.Vec2{w + w*du, h * dv}
		case z > math.Abs(x) && z > math.Abs(y): // back
			t := 1 / z
			du, dv := 1-((x*t)+1)/2, 1-((y*t)+1)/2
			uv = mgl64.Vec2{3*w + w*du, h + h*dv}
		}
		if uv.X() == 0.0 && uv.Y() == 0.0 {
			fmt.Println(r)
		}
		return s.BgTex.ColorAt(uv)
	case "Photosphere":
		u := 0.5 + math.Atan2(r.Dir.Z(), r.Dir.X())/(2*math.Pi)
		v := 0.5 - math.Asin(r.Dir.Y())/math.Pi
		return s.BgTex.ColorAt(mgl64.Vec2{u, v})
	default:
		return s.BgColor
	}
}
