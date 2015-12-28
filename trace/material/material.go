package material

import (
	"fmt"

	"github.com/Bredgren/gotracer/trace/color64"
	"github.com/Bredgren/gotracer/trace/options"
	"github.com/Bredgren/gotracer/trace/texture"
	"github.com/go-gl/mathgl/mgl64"
)

// BRDF (bidirectional reflectance distribution function)
type BRDF int

const (
	// Lambert is a perfectly diffuse shader.
	Lambert BRDF = iota
	// BlinnPhong uses all the material properties to produce a more realistic shader.
	BlinnPhong
)

// Material holds all the properties of a material.
type Material struct {
	Emissive     Property
	Ambient      Property
	Diffuse      Property
	Specular     Property
	Reflective   Property
	Transmissive Property
	Smoothness   Property
	Normal       Property
	IllumMap     *texture.Texture
	Index        float64
	IsLiquid     bool
	Brdf         BRDF
}

// New creates a new Material.
func New(opts *options.Material, fast bool) (*Material, error) {
	emiss, e := getProperty(opts.Emissive, fast)
	if e != nil {
		return nil, fmt.Errorf("creating property for Emissive: %v", e)
	}
	amb, e := getProperty(opts.Ambient, fast)
	if e != nil {
		return nil, fmt.Errorf("creating property for Ambient: %v", e)
	}
	dif, e := getProperty(opts.Diffuse, fast)
	if e != nil {
		return nil, fmt.Errorf("creating property for Diffuse: %v", e)
	}
	spec, e := getProperty(opts.Specular, fast)
	if e != nil {
		return nil, fmt.Errorf("creating property for Specular: %v", e)
	}
	refl, e := getProperty(opts.Reflective, fast)
	if e != nil {
		return nil, fmt.Errorf("creating property for Reflective: %v", e)
	}
	tran, e := getProperty(opts.Transmissive, fast)
	if e != nil {
		return nil, fmt.Errorf("creating property for Transmissive: %v", e)
	}
	smooth, e := getProperty(opts.Smoothness, fast)
	if e != nil {
		return nil, fmt.Errorf("creating property for Smoothness: %v", e)
	}

	brdf, ok := map[string]BRDF{
		"Lambert":     Lambert,
		"Blinn-Phong": BlinnPhong,
	}[opts.Brdf]
	if !ok {
		return nil, fmt.Errorf("unrecognized BRDF: %s", opts.Brdf)
	}

	illum, e := texture.NewEmpty(500, 500)
	if e != nil {
		return nil, fmt.Errorf("creating empty illumination map: %s", e)
	}

	return &Material{
		Emissive:     emiss,
		Ambient:      amb,
		Diffuse:      dif,
		Specular:     spec,
		Reflective:   refl,
		Transmissive: tran,
		Smoothness:   smooth,
		// Normal:,
		IllumMap: illum,
		Index:    opts.Index,
		IsLiquid: opts.IsLiquid,
		Brdf:     brdf,
	}, nil
}

func getProperty(matProp options.MatProperty, fast bool) (Property, error) {
	if fast {
		matProp.Type = "Uniform"
	}

	var prop Property
	switch matProp.Type {
	case "Uniform":
		prop = Uniform(matProp.Color)
	case "Texture":
		// TODO: implement textures
		prop = Uniform(matProp.Color)
	default:
		return nil, fmt.Errorf("unknown material property type: %s", matProp.Type)
	}
	return prop, nil
}

// Property is an interface that returns a color for a particular set of UV-coordinates.
type Property interface {
	Color(uv mgl64.Vec2) color64.Color64
}

// PropertyFunc is a function that implements the Property interface.
type PropertyFunc func(uv mgl64.Vec2) color64.Color64

// Color implements the Property interface for PropertyFunc.
func (u PropertyFunc) Color(uv mgl64.Vec2) color64.Color64 {
	return u(uv)
}

// Uniform creates a Property that always returns the given color.
func Uniform(color options.Color) Property {
	c := color64.Color64{color.R, color.G, color.B}
	return PropertyFunc(func(uv mgl64.Vec2) color64.Color64 {
		return c
	})
}

// func Texture() {
// }
