package raytracer

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

type Camera struct {
	ImageWidth int
	ImageHeight int
	Position mgl64.Vec3
	ViewDir mgl64.Vec3
	UpDir mgl64.Vec3
	FOV float64
	Background Color64

	u mgl64.Vec3
	v mgl64.Vec3
	m mgl64.Mat3
}

func NewCamera(imgW, imgH int, pos, lookAt, upDir mgl64.Vec3, fov float64,
	bg Color64) (c Camera) {
	if fov == 0 {
		fov = 53
	}
	c = Camera{
		ImageWidth: imgW,
		ImageHeight: imgH,
		Position: pos,
		ViewDir: lookAt.Sub(pos).Normalize(),
		FOV: fov,
		Background: bg,
	}

	z := lookAt.Sub(pos).Mul(-1).Normalize()
	if upDir.Len() == 0.0 {
		upDir = mgl64.Vec3{0, 1, 0}
	}
	y := upDir
	x := y.Cross(z).Normalize()
	y = z.Cross(x).Normalize()
	c.m = mgl64.Mat3FromCols(x, y, z)
	c.Update()
	return
}

// RayThrough takes normalized window coordinates and returns the ray that goes
// through that point starting from the camera.
func (c Camera) RayThrough(nx, ny float64) Ray {
	dir := c.ViewDir.Add(c.u.Mul(nx - 0.5)).Add(c.v.Mul(ny - 0.5))
	return NewRay(c.Position, dir)
}

// Update must be called after making changes to the camera.
func (c *Camera) Update() {
	fov := mgl64.DegToRad(c.FOV)
	normalizedHeight := math.Abs(2 * math.Tan(fov / 2))
	aspectRatio := float64(c.ImageWidth) / float64(c.ImageHeight)
	c.u = c.m.Mul3x1(mgl64.Vec3{1, 0, 0}.Mul(normalizedHeight * aspectRatio))
	c.v = c.m.Mul3x1(mgl64.Vec3{0, -1, 0}.Mul(normalizedHeight))
	c.ViewDir = c.m.Mul3x1(mgl64.Vec3{0, 0, -1})
}
