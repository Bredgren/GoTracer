package raytracer

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

type SceneObject interface {
	GetTransform() mgl64.Mat4
	GetInvTransform() mgl64.Mat4
	GetMaterialName() string
	// Intersect takes a ray in local coordinates and returns an intersection and true
	// if the intersects the object, false otherwise.
	Intersect(r Ray) (isect Intersection, hit bool)
}

type SphereObject struct {
	Transform mgl64.Mat4
	MaterialName string

	invTransform mgl64.Mat4
}

func InitSphereObject(s *SphereObject) {
	s.invTransform = s.Transform.Inv()
}

func (s SphereObject) GetTransform() mgl64.Mat4 {
	return s.Transform
}

func (s SphereObject) GetInvTransform() mgl64.Mat4 {
	return s.invTransform
}

func (s SphereObject) GetMaterialName() string {
	return s.MaterialName
}

func (s SphereObject) Intersect(r Ray) (isect Intersection, hit bool) {
	isect = Intersection{Object: s}

	// -(d . o) +- sqrt((d . o)^2 - (d . d)((o . o) - 1)) / (d . d)
	do := r.Direction.Dot(r.Origin)
	dd := r.Direction.Dot(r.Direction)
	oo := r.Origin.Dot(r.Origin)

	discriminant := do * do - dd * (oo - 1)
	if discriminant < 0 {
		return isect, false
	}

	discriminant = math.Sqrt(discriminant)

	t2 := (-do + discriminant) / dd
	if t2 <= Rayε {
		return isect, false
	}

	t1 := (-do - discriminant) / dd
	if t1 > Rayε {
		isect.T = t1
		// Normalize because sphere is at origin
		isect.Normal = r.At(t1)
		InitIntersection(&isect)
		u := 0.5 + (math.Atan2(isect.Normal.Y(), isect.Normal.X()) / (2 * math.Pi))
		v := 0.5 - (math.Asin(isect.Normal.Z()) / math.Pi)
		isect.UVCoords = mgl64.Vec2{u, v}
		return isect, true
	}

	if t2 > Rayε {
		isect.T = t2
		// Normalize because sphere is at origin
		isect.Normal = r.At(t2)
		InitIntersection(&isect)
		u := 0.5 + (math.Atan2(isect.Normal.Y(), isect.Normal.X()) / (2 * math.Pi))
		v := 0.5 - (math.Asin(isect.Normal.Z()) / math.Pi)
		isect.UVCoords = mgl64.Vec2{u, v}
		return isect, true
	}

	return isect, false
}

type BoxObject struct {
	Transform mgl64.Mat4
	MaterialName string

	invTransform mgl64.Mat4
}

func InitBoxObject(b *BoxObject) {
	b.invTransform = b.Transform.Inv()
}

func (b BoxObject) GetTransform() mgl64.Mat4 {
	return b.Transform
}

func (b BoxObject) GetInvTransform() mgl64.Mat4 {
	return b.invTransform
}

func (b BoxObject) GetMaterialName() string {
	return b.MaterialName
}

func (b BoxObject) Intersect(r Ray) (isect Intersection, hit bool) {
	isect = Intersection{Object: b}

	halfSize := 0.5

	bestT := math.Inf(1)
	bestSide := -1

	for side := 0; side < 6; side++ {
		mod0Side := side % 3
		if r.Direction[mod0Side] == 0 {
			continue
		}

		t := (float64(side / 3) - halfSize - r.Origin[mod0Side]) / r.Direction[mod0Side]

		if t < Rayε || t > bestT {
			continue
		}

		mod1Side := (side + 1) % 3
		mod2Side := (side + 2) % 3
		x := r.Origin[mod1Side] + t * r.Direction[mod1Side]
		y := r.Origin[mod2Side] + t * r.Direction[mod2Side]

		if x <= halfSize && x >= -halfSize && y <= halfSize && y >= -halfSize && bestT > t {
			bestT = t
			bestSide = side
		}
	}

	if bestSide < 0 {
		return isect, false
	}

	isect.T = bestT

	// // For UV coords
	// intersectPoint := r.At(isect.T)
	// side1 := (bestSide + 1) % 3
	// side2 := (bestSide + 2) % 3

	if bestSide < 3 {
		x := 0.0
		if bestSide == 0 { x = -1.0; }
		y := 0.0
		if bestSide == 1 { y = -1.0; }
		z := 0.0
		if bestSide == 2 { z = -1.0; }
		isect.Normal = mgl64.Vec3{x, y, z}
		// TODO: Set UV coods
	} else {
		x := 0.0
		if bestSide == 3 { x = 1.0; }
		y := 0.0
		if bestSide == 4 { y = 1.0; }
		z := 0.0
		if bestSide == 5 { z = 1.0; }
		isect.Normal = mgl64.Vec3{x, y, z}
		// TODO: Set UV coods
	}

	InitIntersection(&isect)
	return isect, true
}

type SquareObject struct {
	Transform mgl64.Mat4
	MaterialName string

	invTransform mgl64.Mat4
}

func InitSquareObject(s *SquareObject) {
	s.invTransform = s.Transform.Inv()
}

func (s SquareObject) GetTransform() mgl64.Mat4 {
	return s.Transform
}

func (s SquareObject) GetInvTransform() mgl64.Mat4 {
	return s.invTransform
}

func (s SquareObject) GetMaterialName() string {
	return s.MaterialName
}

func (s SquareObject) Intersect(r Ray) (isect Intersection, hit bool) {
	isect = Intersection{Object: s}

	halfSize := 0.5

	if r.Direction.Z() == 0 {
		return isect, false
	}

	t := -r.Origin.Z() / r.Direction.Z()

	if t <= Rayε {
		return isect, false
	}

	point := r.At(t)

	if point.X() < -halfSize || point.X() > halfSize || point.Y() < -halfSize || point.Y() > halfSize {
		return isect, false
	}

	isect.T = t
	if r.Direction.Z() > 0 {
		isect.Normal = mgl64.Vec3{0, 0, -1}
	} else {
		isect.Normal = mgl64.Vec3{0, 0, 1}
	}

	isect.UVCoords = mgl64.Vec2{point.X() + 0.5, point.Y() + 0.5}

	InitIntersection(&isect)
	return isect, true
}

// type TriangleObject struct {
// 	Transform mgl64.Mat4
// 	InvTransform mgl64.Mat4
// 	MaterialName string
// 	PointA mgl64.Vec3
// 	PointB mgl64.Vec3
// 	PointC mgl64.Vec3
// }

// func (t TriangleObject) GetTransform() mgl64.Mat4 {
// 	return t.Transform
// }

// func (t TriangleObject) GetInvTransform() mgl64.Mat4 {
// 	return t.InvTransform
// }

// func (t TriangleObject) GetMaterialName() string {
// 	return t.MaterialName
// }

// func (t TriangleObject) Intersect(r Ray) (Intersection, bool) {
// 	return Intersection{}, false
// }
