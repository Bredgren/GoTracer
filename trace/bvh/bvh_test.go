package bvh

import (
	"testing"

	"github.com/Bredgren/gotracer/trace/ray"
	"github.com/go-gl/mathgl/mgl64"
)

func TestAABB(t *testing.T) {
	aabb := AABB{Min: mgl64.Vec3{1, 0, -2}, Max: mgl64.Vec3{2, 1, 0}}
	cases := getTestCases(&aabb)
	res := IntersectResult{}
	for i, c := range cases {
		aabb.Intersect(&c.r, &res)
		if c.hit != res.Object {
			t.Errorf("case %d: %#v: got=%#v want=%#v", i, c.r, res.Object, c.hit)
		}
	}
}

var res IntersectResult

func benchmarkAABB(b *testing.B, aabb AABB, r ray.Ray) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		aabb.Intersect(&r, &res)
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
	hit Intersector
}

func getTestCases(aabb *AABB) []testCase {
	var cases []testCase
	for i := 0; i < 3; i++ {
		cases = append(cases, getAxisCases(aabb, i)...)
	}
	return cases
}

func getAxisCases(aabb *AABB, axis int) []testCase {
	var cases []testCase
	// Generates 54 test cases per axis, half for one direction and half for the other.
	// The 27 per direction are each position around the AABB in a 3x3x3 grid.
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
				hit1 := Intersector(nil)
				if (axis == 0 && x <= 0 && y == 0 && z == 0) ||
					(axis == 1 && x == 0 && y <= 0 && z == 0) ||
					(axis == 2 && x == 0 && y == 0 && z <= 0) {
					hit1 = aabb
				}
				hit2 := Intersector(nil)
				if (axis == 0 && x >= 0 && y == 0 && z == 0) ||
					(axis == 1 && x == 0 && y >= 0 && z == 0) ||
					(axis == 2 && x == 0 && y == 0 && z >= 0) {
					hit2 = aabb
				}
				cases = append(cases,
					testCase{
						r:   ray.Ray{Origin: origin, Dir: dir1},
						hit: hit1,
					},
					testCase{
						r:   ray.Ray{Origin: origin, Dir: dir2},
						hit: hit2,
					})
			}
		}
	}

	return cases
}

func TestBvhConstruction(t *testing.T) {
	var objects []Intersector
	tree := NewTree(objects)
	if tree != nil {
		t.Errorf("Empty object slice didn't not create nil tree, got this: %#v", tree)
	}

	o1 := AABB{Min: mgl64.Vec3{0, 0, 0}, Max: mgl64.Vec3{1, 1, 1}}
	objects = append(objects, &o1)
	tree = NewTree(objects)
	if tree.Left != nil || tree.Right != nil || tree.Object != &o1 {
		t.Errorf("Tree with 1 object: incorrect: %#v", tree)
	}

	o2 := AABB{Min: mgl64.Vec3{-2, 0, 0}, Max: mgl64.Vec3{-1, 1, 1}}
	objects = append(objects, &o2)
	tree = NewTree(objects)
	if tree.Left == nil {
		t.Errorf("Tree with 2 objects: Left is nil: %#v", tree)
	}
	if tree.Left.Object != &o2 {
		t.Errorf("Tree with 2 objects: Left is not o2: %#v", tree)
	}
	if tree.Right == nil {
		t.Errorf("Tree with 2 objects: Right is nil: %#v", tree)
	}
	if tree.Right.Object != &o1 {
		t.Errorf("Tree with 2 objects: Left is not o1: %#v", tree)
	}
	aabb := tree.AABB()
	if !aabb.Min.ApproxEqual(o2.Min) && !aabb.Max.ApproxEqual(o1.Max) {
		t.Errorf("Tree with 2 objects: AABB is incorrect: %#v", tree)
	}

	o3 := AABB{Min: mgl64.Vec3{0, 0, -1}, Max: mgl64.Vec3{1, 1, 0}}
	objects = append(objects, &o3)
	tree = NewTree(objects)
	if tree.Left == nil {
		t.Errorf("Tree with 3 objects: Left is nil: %#v", tree)
	}
	if tree.Left.Object != &o2 {
		t.Errorf("Tree with 3 objects: Left is not o2: %#v", tree)
	}
	if tree.Right == nil {
		t.Errorf("Tree with 3 objects: Right is nil: %#v", tree)
	}
	if tree.Right.Right.Object != &o1 {
		t.Errorf("Tree with 3 objects: Right.Right is not o1: %#v", tree)
	}
	if tree.Right.Left.Object != &o3 {
		t.Errorf("Tree with 3 objects: Right.Left is not o3: %#v", tree)
	}
	aabb = tree.AABB()
	if !aabb.Min.ApproxEqual(mgl64.Vec3{-2, 0, -1}) && !aabb.Max.ApproxEqual(o1.Max) {
		t.Errorf("Tree with 3 objects: AABB is incorrect: %#v", tree)
	}
}

func TestBvhIntersect(t *testing.T) {
	var objects []Intersector
	tree := NewTree(objects)
	r1 := ray.Ray{Origin: mgl64.Vec3{0, 0, 5}, Dir: mgl64.Vec3{0, 0, -1}}
	res := IntersectResult{}
	tree.Intersect(&r1, &res)
	if res.Object != nil {
		t.Errorf("Got hit on empty tree. Object is %#v", res.Object)
	}

	o1 := AABB{Min: mgl64.Vec3{-1, -1, -1}, Max: mgl64.Vec3{1, 1, 1}}
	objects = append(objects, &o1)
	tree = NewTree(objects)

	tree.Intersect(&r1, &res)
	if res.Object != &o1 {
		t.Errorf("1 Object: r1 (%#v) didn't hit o1: %#v, tree: %#v", res.Object, r1, tree)
	}

	r2 := ray.Ray{Origin: mgl64.Vec3{-2, 0, 5}, Dir: mgl64.Vec3{0, 0, -1}}
	tree.Intersect(&r2, &res)
	if res.Object == &o1 {
		t.Errorf("1 Object: r2 (%#v) hit o1: %#v tree: %#v", r2, res.Object, tree)
	}

	o2 := AABB{Min: mgl64.Vec3{3, -1, -1}, Max: mgl64.Vec3{5, 1, 1}}
	objects = append(objects, &o2)
	tree = NewTree(objects)

	tree.Intersect(&r1, &res)
	if res.Object != &o1 {
		t.Errorf("2 Objects: r1 (%#v) didn't hit o1: %#v tree: %#v", r1, res.Object, tree)
	}

	tree.Intersect(&r2, &res)
	if res.Object != nil {
		t.Errorf("2 Objects: r2 (%#v) hit something: %#v tree: %#v", r2, res.Object, tree)
	}

	r3 := ray.Ray{Origin: mgl64.Vec3{4, 0, 5}, Dir: mgl64.Vec3{0, 0, -1}}
	tree.Intersect(&r3, &res)
	if res.Object != &o2 {
		t.Errorf("2 Objects: r3 (%#v) didn't hit o2: hit %#v tree: %#v", r3, res.Object, tree)
	}

	o3 := AABB{Min: mgl64.Vec3{3, -1, 2}, Max: mgl64.Vec3{5, 1, 4}}
	objects = append(objects, &o3)
	tree = NewTree(objects)

	tree.Intersect(&r1, &res)
	if res.Object != &o1 {
		t.Errorf("3 Objects: r1 (%#v) didn't hit o1: %#v tree: %#v", r1, res.Object, tree)
	}

	tree.Intersect(&r2, &res)
	if res.Object != nil {
		t.Errorf("3 Objects: r2 (%#v) hit something: %#v tree: %#v", r2, res.Object, tree)
	}

	tree.Intersect(&r3, &res)
	if res.Object != &o3 {
		t.Errorf("3 Objects: r3 (%#v) didn't hit o3: hit %#v tree: %#v", r3, res.Object, tree)
	}
}
