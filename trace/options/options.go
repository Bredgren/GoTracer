package options

// Options contains all the options possible.
type Options struct {
	Global     Global      `title:"Global options" json:"global"`
	Resolution Resolution  `title:"Image width and height in pixels. Larger resolution will increase rendering time." json:"resolution"`
	Background Background  `title:"Determines the color when a ray doesn't hit anything" json:"background"`
	Camera     Camera      `title:"Camera position and orientation" json:"camera"`
	AntiAlias  AntiAlias   `title:"Anti-aliasing settings" json:"anti_alias"`
	Debug      Debug       `title:"Debug options" json:"debug"`
	Lights     []*Light    `title:"List of all lights in the scene" json:"lights"`
	Materials  []*Material `title:"List of all the materials available to objects" json:"materials"`
	Objects    []*Object   `title:"List of all objects in the scene" json:"objects"`
	Layout     []*Layout   `title:"Placement of objects" json:"layout"`
}

// Global contains global options.
type Global struct {
	FastRender       bool `title:"Disable/limit some settings in order to decrease rendering time" id:"fast-render" json:"fast_render"`
	MaxRecursion     int  `title:"Maximum reflective/refractive rays per pixel. Larger values increase quality and rendering time." min:"0" max:"99" class:"fast-render" json:"max_recursion"`
	SoftShadowDetail int  `title:"Soft shadow detail. 0 disables soft shadows. Larger values increase quality and rendering time. " min:"0" max:"99" class:"fast-render" json:"soft_shadow_detail"`
}

// Resolution defines the dimensions of the image.
type Resolution struct {
	W int `title:"Width in pixels. Larger values increase rendering time." min:"1" max:"16000" json:"w"`
	H int `title:"Height in pixels. Larger values increase rendering time." min:"1" max:"16000" json:"h"`
}

// Camera defines the camera position and other options.
type Camera struct {
	Position Vector  `title:"Camera position" json:"position"`
	LookAt   Vector  `title:"The point that the camera is looking at" json:"look_at"`
	UpDir    Vector  `title:"The up direction" json:"up_dir"`
	Fov      float64 `title:"Field of view in degrees" min:"0.1" max:"360.0" json:"fov"`
	Dof      Dof     `title:"Depth of field" json:"dof"`
}

// Dof defines the depth of field options.
type Dof struct {
	Enabled           bool    `title:"Enable depth of field. Enabling DOF will increase rendering time." class:"fast-render" json:"enabled"`
	FocalDistance     float64 `title:"Distance from the camera of the focal point" min:"0.01" max:"999.9" json:"focal_distance"`
	ApertureRadius    float64 `title:"Radius of the aperture" min:"0.00001" max:"999.9" json:"aperture_radius"`
	MaxRays           int     `title:"Maximum DOF rays to cast per pixel" min:"1" json:"max_rays"`
	AdaptiveThreshold float64 `title:"Rays will continue to be cast for each pixel until the contribution to the overall color is less than this, or until MaxRays is reached, whichever is first. Smaller values increase quality and rendering time." min:"0.00001" step:"0.01" json:"adaptive_threshold"`
}

// Background defines the options for rays that don't hit objects.
type Background struct {
	Type  string `title:"Type of background" choice:"Uniform,Skybox" json:"type"`
	Color Color  `title:"Color to use for a solid background (Uniform)" json:"color"`
	Image string `title:"Path/URL of image to use (Skybox)" json:"image"`
}

// AntiAlias defines the anti-alias options.
type AntiAlias struct {
	MaxDivisions int     `title:"Maximum subdivisions of a pixel. Larger values increase quality and rendering time." min:"0" max:"16" class:"fast-render" json:"max_divisions"`
	Threshold    float64 `title:"Stop subdividing pixels when the difference is less than this. Smaller values increase quality and rendering time." min:"0.0" max:"99" step:"0.1" class:"fast-render" json:"threshold"`
}

// Debug defines special options for producing images for debugging.
type Debug struct {
	DebugRender bool   `title:"Produce a debug image" json:"debug_render"`
	Type        string `title:"Type of debug image to produce" choice:"Ray Count,Anti Alias Subdivisions,Singe Pixel" json:"type"`
	SinglePixel bool   `title:"Render one pixel" json:"single_pixel"`
	Pixel       Pixel  `title:"Pixel to render" json:"pixel"`
}

// Pixel specifies a single pixel.
type Pixel struct {
	X int `title:"X value of the pixel" min:"0" valid:"validXPixel" json:"x"`
	Y int `title:"Y value of the pixel" min:"0" valid:"validYPixel" json:"y"`
}

// Light defines options for single light source.
type Light struct {
	Type            string `title:"Type of light" choice:"Directional,Point,Spot" json:"type"`
	Color           Color  `title:"Color of the light" json:"color"`
	Position        Vector `title:"Position of the light (Point and Spot)" json:"position"`
	Direction       Vector `title:"Direction of the light (Direcitonal and Spot)" json:"direction"`
	IlluminationMap bool   `title:"Generate an illumination map (Point and Spot). Increases rendering time." class:"fast-render" json:"illumination_map"`
	Coeff           struct {
		Constant  float64 `json:"constant"`
		Linear    float64 `json:"linear"`
		Quadratic float64 `json:"quadratic"`
	} `title:"Constant, lnear, and quadratic coefficients for point light fall off (Point)" json:"coeff"`
	Angle     float64 `title:"Spotlight angle" Type:"Spot" min:"0.0" max:"360.0" json:"angle"`
	DropOff   float64 `title:"Spotlight drop off angle" Type:"Spot" min:"0.0" max:"360.0" json:"drop_off"`
	FadeAngle float64 `title:"Stoplight fade angle" Type:"Spot" min:"0.0" max:"360.0" json:"fade_angle"`
}

// Material defines options for materials.
type Material struct {
	Name         string      `title:"Unique name for this material" json:"name"`
	Parent       string      `title:"Name of existing material to inherit attributes from" json:"parent"`
	Emissive     MatProperty `title:"Emissive color of the material" json:"emissive"`
	Ambient      MatProperty `title:"Abmient color of the material" json:"ambient"`
	Diffuse      MatProperty `title:"Diffuse color of the material" json:"diffuse"`
	Specular     MatProperty `title:"Specular color of the material" json:"specular"`
	Reflective   MatProperty `title:"Reflectivness of the material" json:"reflective"`
	Transmissive MatProperty `title:"Transmissivness of the material" json:"transmissive"`
	Smoothness   MatProperty `title:"Smoothness of the material. Affects size of speclar spots" json:"smoothness"`
	Index        float64     `title:"Refractive index of the material" step:"0.1" json:"index"`
	Normal       string      `title:"Path/URL of image to use as a normal map" json:"normal"`
	IsLiquid     bool        `title:"Overlapping behavior is only defined for non-liquids inside liquids" json:"is_liquid"`
	Brdf         string      `title:"Shadding algorithm" choice:"Lambert,Blinn-Phong" json:"brdf"`
}

// MatProperty defines the properties of a single material option.
type MatProperty struct {
	Type    string  `title:"Type of material" choice:"Uniform,Texture" json:"type"`
	Color   Color   `title:"Uniform color" json:"color"`
	Texture Texture `title:"Texture options if using the 'Texture' type" json:"texture"`
}

// Texture defines options for a texture.
type Texture struct {
	Source  string  `title:"Path to texture file" json:"source"`
	ScaleX  float64 `title:"Scale of the texture in the X direction. Note that 0 is equivalent to 1." min:"0.0" step:"0.1" json:"scale_x"`
	ScaleY  float64 `title:"Scale of the texture in the Y direction. Note that 0 is equivalent to 1." min:"0.0" step:"0.1" json:"scale_y"`
	OffsetX float64 `title:"Offset of the texture in the X direction" step:"0.1" json:"offset_x"`
	OffsetY float64 `title:"Offset of the texture in the Y direction" step:"0.1" json:"offset_y"`
	Tile    bool    `title:"Whether or not to tile the texture" json:"tile"`
}

// Object definse options for a single object.
type Object struct {
	Name         string    `title:"Name of new or existing object" json:"name"`
	Type         string    `title:"Type of shape. The 'Transform' type is an invisible object" choice:"Transform,Sphere,Cube,Plane,Triangle,Trimesh,Cylinder,Cone" json:"type"`
	Transform    Transform `title:"Tranform of the object" json:"transform"`
	Material     string    `title:"Name of the material to use" json:"material"`
	TopRadius    float64   `title:"Top radius for cone objects" step:"0.1" json:"top_radius"`
	BottomRadius float64   `title:"Bottom radius for cone objects" step:"0.1" json:"bottom_radius"`
	Capped       bool      `title:"Whether to cap the ends of cones/cylinders" json:"capped"`
	Children     []*Object `title:"Child objects that inherit this one's transform" json:"children"`
}

// Layout defines transform options for an object.
type Layout struct {
	Transform Transform `title:"Tranform of the object" json:"transform"`
	Name      string    `title:"Name of existing object from the list of objects" json:"name"`
}

// Transform defines the transform for an object.
type Transform struct {
	Translate   Vector  `title:"Translation" json:"translate"`
	RotateAxis  Vector  `title:"Axis to rotate around" json:"rotate_axis"`
	RotateAngle float64 `title:"Angle to rotate around the axis in degrees" min:"-360.0" max:"360.0" json:"rotate_angle"`
	Scale       Vector  `title:"Scale" json:"scale"`
}

// Color is and RGB color.
type Color struct {
	R float64 `min:"0.0" max:"1.0" step:"0.1" default:"0.0" json:"r"`
	G float64 `min:"0.0" max:"1.0" step:"0.1" default:"0.0" json:"g"`
	B float64 `min:"0.0" max:"1.0" step:"0.1" default:"0.0" json:"b"`
}

// Vector is a 3D vector.
type Vector struct {
	X float64 `min:"-999" max:"999" default:"0.0" json:"x"`
	Y float64 `min:"-999" max:"999" default:"0.0" json:"y"`
	Z float64 `min:"-999" max:"999" default:"0.0" json:"z"`
}

// NewOptions returns an intialized Options.
func NewOptions() *Options {
	opts := &Options{
		Global: Global{
			FastRender:       true,
			MaxRecursion:     2,
			SoftShadowDetail: 0,
		},
		Resolution: Resolution{600, 600},
		Camera: Camera{
			Position: Vector{0, 5, 10},
			LookAt:   Vector{0, 0, 0},
			UpDir:    Vector{0, 1, 0},
			Fov:      58,
			Dof: Dof{
				Enabled:           false,
				FocalDistance:     5.0,
				ApertureRadius:    0.001,
				MaxRays:           1,
				AdaptiveThreshold: 0.1,
			},
		},
		Background: Background{
			Type:  "Uniform",
			Color: Color{0, 0, 0},
			Image: "",
		},
		AntiAlias: AntiAlias{
			MaxDivisions: 0,
			Threshold:    0,
		},
		Debug: Debug{
			DebugRender: false,
		},
		Lights: []*Light{
			{
				Type:      "Directional",
				Color:     Color{1, 1, 1},
				Direction: Vector{-1, -1, -1},
			},
		},
		Materials: []*Material{
			{
				Name:         "white",
				Emissive:     MatProperty{Type: "Uniform", Color: Color{0, 0, 0}},
				Ambient:      MatProperty{Type: "Uniform", Color: Color{0.3, 0.3, 0.3}},
				Diffuse:      MatProperty{Type: "Uniform", Color: Color{1, 1, 1}},
				Specular:     MatProperty{Type: "Uniform", Color: Color{1, 1, 1}},
				Reflective:   MatProperty{Type: "Uniform", Color: Color{0, 0, 0}},
				Transmissive: MatProperty{Type: "Uniform", Color: Color{0, 0, 0}},
				Smoothness:   MatProperty{Type: "Uniform", Color: Color{0.9, 0.9, 0.9}},
				Index:        1,
				Normal:       "",
				IsLiquid:     false,
				Brdf:         "Blinn-Phong",
			},
		},
		Objects: []*Object{
			{
				Name: "floor",
				Type: "Plane",
				Transform: Transform{
					Translate:   Vector{0, 0, 0},
					RotateAxis:  Vector{1, 0, 0},
					RotateAngle: -90,
					Scale:       Vector{10, 10, 0},
				},
				Material: "white",
			},
			{
				Name: "box",
				Type: "Cube",
				Transform: Transform{
					Translate:   Vector{0, 1, 0},
					RotateAxis:  Vector{0, 1, 0},
					RotateAngle: 45,
					Scale:       Vector{2, 2, 2},
				},
				Material: "white",
			},
		},
		Layout: []*Layout{
			{
				Name: "floor",
			},
			{
				Name: "box",
			},
		},
	}
	return opts
}
