package lib

// Options contains all the options possible
type Options struct {
	Global     global      `title:"Global options"`
	Resolution resolution  `title:"Image width and height in pixels"`
	Background background  `title:"Determines the color when a ray doesn't hit anything"`
	Camera     camera      `title:"Camera position and orientation"`
	AntiAlias  antiAlias   `title:"Anti-aliasing settings"`
	Lights     []*light    `title:"List of all lights in the scene"`
	Materials  []*material `title:"List of all the materials available to objects"`
	Objects    []*object   `title:"List of all objects in the scene"`
}

type global struct {
	FastRender       bool `title:"Disable/limit some settings in order to decrease rendering time"`
	MaxRecursion     int  `title:"Maximum reflective/refractive rays per pixel" min:"0" max:"99"`
	SoftShadowDetail int  `title:"Soft shadow detail. 0 disables soft shadows" min:"0" max:"99"`
}

type resolution struct {
	W int `title:"Width in pixels" min:"1" max:"16000"`
	H int `title:"Height in pixels" min:"1" max:"16000"`
}

type camera struct {
	Position vector  `title:"Camera position"`
	LookAt   vector  `title:"The point that the camera is looking at"`
	UpDir    vector  `title:"The up direction"`
	Fov      float64 `title:"Field of view in degrees" min:"0.0" max:"360.0"`
	Dof      float64 `title:"Depth of field" min:"0.0" max:"999.9"`
}

type background struct {
	Type  string `title:"Type of background" choice:"Uniform,Skybox"`
	Color color  `title:"Color to use for a solid background (Uniform)"`
	Image string `title:"Path/URL of image to use (Skybox)"`
}

type antiAlias struct {
	MaxDivisions int     `title:"Maximum subdivisions of a pixel" min:"0" max:"16"`
	Threshold    float64 `title:"Stop subdividing pixels when the difference is less than this" min:"0.0" max:"99"`
}

type light struct {
	Type            string `title:"Type of light" choice:"Directional,Point,Spot"`
	Color           color  `title:"Color of the light"`
	Position        vector `title:"Position of the light (Point and Spot)"`
	Direction       vector `title:"Direction of the light (Direcitonal and Spot)"`
	IlluminationMap bool   `title:"Generate an illumination map (Point and Spot)"`
	Coeff           struct {
		Constant  float64
		Linear    float64
		Quadratic float64
	} `title:"Constant, lnear, and quadratic coefficients for point light fall off (Point)"`
	Angle     float64 `title:"Spotlight angle" Type:"Spot" min:"0.0" max:"360.0"`
	DropOff   float64 `title:"Spotlight drop off angle" Type:"Spot" min:"0.0" max:"360.0"`
	FadeAngle float64 `title:"Stoplight fade angle" Type:"Spot" min:"0.0" max:"360.0"`
}

type material struct {
	Name         string      `title:"Unique name for this material"`
	Parent       string      `title:"Name of existing material to inherit attributes from"`
	Emissive     matProperty `title:"Emissive color of the material"`
	Ambient      matProperty `title:"Abmient color of the material"`
	Diffuse      matProperty `title:"Diffuse color of the material"`
	Specular     matProperty `title:"Specular color of the material"`
	Reflective   matProperty `title:"Reflectivness of the material"`
	Transmissive matProperty `title:"Transmissivness of the material"`
	Smoothness   matProperty `title:"Smoothness of the material. Affects size of speclar spots"`
	Index        float64     `title:"Refractive index of the material"`
	Normal       string      `title:"Path/URL of image to use as a normal map"`
	IsLiquid     bool        `title:"Overlapping behavior is only defined for non-liquids inside liquids"`
	Brdf         string      `title:"Shadding algorithm" choice:"Lambert,Blinn-Phong"`
}

type matProperty struct {
	Type    string `title:"Type of material" choice:"Uniform,Texture"`
	Color   color  `title:"Uniform color"`
	Texture string `title:"Path to texture file"`
}

type object struct {
	Type         string    `title:"Type of shape. The 'Transform' type is an invisible object" choice:"Transform,Sphere,Box,Plane,Triangle,Trimesh,Cylinder,Cone"`
	Transform    transform `title:"Tranform of the object"`
	Material     string    `title:"Name of the material to use"`
	TopRadius    float64   `title:"Top radius for cone objects"`
	BottomRadius float64   `title:"Bottom radius for cone objects"`
	Capped       bool      `title:"Whether to cap the ends of cones/cylinders"`
	Children     []object  `title:"Child objects that inherit this one's transform"`
}

type transform struct {
	Translate   vector  `title:"Translation"`
	RotateAxis  vector  `title:"Axis to rotate around"`
	RotateAngle float64 `title:"Angle to rotate around the axis in degrees" min:"-360.0" max:"360.0"`
	Scale       vector  `title:"Scale"`
}

type color struct {
	R, G, B float64 `min:"0.0" max:"1.0" step:"0.1" default:"0.0"`
}

type vector struct {
	X, Y, Z float64 `min:"-999" max:"999" default:"0.0"`
}

// NewOptions returns an intialized Options.
func NewOptions() *Options {
	opts := &Options{
		Global: global{
			FastRender:       true,
			MaxRecursion:     2,
			SoftShadowDetail: 0,
		},
		Resolution: resolution{600, 600},
		Camera: camera{
			Position: vector{0, 5, 10},
			LookAt:   vector{0, 0, 0},
			UpDir:    vector{0, 1, 0},
			Fov:      58,
		},
		Background: background{
			Type:  "Uniform",
			Color: color{0, 0, 0},
			Image: "",
		},
		AntiAlias: antiAlias{
			MaxDivisions: 0,
			Threshold:    0,
		},
		Lights: []*light{
			{
				Type:      "Directional",
				Color:     color{1, 1, 1},
				Direction: vector{-1, -1, -1},
			},
		},
		Materials: []*material{
			{
				Name:         "white",
				Emissive:     matProperty{Type: "Uniform", Color: color{0, 0, 0}},
				Ambient:      matProperty{Type: "Uniform", Color: color{0.3, 0.3, 0.3}},
				Diffuse:      matProperty{Type: "Uniform", Color: color{1, 1, 1}},
				Specular:     matProperty{Type: "Uniform", Color: color{1, 1, 1}},
				Reflective:   matProperty{Type: "Uniform", Color: color{0, 0, 0}},
				Transmissive: matProperty{Type: "Uniform", Color: color{0, 0, 0}},
				Smoothness:   matProperty{Type: "Uniform", Color: color{0.9, 0.9, 0.9}},
				Index:        1,
				Normal:       "",
				IsLiquid:     false,
				Brdf:         "Blinn-Phong",
			},
		},
		Objects: []*object{
			{
				Type: "Plane",
				Transform: transform{
					Translate:   vector{0, 0, 0},
					RotateAxis:  vector{1, 0, 0},
					RotateAngle: -90,
					Scale:       vector{10, 10, 0},
				},
				Material: "white",
			},
			{
				Type: "Box",
				Transform: transform{
					Translate:   vector{0, 1, 0},
					RotateAxis:  vector{0, 0, 0},
					RotateAngle: 0,
					Scale:       vector{2, 2, 2},
				},
				Material: "white",
			},
		},
	}
	return opts
}
