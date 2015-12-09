package vec

import "github.com/go-gl/mathgl/mgl64"

type vec3Dir uint8

const (
	X vec3Dir = iota
	Y
	Z
)

// Normalize returns a new Vec3 with length 1 in the same direction as v. If v has
// length 0 then the direction of the returned vector is determined by dir.
func Normalize(v mgl64.Vec3, dir vec3Dir) mgl64.Vec3 {
	len := v.Len()
	if len == 0.0 {
		v := mgl64.Vec3{}
		v[dir] = 1.0
		return v
	}
	return mgl64.Vec3{v.X() / len, v.Y() / len, v.Z() / len}
}
