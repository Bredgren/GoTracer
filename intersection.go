package gotracer

import (
	"github.com/go-gl/mathgl/mgl64"
)

type Intersection struct {
	Object   *Object
	Normal   mgl64.Vec3
	T        float64
	UVCoords mgl64.Vec2
}
