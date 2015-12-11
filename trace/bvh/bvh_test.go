package bvh

import (
	"testing"

	"github.com/Bredgren/gotracer/trace/ray"
	"github.com/go-gl/mathgl/mgl64"
)

// func TestBvh(t *testing.T) {
// }

func TestAABB(t *testing.T) {
	aabb := AABB{Min: mgl64.Vec3{1, 0, -2}, Max: mgl64.Vec3{2, 1, 0}}
	cases := getTestCases(aabb)
	for i, c := range cases {
		got := aabb.Intersect(&c.r)
		if got != c.hit {
			t.Errorf("case %d: %#v: got=%#v want=%#v", i, c.r, got, c.hit)
		}
	}
}

var res bool

func benchmarkAABB(b *testing.B, aabb AABB, r ray.Ray) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res = aabb.Intersect(&r)
	}
}

func BenchmarkAABBMissXY(b *testing.B) {
	aabb := AABB{Min: mgl64.Vec3{1, 0, -2}, Max: mgl64.Vec3{2, 1, 0}}
	r := ray.Ray{Origin: mgl64.Vec3{0, 1, -1}, Dir: mgl64.Vec3{1, 1, 0}.Normalize()}
	benchmarkAABB(b, aabb, r)
}

func BenchmarkAABBMissZ(b *testing.B) {
	aabb := AABB{Min: mgl64.Vec3{1, 0, -2}, Max: mgl64.Vec3{2, 1, 0}}
	r := ray.Ray{Origin: mgl64.Vec3{0, 0, 1}, Dir: mgl64.Vec3{1, 0.5, 0}.Normalize()}
	benchmarkAABB(b, aabb, r)
}

func BenchmarkAABBHit(b *testing.B) {
	aabb := AABB{Min: mgl64.Vec3{1, 0, -2}, Max: mgl64.Vec3{2, 1, 0}}
	r := ray.Ray{Origin: mgl64.Vec3{0, 0, -1}, Dir: mgl64.Vec3{1, 0.5, 0}.Normalize()}
	benchmarkAABB(b, aabb, r)
}

type testCase struct {
	r   ray.Ray
	hit bool
}

func getTestCases(aabb AABB) []testCase {
	var cases []testCase
	for i := 0; i < 3; i++ {
		cases = append(cases, getAxisCases(aabb, i)...)
	}
	return cases
}

func getAxisCases(aabb AABB, axis int) []testCase {
	var cases []testCase
	for x := -1; x < 2; x++ {
		xCenter := (aabb.Max.X() + aabb.Min.X()) / 2
		xWidth := (aabb.Max.X() - aabb.Min.X())
		xCoord := xCenter + xWidth*float64(x)
		for y := -1; y < 2; y++ {
			yCenter := (aabb.Max.Y() + aabb.Min.Y()) / 2
			yWidth := (aabb.Max.Y() - aabb.Min.Y())
			yCoord := yCenter + yWidth*float64(y)
			for z := -1; z < 2; z++ {
				zCenter := (aabb.Max.Z() + aabb.Min.Z()) / 2
				zWidth := (aabb.Max.Z() - aabb.Min.Z())
				zCoord := zCenter + zWidth*float64(z)
				origin := mgl64.Vec3{xCoord, yCoord, zCoord}
				dir1 := mgl64.Vec3{}
				dir1[axis] = 1
				dir2 := mgl64.Vec3{}
				dir2[axis] = -1
				cases = append(cases,
					testCase{
						r: ray.Ray{Origin: origin, Dir: dir1},
						hit: (axis == 0 && x <= 0 && y == 0 && z == 0) ||
							(axis == 1 && x == 0 && y <= 0 && z == 0) ||
							(axis == 2 && x == 0 && y == 0 && z <= 0),
					},
					testCase{
						r: ray.Ray{Origin: origin, Dir: dir2},
						hit: (axis == 0 && x >= 0 && y == 0 && z == 0) ||
							(axis == 1 && x == 0 && y >= 0 && z == 0) ||
							(axis == 2 && x == 0 && y == 0 && z >= 0),
					})
			}
		}
	}

	return cases
}
