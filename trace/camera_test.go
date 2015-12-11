package trace

import (
	"math/rand"
	"testing"

	"github.com/Bredgren/gotracer/trace/ray"
	"github.com/go-gl/mathgl/mgl64"
)

func TestNewCamera(t *testing.T) {
	cases := []struct {
		opts   cameraOpts
		aspect float64
		want   Camera
	}{
		{
			opts: cameraOpts{
				Position: vectorOpts{X: 0, Y: 0, Z: 1},
				LookAt:   vectorOpts{X: 0, Y: 0, Z: 0},
				UpDir:    vectorOpts{X: 0, Y: 1, Z: 0},
				Fov:      45,
			},
			aspect: 1.0,
			want: Camera{
				Position: mgl64.Vec3{0, 0, 1},
				ViewDir:  mgl64.Vec3{0, 0, -1},
				u:        mgl64.Vec3{0.8284271247461901, 0, 0},
				v:        mgl64.Vec3{0, -0.8284271247461901, 0},
			},
		},
		{
			opts: cameraOpts{
				Position: vectorOpts{X: 0, Y: 1, Z: 1},
				LookAt:   vectorOpts{X: 0, Y: 0, Z: 0},
				UpDir:    vectorOpts{X: 0, Y: 1, Z: 0},
				Fov:      45,
			},
			aspect: 2.0,
			want: Camera{
				Position: mgl64.Vec3{0, 1, 1},
				ViewDir:  mgl64.Vec3{0, -1, -1}.Normalize(),
				u:        mgl64.Vec3{0.8284271247461901 * 2, 0, 0},
				v:        mgl64.Vec3{0, -0.5857864376269049, 0.5857864376269049},
			},
		},
		{
			opts: cameraOpts{
				Position: vectorOpts{X: 1, Y: -1, Z: 0},
				LookAt:   vectorOpts{X: 0, Y: 0, Z: 0},
				UpDir:    vectorOpts{X: 0, Y: 0, Z: -1},
				Fov:      45,
			},
			aspect: 0.5,
			want: Camera{
				Position: mgl64.Vec3{1, -1, 0},
				ViewDir:  mgl64.Vec3{-1, 1, 0}.Normalize(),
				u:        mgl64.Vec3{-0.5857864376269049 * 0.5, -0.5857864376269049 * 0.5, 0},
				v:        mgl64.Vec3{0, 0, 0.8284271247461901},
			},
		},
	}

	for i, c := range cases {
		got := NewCamera(&c.opts, c.aspect)
		if !got.Position.ApproxEqual(c.want.Position) {
			t.Errorf("NewCamera: case %d: Position: got=%#v want=%#v", i, got.Position, c.want.Position)
		}
		if !got.ViewDir.ApproxEqual(c.want.ViewDir) {
			t.Errorf("NewCamera: case %d: ViewDir: got=%#v want=%#v", i, got.ViewDir, c.want.ViewDir)
		}
		if !got.u.ApproxEqual(c.want.u) {
			t.Errorf("NewCamera: case %d: u: got=%#v want=%#v", i, got.u, c.want.u)
		}
		if !got.v.ApproxEqual(c.want.v) {
			t.Errorf("NewCamera: case %d: v: got=%#v want=%#v", i, got.v, c.want.v)
		}
	}
}

func TestRayThrough(t *testing.T) {
	cases := []struct {
		opts   cameraOpts
		aspect float64
		nx, ny float64
		want   ray.Ray
	}{
		// Center
		{
			opts: cameraOpts{
				Position: vectorOpts{X: 0, Y: 0, Z: 1},
				LookAt:   vectorOpts{X: 0, Y: 0, Z: 0},
				UpDir:    vectorOpts{X: 0, Y: 0, Z: 0},
				Fov:      45,
			},
			aspect: 1.0,
			nx:     0.5,
			ny:     0.5,
			want: ray.Ray{
				Origin: mgl64.Vec3{0, 0, 1},
				Dir:    mgl64.Vec3{0, 0, -1},
			},
		},
		// Top left corner
		{
			opts: cameraOpts{
				Position: vectorOpts{X: 0, Y: 0, Z: 1},
				LookAt:   vectorOpts{X: 0, Y: 0, Z: 0},
				UpDir:    vectorOpts{X: 0, Y: 0, Z: 0},
				Fov:      45,
			},
			aspect: 1.0,
			nx:     0,
			ny:     0,
			want: ray.Ray{
				Origin: mgl64.Vec3{0, 0, 1},
				// top left corner => abs(2 * tan(fov/2)) / 2 = 0.41421356237309503
				Dir: mgl64.Vec3{-0.41421356237309503, 0.41421356237309503, 0}.Sub(mgl64.Vec3{0, 0, 1}),
			},
		},
	}

	var got ray.Ray
	for i, c := range cases {
		camera := NewCamera(&c.opts, c.aspect)
		camera.RayThrough(c.nx, c.ny, &got)
		if !got.Origin.ApproxEqual(c.want.Origin) {
			t.Errorf("RayThrough: case %d: Origin: got=%#v want=%#v", i, got.Origin, c.want.Origin)
		}
		if !got.Dir.ApproxEqual(c.want.Dir) {
			t.Errorf("RayThrough: case %d: Dir: got=%#v want=%#v", i, got.Dir, c.want.Dir)
		}
	}
}

func TestDofRayThrough(t *testing.T) {
	cases := []struct {
		opts     cameraOpts
		aspect   float64
		nx, ny   float64
		attempts int
		want     ray.Ray
	}{
		// Center
		{
			opts: cameraOpts{
				Position: vectorOpts{X: 0, Y: 0, Z: 1},
				LookAt:   vectorOpts{X: 0, Y: 0, Z: 0},
				UpDir:    vectorOpts{X: 0, Y: 0, Z: 0},
				Fov:      45,
				Dof: dofOpts{
					FocalDistance:  5,
					ApertureRadius: 0.1,
				},
			},
			aspect:   1.0,
			nx:       0.5,
			ny:       0.5,
			attempts: 1000,
			want: ray.Ray{
				Origin: mgl64.Vec3{0, 0, 1},
				Dir:    mgl64.Vec3{0, 0, -1},
			},
		},
	}

	var center ray.Ray
	var got ray.Ray
	for i, c := range cases {
		camera := NewCamera(&c.opts, c.aspect)
		camera.RayThrough(c.nx, c.ny, &center)
		focalPoint := center.At(camera.FocalDistance)
		for attempt := 0; attempt < c.attempts; attempt++ {
			camera.DofRayThrough(&center, &got)
			// Must check X and Y not actual radius since, for simplicity, the implementation chooses within a square
			originDist := got.Origin.Sub(c.want.Origin)
			if originDist.X() > camera.ApertureRadius {
				t.Errorf("DofRayThrough: case %d: Origin X distance: got=%#v want=%#v dist=%#v", i, got.Origin, c.want.Origin, originDist.X())
			}
			if originDist.Y() > camera.ApertureRadius {
				t.Errorf("DofRayThrough: case %d: Origin Y distance: got=%#v want=%#v dist=%#v", i, got.Origin, c.want.Origin, originDist.Y())
			}
			// Arbitrary threshold of 0.01
			if got.At(camera.FocalDistance).Sub(focalPoint).Len() > 0.01 {
				dist := got.At(camera.FocalDistance).Sub(focalPoint).Len()
				t.Errorf("DofRayThrough: case %d: Intersects focal point: got=%#v want=%#v dist=%#v", i, got.At(camera.FocalDistance), focalPoint, dist)
			}
		}
	}
}

func BenchmarkDofRayThrough(b *testing.B) {
	camera := NewCamera(&cameraOpts{
		Position: vectorOpts{X: 0, Y: 0, Z: 1},
		LookAt:   vectorOpts{X: 0, Y: 0, Z: 0},
		UpDir:    vectorOpts{X: 0, Y: 0, Z: 0},
		Fov:      45,
		Dof: dofOpts{
			FocalDistance:  5,
			ApertureRadius: 0.1,
		},
	}, 1.0)

	var center ray.Ray
	var r ray.Ray
	camera.RayThrough(rand.Float64(), rand.Float64(), &center)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		camera.DofRayThrough(&center, &r)
	}
}
