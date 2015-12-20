package object

import (
	"fmt"
	"math"

	"github.com/Bredgren/gotracer/trace/bvh"
	"github.com/Bredgren/gotracer/trace/options"
	"github.com/Bredgren/gotracer/trace/ray"
	"github.com/Bredgren/gotracer/trace/vec"
	"github.com/go-gl/mathgl/mgl64"
)

type objFn func(*Object) (bvh.IntersectFn, *bvh.AABB)

var objFnMap = map[string]objFn{
	"Plane": plane,
	"Cube":  cube,
	// "Sphere": sphere,
	// "Cylinder": cylinder,
	// "Cone": cone,
	// "Triangle": triangle,
	// "Trimesh": trimesh,
	// "CSG": csg,
}

// Object reprsents an object in the scene and can be intersected with rays.
type Object struct {
	Transform    mgl64.Mat4
	InvTransform mgl64.Mat4
	MaterialName string
	IsectFn      bvh.IntersectFn
	aabb         *bvh.AABB
}

// Intersect implements the bvh.Intersector interface.
func (o *Object) Intersect(r *ray.Ray, res *bvh.IntersectResult) {
	newDir := mgl64.TransformNormal(r.Dir, o.InvTransform)
	localRay := ray.Ray{
		Origin: mgl64.TransformCoordinate(r.Origin, o.InvTransform),
		Dir:    newDir.Normalize(),
	}
	o.IsectFn(&localRay, res)
	if res.Object != nil {
		if !mgl64.FloatEqual(newDir.Len(), 0) {
			res.T /= newDir.Len()
		}
		res.Normal = mgl64.TransformNormal(res.Normal, o.Transform).Normalize()
	}
}

// AABB implements the bvh.Intersector interface.
func (o *Object) AABB() *bvh.AABB {
	return o.aabb
}

// NewObjects creates new objects from the given options including all the children.
func NewObjects(opts *options.Object) ([]*Object, error) {
	return newObjects(opts, mgl64.Ident4())
}

func newObjects(opts *options.Object, transform mgl64.Mat4) ([]*Object, error) {
	optsT := opts.Transform
	if mgl64.FloatEqual(optsT.Scale.X, 0) {
		optsT.Scale.X = 1
	}
	if mgl64.FloatEqual(optsT.Scale.Y, 0) {
		optsT.Scale.Y = 1
	}
	if mgl64.FloatEqual(optsT.Scale.Z, 0) {
		optsT.Scale.Z = 1
	}
	// transform = mgl64.Translate3D(optsT.Translate.X, optsT.Translate.Y, optsT.Translate.Z).Mul4(transform)
	// transform = mgl64.HomogRotate3D(optsT.RotateAngle*math.Pi/180, vec.Normalize(mgl64.Vec3{optsT.RotateAxis.X, optsT.RotateAxis.Y, optsT.RotateAxis.Z}, vec.Y)).Mul4(transform)
	// transform = mgl64.Scale3D(optsT.Scale.X, optsT.Scale.Y, optsT.Scale.Z).Mul4(transform)
	transform = transform.Mul4(mgl64.Translate3D(optsT.Translate.X, optsT.Translate.Y, optsT.Translate.Z))
	transform = transform.Mul4(mgl64.HomogRotate3D(optsT.RotateAngle*math.Pi/180, vec.Normalize(mgl64.Vec3{optsT.RotateAxis.X, optsT.RotateAxis.Y, optsT.RotateAxis.Z}, vec.Y)))
	transform = transform.Mul4(mgl64.Scale3D(optsT.Scale.X, optsT.Scale.Y, optsT.Scale.Z))

	var objs []*Object
	o := Object{}
	fn, ok := objFnMap[opts.Type]
	if !ok {
		return nil, fmt.Errorf("unknown object type '%s'", opts.Type)
	}
	o.Transform = transform
	o.InvTransform = transform.Inv()
	o.IsectFn, o.aabb = fn(&o)

	objs = append(objs, &o)
	for _, child := range opts.Children {
		os, e := newObjects(child, transform)
		if e != nil {
			return nil, e
		}
		objs = append(objs, os...)
	}
	return objs, nil
}

// Plane is a 2D plane object with a width and height of 1 in the XY-plane centered at the origin.
func plane(o *Object) (bvh.IntersectFn, *bvh.AABB) {
	return func(r *ray.Ray, res *bvh.IntersectResult) {
		res.Object = nil

		if mgl64.FloatEqual(r.Dir.Z(), 0) {
			return // Miss when parallel
		}

		t := -r.Origin.Z() / r.Dir.Z()

		if t < ray.Epsilon {
			return // We're too close
		}

		point := r.At(t)

		if point.X() < -0.5 || point.X() > 0.5 || point.Y() < -0.5 || point.Y() > 0.5 {
			return // Out of bounds
		}

		// Successful hit
		res.Object = o
		res.T = t
		if r.Dir.Z() > 0 {
			res.Normal = mgl64.Vec3{0, 0, -1}
		} else {
			res.Normal = mgl64.Vec3{0, 0, 1}
		}

		res.UV = mgl64.Vec2{point.X() + 0.5, 1 - (point.Y() + 0.5)}
	}, makeAABB(1, 1, 0.1, o.Transform)
}

// Cube has dimensinos 1x1x1 and is centered at the origin.
func cube(o *Object) (bvh.IntersectFn, *bvh.AABB) {
	return func(r *ray.Ray, res *bvh.IntersectResult) {
		res.Object = nil

		res.T = math.Inf(1)
		bestSide := -1

		for side := 0; side < 6; side++ {
			axis := side % 3
			if r.Dir[axis] == 0 {
				continue
			}

			t := (float64(side/3) - 0.5 - r.Origin[axis]) / r.Dir[axis]
			if t < ray.Epsilon || t > res.T {
				continue
			}

			x := r.Origin[(side+1)%3] + t*r.Dir[(side+1)%3]
			y := r.Origin[(side+2)%3] + t*r.Dir[(side+2)%3]
			if x <= 0.5 && x >= -0.5 && y <= 0.5 && y >= -0.5 && res.T > t {
				res.T = t
				bestSide = side
			}
		}

		if bestSide < 0 {
			return
		}

		res.Object = o
		res.Normal = mgl64.Vec3{}

		// Calculate UV coords and Normal
		point := r.At(res.T)
		side1 := float64((bestSide + 1) % 3)
		side2 := float64((bestSide + 2) % 3)
		if bestSide < 3 {
			res.UV = mgl64.Vec2{
				0.5 - point[int(math.Min(side1, side2))],
				0.5 + point[int(math.Max(side1, side2))],
			}
			res.Normal[bestSide%3] = -1
		} else {
			res.UV = mgl64.Vec2{
				0.5 + point[int(math.Min(side1, side2))],
				0.5 + point[int(math.Max(side1, side2))],
			}
			res.Normal[bestSide%3] = 1
		}
	}, makeAABB(1, 1, 1, o.Transform)
}

// Sphere has radius 1 centered at the origin.
type Sphere struct {
}

// Cylinder has height and radius 1 centered at the origin.
type Cylinder struct {
}

// Cone has height and base radius 1 centered at the origin.
type Cone struct {
}

// Triangle is made up of the points (0, 0, 0), (1, 0, 0), (0, 1, 0).
type Triangle struct {
}

// Trimesh is a mesh of many triangles.
type Trimesh struct {
}

// CSG (constructive solid geometry) combines other objects using union, intersection and difference.
type CSG struct {
}

func makeAABB(w, h, d float64, transform mgl64.Mat4) *bvh.AABB {
	points := aabbPoints(w, h, d)
	transformPoints(points[:], transform)
	return aabbFromPoints(points)
}

func aabbPoints(w, h, d float64) [8]mgl64.Vec3 {
	hw := w / 2
	hh := h / 2
	hd := d / 2
	return [8]mgl64.Vec3{
		{-hw, -hh, -hd}, {-hw, -hh, hd}, {hw, -hh, -hd}, {hw, -hh, hd},
		{-hw, hh, -hd}, {-hw, hh, hd}, {hw, hh, -hd}, {hw, hh, hd},
	}
}

func transformPoints(points []mgl64.Vec3, transform mgl64.Mat4) {
	for i, p := range points {
		points[i] = mgl64.TransformCoordinate(p, transform)
	}
}

func aabbFromPoints(points [8]mgl64.Vec3) *bvh.AABB {
	aabb := bvh.AABB{
		Min: mgl64.Vec3{math.Inf(1), math.Inf(1), math.Inf(1)},
		Max: mgl64.Vec3{math.Inf(-1), math.Inf(-1), math.Inf(-1)},
	}

	for _, p := range points {
		for axis := 0; axis < 3; axis++ {
			if p[axis] < aabb.Min[axis] {
				aabb.Min[axis] = p[axis]
			}
			if p[axis] > aabb.Max[axis] {
				aabb.Max[axis] = p[axis]
			}
		}
	}

	return &aabb
}
