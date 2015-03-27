package gotracer

import (
	"github.com/go-gl/mathgl/mgl64"
)

type Intersection struct {
	Normal   mgl64.Vec3
	T        float64
	Material *Material
	UVCoords mgl64.Vec2
}
