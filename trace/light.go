package trace

import (
	"fmt"

	"github.com/Bredgren/gotracer/trace/color64"
	"github.com/Bredgren/gotracer/trace/options"
	"github.com/Bredgren/gotracer/trace/ray"
	"github.com/Bredgren/gotracer/trace/vec"
	"github.com/go-gl/mathgl/mgl64"
)

type Light interface {
	// Attenuation(scene *Scene, point mgl64.Vec3) Color64
	// Direction(from mgl64.Vec3) mgl64.Vec3
	GenIllumMap(*Scene)
}

type Coeff struct {
	Constant  float64
	Linear    float64
	Quadratic float64
}

func NewLight(opts *options.Light) (Light, error) {
	switch opts.Type {
	case "Spot":
		return &Spot{
			Color:     color64.Color64{opts.Color.R, opts.Color.G, opts.Color.B},
			Position:  mgl64.Vec3{opts.Position.X, opts.Position.Y, opts.Position.Z},
			Direction: vec.Normalize(mgl64.Vec3{opts.Direction.X, opts.Direction.Y, opts.Direction.Z}, vec.Y, -1),
			IllumMap:  opts.IlluminationMap,
			Angle:     opts.Angle,
			DropOff:   opts.DropOff,
			FadeAngle: opts.FadeAngle,
			Coeff: Coeff{
				Constant:  opts.Coeff.Constant,
				Linear:    opts.Coeff.Linear,
				Quadratic: opts.Coeff.Quadratic,
			},
		}, nil
	}
	return nil, fmt.Errorf("unkonwn light type: %s", opts.Type)
}

type Spot struct {
	Color     color64.Color64
	Position  mgl64.Vec3
	Direction mgl64.Vec3
	IllumMap  bool
	Coeff
	Angle     float64
	DropOff   float64
	FadeAngle float64
}

func (s *Spot) GenIllumMap(scene *Scene) {
	if !s.IllumMap {
		return
	}
	r := ray.Ray{Origin: s.Position}
	for x := -5.0; x < 5.0; x += 0.001 {
		for z := -4.0; z < 5.0; z += 0.001 {
			r.Dir = mgl64.Vec3{x, -1, z}.Normalize()
			scene.TraceIllumRay(&r, s.Color, 0)
		}
	}
}
