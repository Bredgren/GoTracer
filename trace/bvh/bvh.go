package bvh

import (
	"github.com/Bredgren/gotracer/trace/ray"
	"github.com/go-gl/mathgl/mgl64"
)

// Node is a node in the BVH tree. They have up to 2 child nodes and an optional Object. Nodes
// without children are leaves. Non-leaf Nodes are not expected to contain Objcts.
type Node struct {
	Right Intersector
	Left  Intersector
	AABB
}

type Intersector interface {
	Intersect(*ray.Ray) bool
}

// Bounder is an object that provides AABB.
type Bounder interface {
	Bounds() AABB
}

// AABB is an axis-aligned bounding box.
type AABB struct {
	Min mgl64.Vec3
	Max mgl64.Vec3
}

func (a *AABB) Intersect(r *ray.Ray) bool {
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
		return false
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
		return false
	}

	if tzmin > txmin {
		txmin = tzmin
	}
	if tzmax < txmax {
		txmax = tzmax
	}

	return txmin > 0 || txmax > 0
}

// NewTree takes a list of Objects and returns the root of a BVH tree.
func NewTree([]Bounder) *Node {
	// Note that both children must be non-nil for non-leaf nodes?
	// Sort objects by position and recursively devide in half
	return nil
}

// Intersect takes a ray which be used to find the object that needs to be tested. It does
// not test intersection with the object itself.
func (n *Node) Intersect(r *ray.Ray) (hit bool) {
	// Test self
	if !n.AABB.Intersect(r) {
		return false
	}

	// If ray is going right (+X) then start with left group
	if r.Dir.X() >= 0 {
		hit = n.Left.Intersect(r)
		if !hit {
			hit = n.Right.Intersect(r)
		}
	} else {
		hit = n.Right.Intersect(r)
		if !hit {
			hit = n.Left.Intersect(r)
		}
	}

	return
}

func (n *Node) intersect(r *ray.Ray) int {
	// Return distance r goes to hit n, or -1 if no intersection
	// TODO
	return -1
}
