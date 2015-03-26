package gotracer

import (
	"github.com/go-gl/mathgl/mgl64"
)

const (
	DirectionalLightDist = 1e10
)

type DirectionalLight struct {
	Color          Color64
	OrientationInv mgl64.Vec3
}

func (l DirectionalLight) Attenuation(scene *Scene, point mgl64.Vec3) Color64 {
	return ShadowAttenuation(scene, l.OrientationInv.Mul(DirectionalLightDist), point)
}

func (l DirectionalLight) Direction(from mgl64.Vec3) mgl64.Vec3 {
	return l.OrientationInv
}
