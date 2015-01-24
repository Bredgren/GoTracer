package raytracer

import (
	"github.com/go-gl/mathgl/mgl64"
)

type Ray struct {
	Origin mgl64.Vec3
	Direction mgl64.Vec3
}

type Intersection struct {
	Object SceneObject
	Location mgl64.Vec3
	Normal mgl64.Vec3
}
