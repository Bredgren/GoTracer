package lib

// Options contains all the options possible
type Options struct {
	Global     global      `title:"Global options" json:"global"`
	Resolution resolution  `title:"Image width and height in pixels" json:"resolution"`
	Background background  `title:"Determines the color when a ray doesn't hit anything" json:"background"`
	Camera     camera      `title:"Camera position and orientation" json:"camera"`
	AntiAlias  antiAlias   `title:"Anti-aliasing settings" json:"antiAlias"`
	Lights     []*light    `title:"List of all lights in the scene" json:"lights"`
	Materials  []*material `title:"List of all the materials available to objects" json:"materials"`
	Objects    []*object   `title:"List of all objects in the scene" json:"objects"`
}

type global struct {
	FastRender       bool `title:"Disable/limit some settings in order to decrease rendering time" id:"fast-render" json:"fastRender"`
	MaxRecursion     int  `title:"Maximum reflective/refractive rays per pixel" min:"0" max:"99" class:"fast-render" json:"maxRecursion"`
	SoftShadowDetail int  `title:"Soft shadow detail. 0 disables soft shadows" min:"0" max:"99" class:"fast-render" json:"softShadowDetail"`
}

type resolution struct {
	W int `title:"Width in pixels" min:"1" max:"16000" json:"w"`
	H int `title:"Height in pixels" min:"1" max:"16000" json:"h"`
}

type camera struct {
	Position vector  `title:"Camera position" json:"position"`
	LookAt   vector  `title:"The point that the camera is looking at" json:"lookAt"`
	UpDir    vector  `title:"The up direction" json:"upDir"`
	Fov      float64 `title:"Field of view in degrees" min:"0.0" max:"360.0" json:"fov"`
	Dof      float64 `title:"Depth of field" min:"0.0" max:"999.9" class:"fast-render" json:"dof"`
}

type background struct {
	Type  string `title:"Type of background" choice:"Uniform,Skybox" json:"type"`
	Color color  `title:"Color to use for a solid background (Uniform)" json:"color"`
	Image string `title:"Path/URL of image to use (Skybox)" json:"image"`
}

type antiAlias struct {
	MaxDivisions int     `title:"Maximum subdivisions of a pixel" min:"0" max:"16" class:"fast-render" json:"maxDivisions"`
	Threshold    float64 `title:"Stop subdividing pixels when the difference is less than this" min:"0.0" max:"99" class:"fast-render" json:"threshold"`
}

type light struct {
	Type            string `title:"Type of light" choice:"Directional,Point,Spot" json:"type"`
	Color           color  `title:"Color of the light" json:"color"`
	Position        vector `title:"Position of the light (Point and Spot)" json:"position"`
	Direction       vector `title:"Direction of the light (Direcitonal and Spot)" json:"direction"`
	IlluminationMap bool   `title:"Generate an illumination map (Point and Spot)" class:"fast-render" json:"illuminationMap"`
	Coeff           struct {
		Constant  float64 `json:"constant"`
		Linear    float64 `json:"linear"`
		Quadratic float64 `json:"quadratic"`
	} `title:"Constant, lnear, and quadratic coefficients for point light fall off (Point)" json:"coeff"`
	Angle     float64 `title:"Spotlight angle" Type:"Spot" min:"0.0" max:"360.0" json:"angle"`
	DropOff   float64 `title:"Spotlight drop off angle" Type:"Spot" min:"0.0" max:"360.0" json:"dropOff"`
	FadeAngle float64 `title:"Stoplight fade angle" Type:"Spot" min:"0.0" max:"360.0" json:"fadeAngle"`
}

type material struct {
	Name         string      `title:"Unique name for this material" json:"name"`
	Parent       string      `title:"Name of existing material to inherit attributes from" json:"parent"`
	Emissive     matProperty `title:"Emissive color of the material" json:"emissive"`
	Ambient      matProperty `title:"Abmient color of the material" json:"ambient"`
	Diffuse      matProperty `title:"Diffuse color of the material" json:"diffuse"`
	Specular     matProperty `title:"Specular color of the material" json:"specular"`
	Reflective   matProperty `title:"Reflectivness of the material" json:"reflective"`
	Transmissive matProperty `title:"Transmissivness of the material" json:"transmissive"`
	Smoothness   matProperty `title:"Smoothness of the material. Affects size of speclar spots" json:"smoothness"`
	Index        float64     `title:"Refractive index of the material" json:"index"`
	Normal       string      `title:"Path/URL of image to use as a normal map" json:"normal"`
	IsLiquid     bool        `title:"Overlapping behavior is only defined for non-liquids inside liquids" json:"isLiquid"`
	Brdf         string      `title:"Shadding algorithm" choice:"Lambert,Blinn-Phong" json:"brdf"`
}

type matProperty struct {
	Type    string `title:"Type of material" choice:"Uniform,Texture" json:"type"`
	Color   color  `title:"Uniform color" json:"color"`
	Texture string `title:"Path to texture file" json:"texture"`
}

type object struct {
	Type         string    `title:"Type of shape. The 'Transform' type is an invisible object" choice:"Transform,Sphere,Box,Plane,Triangle,Trimesh,Cylinder,Cone" json:"type"`
	Transform    transform `title:"Tranform of the object" json:"transform"`
	Material     string    `title:"Name of the material to use" json:"material"`
	TopRadius    float64   `title:"Top radius for cone objects" json:"topRadius"`
	BottomRadius float64   `title:"Bottom radius for cone objects" json:"bottomRadius"`
	Capped       bool      `title:"Whether to cap the ends of cones/cylinders" json:"capped"`
	Children     []object  `title:"Child objects that inherit this one's transform" json:"children"`
}

type transform struct {
	Translate   vector  `title:"Translation" json:"translate"`
	RotateAxis  vector  `title:"Axis to rotate around" json:"rotateAxis"`
	RotateAngle float64 `title:"Angle to rotate around the axis in degrees" min:"-360.0" max:"360.0" json:"rotateAngle"`
	Scale       vector  `title:"Scale" json:"scale"`
}

type color struct {
	R float64 `min:"0.0" max:"1.0" step:"0.1" default:"0.0" json:"r"`
	G float64 `min:"0.0" max:"1.0" step:"0.1" default:"0.0" json:"g"`
	B float64 `min:"0.0" max:"1.0" step:"0.1" default:"0.0" json:"b"`
}

type vector struct {
	X float64 `min:"-999" max:"999" default:"0.0" json:"x"`
	Y float64 `min:"-999" max:"999" default:"0.0" json:"y"`
	Z float64 `min:"-999" max:"999" default:"0.0" json:"z"`
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
