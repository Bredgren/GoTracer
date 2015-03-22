package gotracer

import (
	"github.com/go-gl/mathgl/mgl64"
)

type Intersection struct {
	Object   SceneObject
	Normal   mgl64.Vec3
	T        float64
	UVCoords mgl64.Vec2
}

func InitIntersection(i *Intersection) {
	i.Normal = i.Normal.Normalize()
}
