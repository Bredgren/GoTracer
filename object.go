package gotracer

import (
	// "math"

	// "github.com/go-gl/mathgl/mgl64"
)

type SceneObject interface {
	// GetTransform() mgl64.Mat4
	// GetInvTransform() mgl64.Mat4
	// GetMaterialName() string
	// // Intersect takes a ray in local coordinates and returns an intersection and true
	// // if the intersects the object, false otherwise.
	// Intersect(r Ray) (isect Intersection, hit bool)
}

// type SphereObject struct {
// 	Transform mgl64.Mat4
// 	MaterialName string

// 	invTransform mgl64.Mat4
// }

// func InitSphereObject(s *SphereObject) {
// 	s.invTransform = s.Transform.Inv()
// }

// func (s SphereObject) GetTransform() mgl64.Mat4 {
// 	return s.Transform
// }

// func (s SphereObject) GetInvTransform() mgl64.Mat4 {
// 	return s.invTransform
// }

// func (s SphereObject) GetMaterialName() string {
// 	return s.MaterialName
// }

// func (s SphereObject) Intersect(r Ray) (isect Intersection, hit bool) {
// 	isect = Intersection{Object: s}

// 	// -(d . o) +- sqrt((d . o)^2 - (d . d)((o . o) - 1)) / (d . d)
// 	do := r.Direction.Dot(r.Origin)
// 	dd := r.Direction.Dot(r.Direction)
// 	oo := r.Origin.Dot(r.Origin)

// 	discriminant := do * do - dd * (oo - 1)
// 	if discriminant < 0 {
// 		return isect, false
// 	}

// 	discriminant = math.Sqrt(discriminant)

// 	t2 := (-do + discriminant) / dd
// 	if t2 <= Rayε {
// 		return isect, false
// 	}

// 	t1 := (-do - discriminant) / dd
// 	if t1 > Rayε {
// 		isect.T = t1
// 		// Normalize because sphere is at origin
// 		isect.Normal = r.At(t1)
// 		InitIntersection(&isect)
// 		u := 0.5 + (math.Atan2(isect.Normal.Y(), isect.Normal.X()) / (2 * math.Pi))
// 		v := 0.5 - (math.Asin(isect.Normal.Z()) / math.Pi)
// 		isect.UVCoords = mgl64.Vec2{u, v}
// 		return isect, true
// 	}

// 	if t2 > Rayε {
// 		isect.T = t2
// 		// Normalize because sphere is at origin
// 		isect.Normal = r.At(t2)
// 		InitIntersection(&isect)
// 		u := 0.5 + (math.Atan2(isect.Normal.Y(), isect.Normal.X()) / (2 * math.Pi))
// 		v := 0.5 - (math.Asin(isect.Normal.Z()) / math.Pi)
// 		isect.UVCoords = mgl64.Vec2{u, v}
// 		return isect, true
// 	}

// 	return isect, false
// }

// type BoxObject struct {
// 	Transform mgl64.Mat4
// 	MaterialName string

// 	invTransform mgl64.Mat4
// }

// func InitBoxObject(b *BoxObject) {
// 	b.invTransform = b.Transform.Inv()
// }

// func (b BoxObject) GetTransform() mgl64.Mat4 {
// 	return b.Transform
// }

// func (b BoxObject) GetInvTransform() mgl64.Mat4 {
// 	return b.invTransform
// }

// func (b BoxObject) GetMaterialName() string {
// 	return b.MaterialName
// }

// func (b BoxObject) Intersect(r Ray) (isect Intersection, hit bool) {
// 	isect = Intersection{Object: b}

// 	halfSize := 0.5

// 	bestT := math.Inf(1)
// 	bestSide := -1

// 	for side := 0; side < 6; side++ {
// 		mod0Side := side % 3
// 		if r.Direction[mod0Side] == 0 {
// 			continue
// 		}

// 		t := (float64(side / 3) - halfSize - r.Origin[mod0Side]) / r.Direction[mod0Side]

// 		if t < Rayε || t > bestT {
// 			continue
// 		}

// 		mod1Side := (side + 1) % 3
// 		mod2Side := (side + 2) % 3
// 		x := r.Origin[mod1Side] + t * r.Direction[mod1Side]
// 		y := r.Origin[mod2Side] + t * r.Direction[mod2Side]

// 		if x <= halfSize && x >= -halfSize && y <= halfSize && y >= -halfSize && bestT > t {
// 			bestT = t
// 			bestSide = side
// 		}
// 	}

// 	if bestSide < 0 {
// 		return isect, false
// 	}

// 	isect.T = bestT

// 	// For UV coords
// 	intersectPoint := r.At(isect.T)
// 	side1 := float64((bestSide + 1) % 3)
// 	side2 := float64((bestSide + 2) % 3)

// 	if bestSide < 3 {
// 		x := 0.0
// 		if bestSide == 0 { x = -1.0; }
// 		y := 0.0
// 		if bestSide == 1 { y = -1.0; }
// 		z := 0.0
// 		if bestSide == 2 { z = -1.0; }
// 		isect.Normal = mgl64.Vec3{x, y, z}
// 		InitIntersection(&isect)
// 		isect.UVCoords = mgl64.Vec2{
// 			0.5 - intersectPoint[int(math.Min(side1, side2))],
// 			0.5 + intersectPoint[int(math.Max(side1, side2))],
// 		}
// 	} else {
// 		x := 0.0
// 		if bestSide == 3 { x = 1.0; }
// 		y := 0.0
// 		if bestSide == 4 { y = 1.0; }
// 		z := 0.0
// 		if bestSide == 5 { z = 1.0; }
// 		isect.Normal = mgl64.Vec3{x, y, z}
// 		InitIntersection(&isect)
// 		isect.UVCoords = mgl64.Vec2{
// 			0.5 + intersectPoint[int(math.Min(side1, side2))],
// 			0.5 + intersectPoint[int(math.Max(side1, side2))],
// 		}
// 	}

// 	return isect, true
// }

// type SquareObject struct {
// 	Transform mgl64.Mat4
// 	MaterialName string

// 	invTransform mgl64.Mat4
// }

// func InitSquareObject(s *SquareObject) {
// 	s.invTransform = s.Transform.Inv()
// }

// func (s SquareObject) GetTransform() mgl64.Mat4 {
// 	return s.Transform
// }

// func (s SquareObject) GetInvTransform() mgl64.Mat4 {
// 	return s.invTransform
// }

// func (s SquareObject) GetMaterialName() string {
// 	return s.MaterialName
// }

// func (s SquareObject) Intersect(r Ray) (isect Intersection, hit bool) {
// 	isect = Intersection{Object: s}

// 	halfSize := 0.5

// 	if r.Direction.Z() == 0 {
// 		return isect, false
// 	}

// 	t := -r.Origin.Z() / r.Direction.Z()

// 	if t <= Rayε {
// 		return isect, false
// 	}

// 	point := r.At(t)

// 	if point.X() < -halfSize || point.X() > halfSize || point.Y() < -halfSize || point.Y() > halfSize {
// 		return isect, false
// 	}

// 	isect.T = t
// 	if r.Direction.Z() > 0 {
// 		isect.Normal = mgl64.Vec3{0, 0, -1}
// 	} else {
// 		isect.Normal = mgl64.Vec3{0, 0, 1}
// 	}

// 	isect.UVCoords = mgl64.Vec2{point.X() + 0.5, 1 - (point.Y() + 0.5)}

// 	InitIntersection(&isect)
// 	return isect, true
// }

// type CylinderObject struct {
// 	Transform mgl64.Mat4
// 	MaterialName string
// 	Capped bool

// 	invTransform mgl64.Mat4
// }

// func InitCylinderObject(c *CylinderObject) {
// 	c.invTransform = c.Transform.Inv()
// }

// func (c CylinderObject) GetTransform() mgl64.Mat4 {
// 	return c.Transform
// }

// func (c CylinderObject) GetInvTransform() mgl64.Mat4 {
// 	return c.invTransform
// }

// func (c CylinderObject) GetMaterialName() string {
// 	return c.MaterialName
// }

// func (c CylinderObject) Intersect(r Ray) (isect Intersection, hit bool) {
// 	isect = Intersection{Object: c}

// 	if c.IntersectCaps(r, &isect) {
// 		i := Intersection{Object: c}
// 		if c.IntersectBody(r, &i) {
// 			if i.T < isect.T {
// 				isect = i
// 			}
// 		}
// 		InitIntersection(&isect)
// 		return isect, true
// 	}

// 	hit = c.IntersectBody(r, &isect)
// 	InitIntersection(&isect)
// 	return isect, hit
// }

// func (c CylinderObject) IntersectCaps(r Ray, isect *Intersection) (hit bool) {
// 	if !c.Capped {
// 		return false
// 	}

// 	pz := r.Origin.Z()
// 	dz := r.Direction.Z()

// 	if mgl64.FloatEqual(dz, 0) {
// 		return false
// 	}

// 	t1 := (1 - pz) / dz
// 	t2 := -pz / dz

// 	if dz > 0 {
// 		t1, t2 = t2, t1
// 	}

// 	if t2 < Rayε {
// 		return false
// 	}

// 	if t1 >= Rayε {
// 		p := r.At(t1)
// 		if p.X() * p.X() + p.Y() * p.Y() <= 1 {
// 			isect.T = t1
// 			if dz > 0 {
// 				isect.Normal = mgl64.Vec3{0, 0, -1}
// 			} else {
// 				isect.Normal = mgl64.Vec3{0, 0, 1}
// 			}
// 			isect.UVCoords = mgl64.Vec2{p.X() / 2 + 0.5, p.Y() / 2 + 0.5}
// 			return true
// 		}
// 	}

// 	p := r.At(t2)
// 	if p.X() * p.X() + p.Y() * p.Y() <= 1 {
// 		isect.T = t2
// 		if dz > 0 {
// 			isect.Normal = mgl64.Vec3{0, 0, 1}
// 		} else {
// 			isect.Normal = mgl64.Vec3{0, 0, -1}
// 		}
// 		isect.UVCoords = mgl64.Vec2{p.X() / 2 + 0.5, p.Y() / 2 + 0.5}
// 		return true
// 	}

// 	return false
// }

// func (cyl CylinderObject) IntersectBody(r Ray, isect *Intersection) (hit bool) {
// 	x0 := r.Origin.X()
// 	y0 := r.Origin.Y()
// 	x1 := r.Direction.X()
// 	y1 := r.Direction.Y()

// 	a := x1 * x1 + y1 * y1
// 	b := 2 * (x0 * x1 + y0 * y1)
// 	c := x0 * x0 + y0 * y0 - 1

// 	if mgl64.FloatEqual(a, 0) {
// 		// x1 = 0, y1 = 0 -> ray aligned with body
// 		return false
// 	}

// 	discriminant := b * b - 4 * a * c
// 	if discriminant < 0 {
// 		return false
// 	}

// 	discriminant = math.Sqrt(discriminant)

// 	t2 := (-b + discriminant) / (2 * a)
// 	if t2 <= Rayε {
// 		return false
// 	}

// 	t1 := (-b - discriminant) / (2 * a)
// 	if t1 > Rayε {
// 		p := r.At(t1)
// 		z := p.Z()
// 		if z >= 0 && z <= 1 {
// 			isect.T = t1
// 			isect.Normal = mgl64.Vec3{p.X(), p.Y(), 0}
// 			isect.UVCoords = mgl64.Vec2{0.5 + (math.Atan2(p.Y(), p.X()) / (2 * math.Pi)), 1 - p.Z()}
// 			return true
// 		}
// 	}

// 	p := r.At(t2)
// 	z := p.Z()
// 	if z >= 0 && z <= 1 {
// 		isect.T = t2
// 		normal := mgl64.Vec3{p.X(), p.Y(), 0}
// 		if !cyl.Capped && normal.Dot(r.Direction) > 0 {
// 			normal = normal.Mul(-1)
// 		}
// 		isect.Normal = normal
// 		isect.UVCoords = mgl64.Vec2{0.5 + (math.Atan2(p.Y(), p.X()) / (2 * math.Pi)), 1 - p.Z()}
// 		return true
// 	}

// 	return false
// }

// type ConeObject struct {
// 	Transform mgl64.Mat4
// 	MaterialName string
// 	Capped bool
// 	BaseRadius float64
// 	TopRadius float64

// 	invTransform mgl64.Mat4
// 	betaSquared float64
// 	gamma float64
// }

// func InitConeObject(c *ConeObject) {
// 	c.invTransform = c.Transform.Inv()
// 	c.BaseRadius = math.Max(c.BaseRadius, 0.0001)
// 	c.TopRadius = math.Max(c.TopRadius, 0.0001)
// 	beta := c.BaseRadius - c.TopRadius
// 	beta = math.Max(beta, 0.001)
// 	if math.Abs(beta) < 0.001 {
// 		beta = 0.001
// 	}
// 	if beta < 0 {
// 		c.gamma = c.BaseRadius / beta
// 	} else {
// 		c.gamma = c.TopRadius / beta
// 	}
// 	c.betaSquared = beta * beta
// 	if c.gamma < 0 {
// 		c.gamma = c.gamma - 1
// 	}
// }

// func (c ConeObject) GetTransform() mgl64.Mat4 {
// 	return c.Transform
// }

// func (c ConeObject) GetInvTransform() mgl64.Mat4 {
// 	return c.invTransform
// }

// func (c ConeObject) GetMaterialName() string {
// 	return c.MaterialName
// }

// func (co ConeObject) Intersect(r Ray) (isect Intersection, hit bool) {
// 	isect = Intersection{Object: co}

// 	ox := r.Origin.X()
// 	oy := r.Origin.Y()
// 	oz := r.Origin.Z()
// 	dx := r.Direction.X()
// 	dy := r.Direction.Y()
// 	dz := r.Direction.Z()

// 	a := dx * dx + dy * dy - co.betaSquared * dz * dz
// 	if mgl64.FloatEqual(a, 0) {
// 		InitIntersection(&isect)
// 		return isect, false
// 	}
// 	b := 2 * (ox * dx + oy * dy - co.betaSquared * ((co.gamma + oz) * dz))
// 	c := ox * ox + oy * oy - co.betaSquared * (co.gamma + oz) * (co.gamma + oz)

// 	discriminant := b * b - 4 * a * c
// 	if discriminant < 0 {
// 		InitIntersection(&isect)
// 		return isect, false
// 	}
// 	discriminant = math.Sqrt(discriminant)

// 	isect.T = -1

// 	nearRoot := (-b + discriminant) / (2 * a)
// 	p := r.At(nearRoot)
// 	if isGoodRoot(p) && nearRoot > Rayε {
// 		isect.T = nearRoot
// 		isect.Normal = mgl64.Vec3{2 * p.X(), 2 * p.Y(), -2 * co.betaSquared * (p.Z() + co.gamma)}
// 	}

// 	farRoot := (-b - discriminant) / (2 * a)
// 	p = r.At(farRoot)
// 	if isGoodRoot(p) && (farRoot < isect.T || isect.T < 0) && farRoot > Rayε {
// 		isect.T = farRoot
// 		isect.Normal = mgl64.Vec3{2 * p.X(), 2 * p.Y(), -2 * co.betaSquared * (p.Z() + co.gamma)}
// 	}

// 	if !co.Capped && isect.Normal.Dot(r.Direction) > 0 {
// 		isect.Normal = isect.Normal.Mul(-1)
// 	}

// 	t1 := -oz / dz
// 	t2 := (1 - oz) / dz

// 	p = r.At(t1)
// 	if co.Capped {
// 		if p.X() * p.X() + p.Y() * p.Y() <= co.TopRadius * co.TopRadius {
// 			if (t1 < isect.T || isect.T < 0) && t1 > Rayε {
// 				isect.T = t1
// 				isect.Normal = mgl64.Vec3{0, 0, -1}
// 			}
// 		}
// 		q := r.At(t2)
// 		if q.X() * q.X() + q.Y() * q.Y() <= co.BaseRadius * co.BaseRadius {
// 			if (t2 < isect.T || isect.T < 0) && t2 > Rayε {
// 				isect.T = t2
// 				isect.Normal = mgl64.Vec3{0, 0, 1}
// 			}
// 		}
// 	}

// 	if isect.T <= Rayε {
// 		InitIntersection(&isect)
// 		return isect, false
// 	}

// 	p = r.At(isect.T)
// 	rad := math.Max(co.BaseRadius, co.TopRadius)
// 	u := p.X() / (2 * rad) + 0.5
// 	v := p.Y() / (2 * rad) + 0.5
// 	isect.UVCoords = mgl64.Vec2{u, v}

// 	InitIntersection(&isect)
// 	return isect, true
// }

// func isGoodRoot(root mgl64.Vec3) bool {
// 	return root.Z() >= 0 && root.Z() <= 1
// }

// type TorusObject struct {
// 	Transform mgl64.Mat4
// 	MaterialName string

// 	invTransform mgl64.Mat4
// }

// func InitTorusObject(t *TorusObject) {
// 	t.invTransform = t.Transform.Inv()
// }

// func (t TorusObject) GetTransform() mgl64.Mat4 {
// 	return t.Transform
// }

// func (t TorusObject) GetInvTransform() mgl64.Mat4 {
// 	return t.invTransform
// }

// func (t TorusObject) GetMaterialName() string {
// 	return t.MaterialName
// }

// func (t TorusObject) Intersect(r Ray) (isect Intersection, hit bool) {
// 	isect = Intersection{Object: t}

// 	InitIntersection(&isect)
// 	return isect, false
// }

// // type TriangleObject struct {
// // 	Transform mgl64.Mat4
// // 	InvTransform mgl64.Mat4
// // 	MaterialName string
// // 	PointA mgl64.Vec3
// // 	PointB mgl64.Vec3
// // 	PointC mgl64.Vec3
// // }

// // func (t TriangleObject) GetTransform() mgl64.Mat4 {
// // 	return t.Transform
// // }

// // func (t TriangleObject) GetInvTransform() mgl64.Mat4 {
// // 	return t.InvTransform
// // }

// // func (t TriangleObject) GetMaterialName() string {
// // 	return t.MaterialName
// // }

// // func (t TriangleObject) Intersect(r Ray) (Intersection, bool) {
// // 	return Intersection{}, false
// // }
