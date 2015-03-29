package gotracer

import (
	"github.com/go-gl/mathgl/mgl64"
)

type SkyboxBackground struct {
	tex Texture
}

func (b SkyboxBackground) GetColor(ray *Ray) Color64 {
	return b.tex.ColorAt(mgl64.Vec2{0, 0})
}
