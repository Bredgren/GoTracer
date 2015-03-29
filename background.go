package gotracer

type Background interface {
	GetColor(ray *Ray) Color64
}

type UniformBackground struct {
	Color Color64
}

func (b UniformBackground) GetColor(ray *Ray) Color64 {
	return b.Color
}
