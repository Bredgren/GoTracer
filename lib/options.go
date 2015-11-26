package lib

// Options contains all the options possible
type Options struct {
	Resolution struct {
		W, H int `min:"1" max:"16000" default:"600"`
	} `desc:"Image width and height in pixels"`
	Camera struct {
		Position Vector `desc:"Camera position"`
		LookAt   Vector `desc:"The point that the camera is looking at"`
		UpDir    Vector `desc:"The up direction"`
		Fov      Angle  `desc:"Field of view"`
	} `desc:"Camera position and orientation"`
	Background struct {
		Type  string `desc:"Type of background" Types:"Uniform,Skybox" default:"Uniform"`
		Color Color  `desc:"Color to use for a solid background" Type:"Uniform"`
		Image string `desc:"Path/URL of image to use as skybox" Type:"Skybox"`
	} `desc:"Determines the color when a ray doesn't hit anything"`
	AntiAlias struct {
		MaxDivisions int     `desc:"Maximum subdivisions of a pixel" min:"0" max:"16" default:"0"`
		Threshold    float64 `desc:"Stop subdividing pixels when the difference is less than this" min:"0.0" max:"99"`
	} `desc:"Anti-aliasing settings"`
	MaxRecursion int         `desc:"Maximum reflective and refractive rays per pixel" min:"0" max:"99" default:"2"`
	SoftShadows  int         `desc:"Soft shadow detail. 0 disables soft shadows" min:"0" max:"99" default:"0"`
	Lights       []*Light    `desc:"List of all lights in the scene"`
	Materials    []*Material `desc:"List of all the materials available to objects"`
	Objects      []*Object   `dest:"List of all objects in the scene"`
}

// Light is a light source in the scene.
type Light struct {
	Type      string `desc:"Type of light" Types:"Directional,Point,Spot"`
	Color     Color  `desc:"Color of the light"`
	Position  Vector `desc:"Position of the light"`
	Direction Vector `desc:"Direction of the light" Type:"Directional,Spot"`
	Coeff     struct {
		Constant  float64
		Linear    float64
		Quadratic float64
	} `desc:"Constant, lnear, and quadratic coefficients for point light fall off" Type:"Point"`
	Angle     float64 `desc:"Spotlight angle" Type:"Spot"`
	DropOff   float64 `desc:"Spotlight drop off angle" Type:"Spot"`
	FadeAngle float64 `desc:"Stoplight fade angle" Type:"Spot"`
}

// Material is a description of a material that can be used by Objects.
type Material struct {
	Name         string      `desc:"Unique name for the material"`
	Parent       string      `desc:"Name of existing material to inherit attributes from"`
	Emissive     MatProperty `desc:"Emissive color of the material"`
	Ambient      MatProperty `desc:"Abmient color of the material"`
	Diffuse      MatProperty `desc:"Diffuse color of the material"`
	Specular     MatProperty `desc:"Specular color of the material"`
	Reflective   MatProperty `desc:"Reflectivness of the material"`
	Transmissive MatProperty `desc:"Transmissivness of the material"`
	Smoothness   MatProperty `desc:"Smoothness of the material. Affects size of speclar spots"`
	Index        float64     `desc:"Refractive index of the material"`
	Normal       string      `desc:"Path/URL of image to use as a normal map"`
	IsLiquid     bool        `desc:"Overlapping behavior is only defined for non-liquids inside liquids"`
	Brdf         string      `desc:"Lambert, Blinn-Phong"`
}

// Object is a rendreable object.
type Object struct {
	Type         string   `desc:"Type of shape" Types:"Transform,Sphere,Box,Square,Cylinder,Cone"`
	Translate    Vector   `desc:"Transform translation" Type:"Transform"`
	RotateAxis   Vector   `desc:"Transform axis to rotate around" Type:"Transform"`
	RotateAngle  Angle    `desc:"Transform angle to rotate around the axis" Type:"Transform"`
	Scale        Vector   `desc:"Transform scale" Type:"Transform"`
	SubObjects   []Object `desc:"Objects to apply the transform to" Type:"Transform"`
	Material     string   `desc:"Name of the material to use" Type:"Sphere,Box,Square,Cylinder,Cone"`
	TopRadius    float64  `desc:"Top radius of the cone" Type:"Cone"`
	BottomRadius float64  `desc:"Bottom radius of the cone" Type:"Cone"`
	Capped       bool     `desc:"Whether to cap the cone/cylinder" Type:"Cynlinder,Cone"`
}

// MatProperty is a Material property which can either be a texture or a uniform color
type MatProperty struct {
	Type    string `desc:"Type of material" Types:"Uniform,Texture"`
	Texture string `desc:"Path to texture file" Type:"Texture"`
	Color   Color  `desc:"Uniform color" Type:"Uniform"`
}

// type Enum interface {
// 	EnumVals() []int
// }

// type Background int

// const (
// 	Uniform Background = iota
// 	Skybox
// )

// func (b Background) String() string {
// 	switch b {
// 	case Uniform:
// 		return "Uniform"
// 	case Skybox:
// 		return "Skybox"
// 	}
// 	return ""
// }

// func (b Background) EnumVals() []int {
// 	return []int{int(Uniform), int(Skybox)}
// }

// Angle is in degrees.
type Angle struct {
	Degrees float64 `min:"0.0" max:"360.0"`
}

// Color is an RGB color.
type Color struct {
	R, G, B float64 `min:"0.0" max:"1.0" default:"0.0"`
}

// Vector is a 3D vector in space.
type Vector struct {
	X, Y, Z float64 `min:"-999" max:"999" default:"0.0"`
}

// Path is a URL or file path.
// type Path string

// NewOptions returns an intialized Options.
func NewOptions() *Options {
	opts := &Options{}
	// w, _ := reflect.TypeOf(opts.Resolution).FieldByName("W")
	// opts.Resolution.W, _ = strconv.Atoi(w.Tag.Get("default"))
	// h, _ := reflect.TypeOf(opts.Resolution).FieldByName("H")
	// opts.Resolution.H, _ = strconv.Atoi(h.Tag.Get("default"))

	// p, _ := reflect.TypeOf(opts.Camera).FieldByName("Position")
	// opts.Camera.Position.X, _ = strconv.Atoi(p.Tag.Get("default"))
	// Camera struct {
	// 	Position Vector  `desc:"Camera position"`
	// 	LookAt   Vector  `desc:"The point that the camera is looking at"`
	// 	UpDir    Vector  `desc:"The up direction"`
	// 	Fov      Degrees `desc:"Field of view" min:"1" max:"180" default:"53"`
	// } `desc:"Camera position and orientation"`
	opts.Materials = append(opts.Materials, &Material{Name: "abc"})
	opts.Materials = append(opts.Materials, &Material{Name: "def"})
	return opts
}
