package raytracer

type Material struct {
	Name string
	Emissive Color64
	Ambient Color64
	Specular Color64
	Reflective Color64
	Diffuse Color64
	Transmissive Color64
	Shininess float64
	Index float64
}
