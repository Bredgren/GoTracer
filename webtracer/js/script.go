package main

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/Bredgren/gotracer/lib"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
)

var jq = jquery.NewJQuery
var console = js.Global.Get("console")

func main() {
	js.Global.Set("onBodyLoad", onBodyLoad)
}

func onBodyLoad() {
	initCallbacks()
	options := lib.NewOptions()
	initOptions(options)

	imgCon := jq(".image-container")
	console.Call("log", imgCon)
	img := js.Global.Get("Image").New()
	img.Set("onload", func() {
		imgCon.Append(img)
		setImageSize(img.Get("width"), img.Get("height"))
	})
	img.Set("src", "/img/render542.png")
}

func initCallbacks() {
	jq(".tab").Call(jquery.CLICK, onToggleControls)
	jq("#fast-render").Call(jquery.CLICK, onToggleMode)
	jq("#save").Call(jquery.CLICK, onSave)
	jq("#load").Call(jquery.CLICK, onLoad)
	jq("#download").Call(jquery.CLICK, onDownload)
	jq(".option-group").Call(jquery.CLICK, func(o *js.Object) {
		// j.SlideToggle()
		console.Call("log", o)
	})
}

func setImageSize(w, h *js.Object) {
	imgCon := jq(".image-container")
	imgCon.SetCss("width", w)
	imgCon.SetCss("height", h)
}

func initOptions(options *lib.Options) {
	opts := jq("#options")
	console.Call("log", opts)
	console.Call("log", options)
	st := reflect.TypeOf(*options)
	sv := reflect.ValueOf(options).Elem()
	for i := 0; i < st.NumField(); i++ {
		opts.Append(handleOptionField(st.Field(i), sv.Field(i)))
	}
}

func handleOptionField(field reflect.StructField, val reflect.Value) (opts jquery.JQuery) {
	switch field.Type.Kind() {
	case reflect.Struct:
		opts = optionGroup(field.Name)
		g := jq("<div>").AddClass("option-group-options")
		t := field.Type
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			v := val.Field(i)
			g.Append(handleOptionField(f, v).SetAttr("title", f.Tag.Get("desc")))
		}
		opts.Call(jquery.CLICK, func(event jquery.Event) {
			g.SlideToggle("fast")
			event.StopPropagation()
		})
		opts.Append(g)
	case reflect.Slice:
		console.Call("log", "is slice")
		opts = vectorOption()
	case reflect.Int, reflect.Float64:
		if field.Type.Implements(reflect.TypeOf((*lib.Enum)(nil)).Elem()) {
			console.Call("log", field.Type.Name(), "is enum")
			// opts = choiceOption(field.Name)
		} else {
			opts = numberOption(field.Name, field.Tag.Get("min"), field.Tag.Get("max"), field.Tag.Get("default"))
			opts.Call(jquery.CHANGE, func(event jquery.Event) {
				if field.Type.Kind() == reflect.Int {
					i, e := strconv.Atoi(event.Target.Get("value").String())
					if e != nil {
						panic(e)
					}
					val.SetInt(int64(i))
				} else {
					i, e := strconv.ParseFloat(event.Target.Get("value").String(), 64)
					if e != nil {
						panic(e)
					}
					val.SetFloat(i)
				}
			})
		}
	case reflect.String:
		if field.Name == "Type" {
			opts = choiceOption(field.Name, strings.Split(field.Tag.Get("Types"), ","))
			opts.Call(jquery.CHANGE, func(event jquery.Event) {
				console.Call("log", event.Target.Get("value").String())
			})
		} else {
			opts = textOption(field.Name)
		}
	default:
		console.Call("log", "is other")
		opts = numberOption(field.Name, "1", "1", "1")
	}
	opts.SetAttr("title", field.Tag.Get("desc"))
	return
}

// type group struct {
// }

// type gOption struct {
// 	structField
// 	uiField
// 	parentGroup
// }

func optionGroup(title string) jquery.JQuery {
	g := jq("<div>").AddClass("option-group")
	h := jq("<div>").AddClass("option-group-header").SetText(title)
	return g.Append(h)
}

func numberOption(label, min, max, value string) jquery.JQuery {
	o := option(label)
	in := jq("<input>").SetAttr("type", "number").SetAttr("min", min).SetAttr("max", max).SetAttr("value", value)
	return o.Append(in)
}

func textOption(label string) jquery.JQuery {
	o := option(label)
	in := jq("<input>").SetAttr("type", "text")
	return o.Append(in)
}

func choiceOption(label string, choices []string) jquery.JQuery {
	o := option(label)
	s := jq("<select>")
	for _, c := range choices {
		s.Append(jq("<option>").SetText(c))
	}
	return o.Append(s)
}

func option(label string) jquery.JQuery {
	o := jq("<div>").AddClass("option")
	o.Append(label + " ")
	o.Call(jquery.CLICK, func(event jquery.Event) {
		event.StopPropagation()
	})
	return o
}

func vectorOption() jquery.JQuery {
	return jq("<span>") //.Append(numberOption("x", -999, 999, 0)).Append(numberOption("y", -999, 999, 0)).
	//	Append(numberOption("z", -999, 999, 0))
}

// func initOptions() {
// 	opts := jq("#options")
// 	// initDebugOptions(opts)
// 	g := globalOptionGroup()
// 	opts.Append(g)
// }

// func globalOptionGroup() jquery.JQuery {
// 	global := optionGroup("Global")
// 	g := subOptionGroup("Resolution")
// 	g.Append(numberOption("w", 1, 16000, 800))
// 	g.Append(numberOption("h", 1, 16000, 600))
// 	global.Append(g)
// 	g = subOptionGroup("Camera Position")
// 	g.Append(vectorOption())
// 	global.Append(g)
// 	g = subOptionGroup("Camera LookAt")
// 	g.Append(vectorOption())
// 	global.Append(g)
// 	g = subOptionGroup("Camera Up Direction")
// 	g.Append(vectorOption())
// 	global.Append(g)
// 	g = subOptionGroup("Camera FOV")
// 	g.Append(numberOption("angle", 1, 180, 53))
// 	global.Append(g)
// 	g = subOptionGroup("Background")
// 	// g.Append(colorOption())
// 	global.Append(g)
// 	return global
// }

// func subOptionGroup(label string) jquery.JQuery {
// 	l := jq("<div>").AddClass("suboption-label").SetText(label)
// 	return jq("<div>").AddClass("suboption-group").Append(l)
// }

func onToggleControls() {
	console.Call("log", "toggle")
	c := jq(".controls")
	if c.HasClass("open") {
		c.RemoveClass("open")
		c.AddClass("closed")
	} else {
		c.RemoveClass("closed")
		c.AddClass("open")
	}
}

func onToggleMode() {
	btn := jq("#fast-render")
	switch {
	case btn.HasClass("toggle-button-down"):
		btn.RemoveClass("toggle-button-down")
		btn.AddClass("toggle-button-up")
	case btn.HasClass("toggle-button-up"):
		btn.RemoveClass("toggle-button-up")
		btn.AddClass("toggle-button-down")
	}
}

func onSave() {
	console.Call("log", "save")
}

func onLoad() {
	console.Call("log", "load")
}

func onDownload() {
	console.Call("log", "download")
}
