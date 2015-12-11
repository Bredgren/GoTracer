package ray

import "github.com/go-gl/mathgl/mgl64"

const (
	ε = 0.00001
)

// Type identifies the type of ray.
type Type int

const (
	// Unspecified is a ray with an unspecified ty.pe
	Unspecified Type = iota
	// Camera is a ray whose origin is the camera.
	Camera
	// Collision is a ray used to detect collision.
	Collision
	// Reflection is a ray generated from a reflection.
	Reflection
	// Refraction is a ray generated from a refraction.
	Refraction
	// Shadow is a ray toward a light source.
	Shadow
	// Illumination is a ray whose origin is a light source.
	Illumination
	// NumTypes is the number of ray types. Useful for using an array to map types to some value instead
	// of a map.
	NumTypes
)

// Counts holds counts of each ray type
type Counts [NumTypes]int

// Ray is a 3D ray with an origin and a direction.
type Ray struct {
	Origin mgl64.Vec3
	Dir    mgl64.Vec3
}

// At returns the point marked by the ray at t.
func (r *Ray) At(t float64) mgl64.Vec3 {
	return r.Origin.Add(r.Dir.Mul(t))
}
