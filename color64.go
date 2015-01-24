package raytracer

type Color64 [3]float64

func (c Color64) R() float64 {
	return c[0]
}

func (c Color64) G() float64 {
	return c[1]
}

func (c Color64) B() float64 {
	return c[2]
}
