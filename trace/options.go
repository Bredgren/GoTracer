package trace

// Options contains all the options possible
type Options struct {
	Global     globalOpts      `title:"Global options" json:"global"`
	Resolution resolutionOpts  `title:"Image width and height in pixels" json:"resolution"`
	Background backgroundOpts  `title:"Determines the color when a ray doesn't hit anything" json:"background"`
	Camera     cameraOpts      `title:"Camera position and orientation" json:"camera"`
	AntiAlias  antiAliasOpts   `title:"Anti-aliasing settings" json:"anti_alias"`
	Debug      debugOpts       `title:"Debug options" json:"debug"`
	Lights     []*lightOpts    `title:"List of all lights in the scene" json:"lights"`
	Materials  []*materialOpts `title:"List of all the materials available to objects" json:"materials"`
	Objects    []*objectOpts   `title:"List of all objects in the scene" json:"objects"`
}

type globalOpts struct {
	FastRender       bool `title:"Disable/limit some settings in order to decrease rendering time" id:"fast-render" json:"fast_render"`
	MaxRecursion     int  `title:"Maximum reflective/refractive rays per pixel" min:"0" max:"99" class:"fast-render" json:"max_recursion"`
	SoftShadowDetail int  `title:"Soft shadow detail. 0 disables soft shadows" min:"0" max:"99" class:"fast-render" json:"soft_shadow_detail"`
}

type resolutionOpts struct {
	W int `title:"Width in pixels" min:"1" max:"16000" json:"w"`
	H int `title:"Height in pixels" min:"1" max:"16000" json:"h"`
}

type cameraOpts struct {
	Position vectorOpts `title:"Camera position" json:"position"`
	LookAt   vectorOpts `title:"The point that the camera is looking at" json:"look_at"`
	UpDir    vectorOpts `title:"The up direction" json:"up_dir"`
	Fov      float64    `title:"Field of view in degrees" min:"0.0" max:"360.0" json:"fov"`
	Dof      dofOpts    `title:"Depth of field" json:"dof"`
}

type dofOpts struct {
	Enabled           bool    `title:"Enable depth of field" class:"fast-render" json:"enabled"`
	FocalDistance     float64 `title:"Distance from the camera of the focal point" min:"0.01" max:"999.9" json:"focal_distance"`
	ApertureRadius    float64 `title:"Radius of the aperture" min:"0.00001" max:"999.9" json:"aperture_radius"`
	AdaptiveThreshold float64 `title:"Rays will continue to be created for each pixel until the contribution to the overall color is less than this" min:"0.00001" json:"adaptive_threshold"`
}

type backgroundOpts struct {
	Type  string    `title:"Type of background" choice:"Uniform,Skybox" json:"type"`
	Color colorOpts `title:"Color to use for a solid background (Uniform)" json:"color"`
	Image string    `title:"Path/URL of image to use (Skybox)" json:"image"`
}

type antiAliasOpts struct {
	MaxDivisions int     `title:"Maximum subdivisions of a pixel" min:"0" max:"16" class:"fast-render" json:"max_divisions"`
	Threshold    float64 `title:"Stop subdividing pixels when the difference is less than this" min:"0.0" max:"99" class:"fast-render" json:"threshold"`
}

type debugOpts struct {
	DebugRender bool      `title:"Produce a debug image" json:"debug_render"`
	Type        string    `title:"Type of debug image to produce" choice:"Ray Count,Anti Alias Subdivisions,Singe Pixel" json:"type"`
	SinglePixel bool      `title:"Render one pixel" json:"single_pixel"`
	Pixel       pixelOpts `title:"Pixel to render" json:"pixel"`
}

type pixelOpts struct {
	X int `title:"X value of the pixel" min:"0" valid:"validXPixel" json:"x"`
	Y int `title:"Y value of the pixel" min:"0" valid:"validYPixel" json:"y"`
}

type lightOpts struct {
	Type            string     `title:"Type of light" choice:"Directional,Point,Spot" json:"type"`
	Color           colorOpts  `title:"Color of the light" json:"color"`
	Position        vectorOpts `title:"Position of the light (Point and Spot)" json:"position"`
	Direction       vectorOpts `title:"Direction of the light (Direcitonal and Spot)" json:"direction"`
	IlluminationMap bool       `title:"Generate an illumination map (Point and Spot)" class:"fast-render" json:"illumination_map"`
	Coeff           struct {
		Constant  float64 `json:"constant"`
		Linear    float64 `json:"linear"`
		Quadratic float64 `json:"quadratic"`
	} `title:"Constant, lnear, and quadratic coefficients for point light fall off (Point)" json:"coeff"`
	Angle     float64 `title:"Spotlight angle" Type:"Spot" min:"0.0" max:"360.0" json:"angle"`
	DropOff   float64 `title:"Spotlight drop off angle" Type:"Spot" min:"0.0" max:"360.0" json:"drop_off"`
	FadeAngle float64 `title:"Stoplight fade angle" Type:"Spot" min:"0.0" max:"360.0" json:"fade_angle"`
}

type materialOpts struct {
	Name         string          `title:"Unique name for this material" json:"name"`
	Parent       string          `title:"Name of existing material to inherit attributes from" json:"parent"`
	Emissive     matPropertyOpts `title:"Emissive color of the material" json:"emissive"`
	Ambient      matPropertyOpts `title:"Abmient color of the material" json:"ambient"`
	Diffuse      matPropertyOpts `title:"Diffuse color of the material" json:"diffuse"`
	Specular     matPropertyOpts `title:"Specular color of the material" json:"specular"`
	Reflective   matPropertyOpts `title:"Reflectivness of the material" json:"reflective"`
	Transmissive matPropertyOpts `title:"Transmissivness of the material" json:"transmissive"`
	Smoothness   matPropertyOpts `title:"Smoothness of the material. Affects size of speclar spots" json:"smoothness"`
	Index        float64         `title:"Refractive index of the material" json:"index"`
	Normal       string          `title:"Path/URL of image to use as a normal map" json:"normal"`
	IsLiquid     bool            `title:"Overlapping behavior is only defined for non-liquids inside liquids" json:"is_liquid"`
	Brdf         string          `title:"Shadding algorithm" choice:"Lambert,Blinn-Phong" json:"brdf"`
}

type matPropertyOpts struct {
	Type    string    `title:"Type of material" choice:"Uniform,Texture" json:"type"`
	Color   colorOpts `title:"Uniform color" json:"color"`
	Texture string    `title:"Path to texture file" json:"texture"`
}

type objectOpts struct {
	Type         string        `title:"Type of shape. The 'Transform' type is an invisible object" choice:"Transform,Sphere,Box,Plane,Triangle,Trimesh,Cylinder,Cone" json:"type"`
	Transform    transformOpts `title:"Tranform of the object" json:"transform"`
	Material     string        `title:"Name of the material to use" json:"material"`
	TopRadius    float64       `title:"Top radius for cone objects" json:"top_radius"`
	BottomRadius float64       `title:"Bottom radius for cone objects" json:"bottom_radius"`
	Capped       bool          `title:"Whether to cap the ends of cones/cylinders" json:"capped"`
	Children     []objectOpts  `title:"Child objects that inherit this one's transform" json:"children"`
}

type transformOpts struct {
	Translate   vectorOpts `title:"Translation" json:"translate"`
	RotateAxis  vectorOpts `title:"Axis to rotate around" json:"rotate_axis"`
	RotateAngle float64    `title:"Angle to rotate around the axis in degrees" min:"-360.0" max:"360.0" json:"rotate_angle"`
	Scale       vectorOpts `title:"Scale" json:"scale"`
}

type colorOpts struct {
	R float64 `min:"0.0" max:"1.0" step:"0.1" default:"0.0" json:"r"`
	G float64 `min:"0.0" max:"1.0" step:"0.1" default:"0.0" json:"g"`
	B float64 `min:"0.0" max:"1.0" step:"0.1" default:"0.0" json:"b"`
}

type vectorOpts struct {
	X float64 `min:"-999" max:"999" default:"0.0" json:"x"`
	Y float64 `min:"-999" max:"999" default:"0.0" json:"y"`
	Z float64 `min:"-999" max:"999" default:"0.0" json:"z"`
}

// NewOptions returns an intialized Options.
func NewOptions() *Options {
	opts := &Options{
		Global: globalOpts{
			FastRender:       true,
			MaxRecursion:     2,
			SoftShadowDetail: 0,
		},
		Resolution: resolutionOpts{600, 600},
		Camera: cameraOpts{
			Position: vectorOpts{0, 5, 10},
			LookAt:   vectorOpts{0, 0, 0},
			UpDir:    vectorOpts{0, 1, 0},
			Fov:      58,
			Dof: dofOpts{
				Enabled:           false,
				FocalDistance:     5.0,
				ApertureRadius:    0.001,
				AdaptiveThreshold: 0.1,
			},
		},
		Background: backgroundOpts{
			Type:  "Uniform",
			Color: colorOpts{0, 0, 0},
			Image: "",
		},
		AntiAlias: antiAliasOpts{
			MaxDivisions: 0,
			Threshold:    0,
		},
		Debug: debugOpts{
			DebugRender: false,
		},
		Lights: []*lightOpts{
			{
				Type:      "Directional",
				Color:     colorOpts{1, 1, 1},
				Direction: vectorOpts{-1, -1, -1},
			},
		},
		Materials: []*materialOpts{
			{
				Name:         "white",
				Emissive:     matPropertyOpts{Type: "Uniform", Color: colorOpts{0, 0, 0}},
				Ambient:      matPropertyOpts{Type: "Uniform", Color: colorOpts{0.3, 0.3, 0.3}},
				Diffuse:      matPropertyOpts{Type: "Uniform", Color: colorOpts{1, 1, 1}},
				Specular:     matPropertyOpts{Type: "Uniform", Color: colorOpts{1, 1, 1}},
				Reflective:   matPropertyOpts{Type: "Uniform", Color: colorOpts{0, 0, 0}},
				Transmissive: matPropertyOpts{Type: "Uniform", Color: colorOpts{0, 0, 0}},
				Smoothness:   matPropertyOpts{Type: "Uniform", Color: colorOpts{0.9, 0.9, 0.9}},
				Index:        1,
				Normal:       "",
				IsLiquid:     false,
				Brdf:         "Blinn-Phong",
			},
		},
		Objects: []*objectOpts{
			{
				Type: "Plane",
				Transform: transformOpts{
					Translate:   vectorOpts{0, 0, 0},
					RotateAxis:  vectorOpts{1, 0, 0},
					RotateAngle: -90,
					Scale:       vectorOpts{10, 10, 0},
				},
				Material: "white",
			},
			{
				Type: "Box",
				Transform: transformOpts{
					Translate:   vectorOpts{0, 1, 0},
					RotateAxis:  vectorOpts{0, 0, 0},
					RotateAngle: 0,
					Scale:       vectorOpts{2, 2, 2},
				},
				Material: "white",
			},
		},
	}
	return opts
}
