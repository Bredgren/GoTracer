package bvh

import (
	"math"
	"sort"

	"github.com/Bredgren/gotracer/trace/ray"
	"github.com/go-gl/mathgl/mgl64"
)

// Node is a node in the BVH tree. They have up to 2 child nodes and an optional Object. Nodes
// without children are leaves. Non-leaf Nodes are not expected to contain Objcts.
type Node struct {
	Right     *Node
	Left      *Node
	Object    Intersector
	splitAxis int
}

// Intersector is an object that can be tested for intersection with a ray.
type Intersector interface {
	// Intersect must make no assumtions about the initial state of IntersectionResult and
	// must be sure to set all fields appropriately.
	Intersect(*ray.Ray, *IntersectResult)
	AABB() *AABB
}

// IntersectFn implements intersection with a ray.
type IntersectFn func(*ray.Ray, *IntersectResult)

// AABBFn returns an AABB.
type AABBFn func() *AABB

// IntersectResult holds the results of an intersection.
type IntersectResult struct {
	Object Intersector
	T      float64
	Normal mgl64.Vec3
	UV     mgl64.Vec2
}

// AABB is an axis-aligned bounding box.
type AABB struct {
	Min mgl64.Vec3
	Max mgl64.Vec3
}

// Intersect implements the Intersect method of the Intersector interface..
func (a *AABB) Intersect(r *ray.Ray, res *IntersectResult) {
	res.Object = nil
	var txmin, txmax, tymin, tymax, tzmin, tzmax float64

	invX := 1 / r.Dir.X()
	if invX >= 0 {
		txmin = (a.Min.X() - r.Origin.X()) * invX
		txmax = (a.Max.X() - r.Origin.X()) * invX
	} else {
		txmin = (a.Max.X() - r.Origin.X()) * invX
		txmax = (a.Min.X() - r.Origin.X()) * invX
	}

	invY := 1 / r.Dir.Y()
	if invY >= 0 {
		tymin = (a.Min.Y() - r.Origin.Y()) * invY
		tymax = (a.Max.Y() - r.Origin.Y()) * invY
	} else {
		tymin = (a.Max.Y() - r.Origin.Y()) * invY
		tymax = (a.Min.Y() - r.Origin.Y()) * invY
	}

	if (txmin > tymax) || (tymin > txmax) {
		return
	}

	if tymin > txmin {
		txmin = tymin
	}
	if tymax < txmax {
		txmax = tymax
	}

	invZ := 1 / r.Dir.Z()
	if invZ >= 0 {
		tzmin = (a.Min.Z() - r.Origin.Z()) * invZ
		tzmax = (a.Max.Z() - r.Origin.Z()) * invZ
	} else {
		tzmin = (a.Max.Z() - r.Origin.Z()) * invZ
		tzmax = (a.Min.Z() - r.Origin.Z()) * invZ
	}

	if (txmin > tzmax) || (tzmin > txmax) {
		return
	}

	if tzmin > txmin {
		txmin = tzmin
	}
	if tzmax < txmax {
		txmax = tzmax
	}

	if txmin > 0 || txmax > 0 {
		res.Object = a
	}
	return
}

// AABB implements the AABB method for the Intersector interface.
func (a *AABB) AABB() *AABB {
	return a
}

// Center retuns the center of the AABB
func (a *AABB) Center() mgl64.Vec3 {
	return a.Min.Add(a.Max).Mul(0.5)
}

type intersectorByX []Intersector

func (in intersectorByX) Len() int {
	return len(in)
}

func (in intersectorByX) Less(i, j int) bool {
	centerI := in[i].AABB().Center()
	centerJ := in[j].AABB().Center()
	return centerI.X() < centerJ.X()
}

func (in intersectorByX) Swap(i, j int) {
	in[i], in[j] = in[j], in[i]
}

type intersectorByY []Intersector

func (in intersectorByY) Len() int {
	return len(in)
}

func (in intersectorByY) Less(i, j int) bool {
	centerI := in[i].AABB().Center()
	centerJ := in[j].AABB().Center()
	return centerI.Y() < centerJ.Y()
}

func (in intersectorByY) Swap(i, j int) {
	in[i], in[j] = in[j], in[i]
}

type intersectorByZ []Intersector

func (in intersectorByZ) Len() int {
	return len(in)
}

func (in intersectorByZ) Less(i, j int) bool {
	centerI := in[i].AABB().Center()
	centerJ := in[j].AABB().Center()
	return centerI.Z() < centerJ.Z()
}

func (in intersectorByZ) Swap(i, j int) {
	in[i], in[j] = in[j], in[i]
}

// NewTree takes a list of Intersectors and returns the root of a BVH tree. An empty slice
// returns nil.
func NewTree(objects []Intersector) *Node {
	if len(objects) == 0 {
		return nil
	}

	if len(objects) == 1 {
		return &Node{
			Left:   nil,
			Right:  nil,
			Object: objects[0],
		}
	}

	// Get AABB around all objects
	aabb := AABB{
		Min: mgl64.Vec3{math.Inf(1), math.Inf(1), math.Inf(1)},
		Max: mgl64.Vec3{math.Inf(-1), math.Inf(-1), math.Inf(-1)},
	}
	for _, o := range objects {
		ab := o.AABB()
		for i := 0; i < 3; i++ {
			if ab.Min[i] < aabb.Min[i] {
				aabb.Min[i] = ab.Min[i]
			}
			if ab.Max[i] > aabb.Max[i] {
				aabb.Max[i] = ab.Max[i]
			}
		}
	}

	// Sort according to which dimension is largest to try and make things as square as possible
	width := aabb.Max.X() - aabb.Min.X()
	height := aabb.Max.Y() - aabb.Min.Y()
	depth := aabb.Max.Z() - aabb.Min.Z()
	axis := 0
	if height > width && height > depth {
		sort.Sort(intersectorByY(objects))
		axis = 1
	} else if depth > width && depth > height {
		sort.Sort(intersectorByZ(objects))
		axis = 2
	} else {
		sort.Sort(intersectorByX(objects))
	}

	return &Node{
		Left:      NewTree(objects[:int(len(objects)/2)]),
		Right:     NewTree(objects[int(len(objects)/2):]),
		Object:    &aabb,
		splitAxis: axis,
	}
}

// Intersect takes a ray which be used to find the object that needs to be tested. It does
// not test intersection with the object itself.
func (n *Node) Intersect(r *ray.Ray, res *IntersectResult) {
	if n == nil {
		res.Object = nil
		return
	}

	if n.Object.Intersect(r, res); res.Object == nil || (n.Left == nil && n.Right == nil) {
		// Either we're a node and the ray missed or we're a leaf node, in which case we don't
		// need to touch res since it already holds the result we want.
		return
	}

	// Being careful about res.Object since we're reusing it. At this point it is non-nil
	// and we are guarenteed to call one of the Intersect methods below which must reset
	// res.Object appropriately, therefore we don't acutally need to worry about it here.

	var tmpRes IntersectResult
	// If ray is going in the positive direction then start with left group because that's
	// the one we're more likely to hit first.
	if r.Dir[n.splitAxis] >= 0 {
		n.Left.Intersect(r, res)
		n.Right.Intersect(r, &tmpRes)
	} else {
		n.Right.Intersect(r, res)
		n.Left.Intersect(r, &tmpRes)
	}

	// Be sure to pick the non-nil or closest intersection
	if tmpRes.Object != nil && (res.Object == nil || tmpRes.T < res.T) {
		*res = tmpRes
	}

	return
}

// AABB implements the AABB method for the Intersector interface.
func (n *Node) AABB() *AABB {
	return n.Object.AABB()
}
