package raytracer

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

type Camera struct {
	ImageWidth int
	ImageHeight int
	Position mgl64.Vec3
	LookAt mgl64.Vec3
	ViewDir mgl64.Vec3
	UpDir mgl64.Vec3
	FOV float64
	Background Color64

	u mgl64.Vec3
	v mgl64.Vec3
	m mgl64.Mat3
}

func NewCamera(c Camera) (camera Camera) {
	camera = c
	if camera.FOV == 0 {
		camera.FOV = 53
	}
	camera.ViewDir = camera.LookAt.Sub(camera.Position).Normalize()

	z := camera.LookAt.Sub(camera.Position).Mul(-1).Normalize()
	if camera.UpDir.Len() == 0.0 {
		camera.UpDir = mgl64.Vec3{0, 1, 0}
	}

	y := camera.UpDir
	x := y.Cross(z).Normalize()
	y = z.Cross(x).Normalize()
	camera.m = mgl64.Mat3FromCols(x, y, z)

	fov := mgl64.DegToRad(camera.FOV)
	normalizedHeight := math.Abs(2 * math.Tan(fov / 2))
	aspectRatio := float64(camera.ImageWidth) / float64(camera.ImageHeight)
	camera.u = camera.m.Mul3x1(mgl64.Vec3{1, 0, 0}.Mul(normalizedHeight * aspectRatio))
	camera.v = camera.m.Mul3x1(mgl64.Vec3{0, -1, 0}.Mul(normalizedHeight))
	camera.ViewDir = camera.m.Mul3x1(mgl64.Vec3{0, 0, -1})

	return camera
}

// RayThrough takes normalized window coordinates and returns the ray that goes
// through that point starting from the camera.
func (c Camera) RayThrough(nx, ny float64) Ray {
	dir := c.ViewDir.Add(c.u.Mul(nx - 0.5)).Add(c.v.Mul(ny - 0.5))
	return NewRay(PrimaryRay, c.Position, dir)
}
