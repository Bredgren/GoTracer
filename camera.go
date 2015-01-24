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
	// UpDir mgl64.Vec3
	FOV float64

	u mgl64.Vec3
	v mgl64.Vec3
}

func NewCamera(imgW, imgH int, pos, lookAt mgl64.Vec3, fov float64) (c Camera) {
	viewDir := lookAt.Sub(pos)
	c = Camera{
		ImageWidth: imgW,
		ImageHeight: imgH,
		Position: pos,
		ViewDir: viewDir,
		FOV: fov,
	}
	c.Update()
	return
}

// RayThrough takes normalized window coordinates and returns the ray that goes
// through that point starting from the camera.
func (c Camera) RayThrough(nx, ny float64) Ray {
	nx -= 0.5
	ny -= 0.5
	dir := c.ViewDir.Add(c.u.Mul(nx)).Add(c.v.Mul(ny)).Normalize()
	return Ray{c.Position, dir}
}

// Update must be called after making changes to the camera.
func (c *Camera) Update() {
	fov := c.FOV / math.Pi
	normalizedHeight := 2 * math.Tan(fov / 2)
	aspectRatio := float64(c.ImageWidth) / float64(c.ImageHeight)
	c.u = /*m * */ mgl64.Vec3{1, 0, 0}.Mul(normalizedHeight * aspectRatio)
	c.v = /*m * */ mgl64.Vec3{0, 1, 0}.Mul(normalizedHeight)
	c.ViewDir = /*m * */ mgl64.Vec3{0, 0, -1}
}
