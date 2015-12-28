package trace

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"strconv"

	"github.com/Bredgren/gotracer/trace/options"
	"github.com/Bredgren/gotracer/trace/ray"
	"github.com/Bredgren/gotracer/trace/vec"
	"github.com/go-gl/mathgl/mgl64"
)

// Camera is a viewpoint into a scene and is able to generate Camera rays.
type Camera struct {
	Position mgl64.Vec3
	ViewDir  mgl64.Vec3

	UseDof         bool
	FocalDistance  float64
	ApertureRadius float64

	u mgl64.Vec3
	v mgl64.Vec3
}

// NewCamera initializes a Camera based on the options and returns it.
func NewCamera(opts *options.Camera, aspectRatio float64) *Camera {
	fovMin, fovMax := fovMinMax(opts)
	fov := mgl64.Clamp(opts.Fov, fovMin, fovMax) * math.Pi / 180

	// Convert option vectors
	lookAt := mgl64.Vec3{opts.LookAt.X, opts.LookAt.Y, opts.LookAt.Z}
	position := mgl64.Vec3{opts.Position.X, opts.Position.Y, opts.Position.Z}
	upDir := vec.Normalize(mgl64.Vec3{opts.UpDir.X, opts.UpDir.Y, opts.UpDir.Z}, vec.Y, 1)

	// Camera space transform
	m := mgl64.LookAtV(position, lookAt, upDir).Inv()

	viewDir := mgl64.TransformNormal(mgl64.Vec3{0, 0, -1}, m).Normalize()

	// Assumes distance to camera plane is 1
	normalizedHeight := math.Abs(2 * math.Tan(fov/2))
	u := mgl64.TransformNormal(mgl64.Vec3{1, 0, 0}.Mul(normalizedHeight*aspectRatio), m)
	v := mgl64.TransformNormal(mgl64.Vec3{0, -1, 0}.Mul(normalizedHeight), m)

	return &Camera{
		Position:       position,
		ViewDir:        viewDir,
		UseDof:         opts.Dof.Enabled,
		FocalDistance:  opts.Dof.FocalDistance,
		ApertureRadius: opts.Dof.ApertureRadius,
		u:              u,
		v:              v,
	}
}

// Gets min and max for FOV from struct field.
func fovMinMax(opts *options.Camera) (min, max float64) {
	t := reflect.TypeOf(*opts)
	fovF, ok := t.FieldByName("Fov")
	if !ok {
		panic(fmt.Errorf("Field Fov not found in *cameraOpts"))
	}
	minTag := fovF.Tag.Get("min")
	min, e := strconv.ParseFloat(minTag, 64)
	if e != nil {
		panic(fmt.Errorf("Struct Tag 'min' for field Fov of *cameraOpts: %v", e))
	}
	maxTag := fovF.Tag.Get("max")
	max, e = strconv.ParseFloat(maxTag, 64)
	if e != nil {
		panic(fmt.Errorf("Struct Tag 'max' for field Fov of *cameraOpts: %v", e))
	}
	return min, max
}

// RayThrough takes normalized window coordinates and returns the ray that goes
// through that point starting from the camera.
func (c *Camera) RayThrough(nx, ny float64, r *ray.Ray) {
	r.Origin = c.Position
	r.Dir = c.ViewDir.Add(c.u.Mul(nx - 0.5)).Add(c.v.Mul(ny - 0.5)).Normalize()
}

// DofRayThrough takes a center ray (which is not midified) and modifies r to be a randomized
// ray, according do depth of field settings, whose origin is slightly off the center but
func (c *Camera) DofRayThrough(center, r *ray.Ray) {
	focalPoint := center.At(c.FocalDistance)

	offsetU := c.u.Mul(rand.Float64()*c.ApertureRadius*2 - c.ApertureRadius)
	offsetV := c.v.Mul(rand.Float64()*c.ApertureRadius*2 - c.ApertureRadius)
	offsetPosition := c.Position.Add(offsetU).Add(offsetV)

	r.Origin = offsetPosition
	r.Dir = focalPoint.Sub(offsetPosition).Normalize()
}
