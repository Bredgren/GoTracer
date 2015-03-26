package gotracer

import (
	"testing"

	"github.com/go-gl/mathgl/mgl64"
)

type testPixel struct {
	pixel mgl64.Vec2
	color Color64
}

var testPixels []testPixel = []testPixel{
	testPixel{mgl64.Vec2{0, 0}, Color64{1, 0, 0}},
	testPixel{mgl64.Vec2{1, 0}, Color64{0, 1, 0}},
	testPixel{mgl64.Vec2{0, 1}, Color64{0, 0, 1}},
	testPixel{mgl64.Vec2{1, 1}, Color64{1, 1, 1}},
	testPixel{mgl64.Vec2{0.5, 0.5}, Color64{0.5, 0.5, 0.5}},
}

func TestTexture(t *testing.T) {
	var tex *Texture = NewTexture("texture/test.png")
	if tex.Width != 2 && tex.Height != 2 {
		t.Errorf("Texture size is %vx%v, expectd 2x2", tex.Width, tex.Height)
	}
	for _, p := range testPixels {
		c := tex.ColorAt(p.pixel)
		if c != p.color {
			t.Errorf("Color at %v is %v, expected %v", p.pixel, c, p.color)
		}
	}
}


func BenchmarkColorAt(b *testing.B) {
	var tex *Texture = NewTexture("texture/test.png")

	coord := mgl64.Vec2{0, 0}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tex.ColorAt(coord)
	}
}
