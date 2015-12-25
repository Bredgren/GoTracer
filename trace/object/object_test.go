package object

import (
	"math"
	"testing"

	"github.com/Bredgren/gotracer/trace/bvh"
	"github.com/Bredgren/gotracer/trace/ray"
	"github.com/go-gl/mathgl/mgl64"
)

func TestPlane(t *testing.T) {
	o := Object{Transform: mgl64.Ident4(), InvTransform: mgl64.Ident4().Inv()}
	f, _ := plane(&o)
	res := bvh.IntersectResult{}

	cases := []struct {
		r  ray.Ray
		o  *Object
		uv mgl64.Vec2
		t  float64
	}{
		{
			r:  ray.Ray{Origin: mgl64.Vec3{0, 0, 5}, Dir: mgl64.Vec3{0, 0, -1}},
			o:  &o,
			uv: mgl64.Vec2{0.5, 0.5},
			t:  5.0,
		},
		{
			r: ray.Ray{Origin: mgl64.Vec3{1, 0, 5}, Dir: mgl64.Vec3{0, 0, -1}},
		},
		{
			r:  ray.Ray{Origin: mgl64.Vec3{0, 0, -5}, Dir: mgl64.Vec3{0, 0, 1}},
			o:  &o,
			uv: mgl64.Vec2{0.5, 0.5},
			t:  5.0,
		},
		{
			r: ray.Ray{Origin: mgl64.Vec3{1, 0, 0}, Dir: mgl64.Vec3{-1, 0, 0}},
		},
	}

	for _, c := range cases {
		f(&c.r, &res)
		if c.o == nil && res.Object != nil || c.o != nil && res.Object == nil {
			t.Errorf("Ray: %#v: expected object=%#v got %#v", c.r, c.o, res.Object)
		}
		if c.o != nil && res.UV != c.uv {
			t.Errorf("Ray: %#v: expected uv=%#v got %#v", c.r, c.uv, res.UV)
		}
		if c.o != nil && res.T != c.t {
			t.Errorf("Ray: %#v: expected t=%#v got %#v", c.r, c.t, res.T)
		}
	}
}

func BenchmarkPlane(b *testing.B) {
	o := Object{Transform: mgl64.Ident4(), InvTransform: mgl64.Ident4().Inv()}
	f, _ := plane(&o)
	res := bvh.IntersectResult{}
	r := ray.Ray{Origin: mgl64.Vec3{0, 0, 5}, Dir: mgl64.Vec3{0, 0, -1}}
	for i := 0; i <= b.N; i++ {
		f(&r, &res)
	}
}

func TestCube(t *testing.T) {
	o := Object{Transform: mgl64.Ident4(), InvTransform: mgl64.Ident4().Inv()}
	f, _ := cube(&o)
	res := bvh.IntersectResult{}

	cases := []struct {
		r  ray.Ray
		o  *Object
		uv mgl64.Vec2
		t  float64
	}{
		{
			r:  ray.Ray{Origin: mgl64.Vec3{0, 0, 5}, Dir: mgl64.Vec3{0, 0, -1}},
			o:  &o,
			uv: mgl64.Vec2{0.5, 0.5},
			t:  4.5,
		},
		{
			r: ray.Ray{Origin: mgl64.Vec3{1, 0, 5}, Dir: mgl64.Vec3{0, 0, -1}},
		},
		{
			r:  ray.Ray{Origin: mgl64.Vec3{5, 0, 0}, Dir: mgl64.Vec3{-1, 0, 0}},
			o:  &o,
			uv: mgl64.Vec2{0.5, 0.5},
			t:  4.5,
		},
	}

	for _, c := range cases {
		f(&c.r, &res)
		if c.o == nil && res.Object != nil || c.o != nil && res.Object == nil {
			t.Errorf("Ray: %#v: expected object=%#v got %#v", c.r, c.o, res.Object)
		}
		if c.o != nil && res.UV != c.uv {
			t.Errorf("Ray: %#v: expected uv=%#v got %#v", c.r, c.uv, res.UV)
		}
		if c.o != nil && res.T != c.t {
			t.Errorf("Ray: %#v: expected t=%#v got %#v", c.r, c.t, res.T)
		}
	}
}

func BenchmarkCube(b *testing.B) {
	o := Object{Transform: mgl64.Ident4(), InvTransform: mgl64.Ident4().Inv()}
	f, _ := cube(&o)
	res := bvh.IntersectResult{}
	r := ray.Ray{Origin: mgl64.Vec3{0, 0, 5}, Dir: mgl64.Vec3{0, 0, -1}}
	for i := 0; i <= b.N; i++ {
		f(&r, &res)
	}
}

func TestSphere(t *testing.T) {
	o := Object{Transform: mgl64.Ident4(), InvTransform: mgl64.Ident4().Inv()}
	f, _ := sphere(&o)
	res := bvh.IntersectResult{}

	cases := []struct {
		r  ray.Ray
		o  *Object
		uv mgl64.Vec2
		t  float64
	}{
		{
			r:  ray.Ray{Origin: mgl64.Vec3{0, 0, 5}, Dir: mgl64.Vec3{0, 0, -1}},
			o:  &o,
			uv: mgl64.Vec2{0.5, 0.0},
			t:  4,
		},
		{
			r: ray.Ray{Origin: mgl64.Vec3{2, 0, 5}, Dir: mgl64.Vec3{0, 0, -1}},
		},
		{
			r:  ray.Ray{Origin: mgl64.Vec3{5, 0, 0}, Dir: mgl64.Vec3{-1, 0, 0}},
			o:  &o,
			uv: mgl64.Vec2{0.5, 0.5},
			t:  4,
		},
	}

	for _, c := range cases {
		f(&c.r, &res)
		if c.o == nil && res.Object != nil || c.o != nil && res.Object == nil {
			t.Errorf("Ray: %#v: expected object=%#v got %#v", c.r, c.o, res.Object)
		}
		if c.o != nil && res.UV != c.uv {
			t.Errorf("Ray: %#v: expected uv=%#v got %#v", c.r, c.uv, res.UV)
		}
		if c.o != nil && res.T != c.t {
			t.Errorf("Ray: %#v: expected t=%#v got %#v", c.r, c.t, res.T)
		}
	}
}

func BenchmarkSphere(b *testing.B) {
	o := Object{Transform: mgl64.Ident4(), InvTransform: mgl64.Ident4().Inv()}
	f, _ := sphere(&o)
	res := bvh.IntersectResult{}
	r := ray.Ray{Origin: mgl64.Vec3{0, 0, 5}, Dir: mgl64.Vec3{0, 0, -1}}
	for i := 0; i <= b.N; i++ {
		f(&r, &res)
	}
}

func TestMakeAABB(t *testing.T) {
	halfHypot := math.Hypot(2, 2) / 2
	cases := []struct {
		w, h, d   float64
		transform mgl64.Mat4
		want      bvh.AABB
	}{
		{
			w: 1, h: 1, d: 1,
			transform: mgl64.Translate3D(1, 0, 0).Mul4(mgl64.HomogRotate3D(math.Pi/4, mgl64.Vec3{0, 1, 0})).Mul4(mgl64.Scale3D(2, 2, 2)),
			want: bvh.AABB{
				Min: mgl64.Vec3{1 - halfHypot, -1, -halfHypot},
				Max: mgl64.Vec3{1 + halfHypot, 1, halfHypot},
			},
		},
	}
	for _, c := range cases {
		got := makeAABB(c.w, c.h, c.d, c.transform)
		if !got.Min.ApproxEqual(c.want.Min) || !got.Max.ApproxEqual(c.want.Max) {
			t.Errorf("MakeAABB: want=%#v got=%#v", c.want, got)
		}
	}
}
