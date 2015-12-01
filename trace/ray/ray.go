package ray

import "github.com/go-gl/mathgl/mgl64"

const (
	ε = 0.00001
)

type Type int

const (
	Unspecified Type = iota
	Camera
	Collision
	Reflection
	Refraction
	Shadow
)

type Ray struct {
	Type   Type
	Origin mgl64.Vec3
	Dir    mgl64.Vec3
}
