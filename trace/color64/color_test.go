package color64

import "testing"

var dif bool

func BenchmarkDifferent(b *testing.B) {
	c1 := Color64{0.2, 0.2, 0.2}
	c2 := Color64{0.3, 0.3, 0.3}
	for i := 0; i < b.N; i++ {
		dif = Different(c1, c2, 0.01)
	}
}

func BenchmarkDifferent1(b *testing.B) {
	c1 := Color64{0.2, 0.2, 0.2}
	c2 := Color64{0.3, 0.3, 0.3}
	for i := 0; i < b.N; i++ {
		dif = different1(c1, c2, 0.01)
	}
}

func BenchmarkDifferent2(b *testing.B) {
	c1 := Color64{0.2, 0.2, 0.2}
	c2 := Color64{0.3, 0.3, 0.3}
	for i := 0; i < b.N; i++ {
		dif = different2(c1, c2, 0.01)
	}
}

func BenchmarkDifferent3(b *testing.B) {
	c1 := Color64{0.2, 0.2, 0.2}
	c2 := Color64{0.3, 0.3, 0.3}
	for i := 0; i < b.N; i++ {
		dif = different3(c1, c2, 0.01)
	}
}

func BenchmarkDifferent4(b *testing.B) {
	c1 := Color64{0.2, 0.2, 0.2}
	c2 := Color64{0.3, 0.3, 0.3}
	for i := 0; i < b.N; i++ {
		dif = different4(c1, c2, 0.01)
	}
}
