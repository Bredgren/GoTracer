package trace

import "github.com/Bredgren/gotracer/trace/ray"

// Camera is a viewpoint into a scene and is able to generate Camera rays.
type Camera struct {
}

func NewCamera(opts cameraOpts) *Camera {
	return &Camera{}
}

// RayThrough takes normalized window coordinates and returns the ray that goes
// through that point starting from the camera.
func (c *Camera) RayThrough(nx, ny float64) ray.Ray {
	return ray.Ray{Type: ray.Camera}
}
