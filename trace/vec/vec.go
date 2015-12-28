package vec

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

type vec3Dir uint8

const (
	X vec3Dir = iota
	Y
	Z
)

// Normalize returns a new Vec3 with length 1 in the same direction as v. If v has
// length 0 then the direction of the returned vector is determined by dir.
func Normalize(v mgl64.Vec3, dir vec3Dir, sign float64) mgl64.Vec3 {
	len := v.Len()
	if mgl64.FloatEqual(len, 0.0) {
		v := mgl64.Vec3{}
		v[dir] = math.Copysign(1.0, sign)
		return v
	}
	return mgl64.Vec3{v.X() / len, v.Y() / len, v.Z() / len}
}
