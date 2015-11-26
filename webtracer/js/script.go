package main

import (
	"reflect"

	"github.com/Bredgren/gohtmlctrl/htmlctrl"
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
	console.Call("log", options)

	imgCon := jq(".image-container")
	img := js.Global.Get("Image").New()
	img.Set("onload", func() {
		imgCon.Append(img)
		console.Call("log", img.Get("width"), img.Get("height"))
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
}

func setImageSize(w, h *js.Object) {
	imgCon := jq(".image-container")
	imgCon.SetCss("width", w)
	imgCon.SetCss("height", h)
}

type optionItem struct {
	jq jquery.JQuery
	st reflect.StructField
}

func initOptions(options *lib.Options) {
	opts := jq("#options")
	o, e := htmlctrl.Struct(options, "Options", "all-options", "")
	if e != nil {
		console.Call("error", e.Error())
		return
	}
	opts.Append(o)
	// st := reflect.TypeOf(*options)
	// sv := reflect.ValueOf(options).Elem()
	// var siblings []optionItem
	// for i := 0; i < st.NumField(); i++ {
	// 	o := handleOptionField(st.Field(i), sv.Field(i), &siblings)
	// 	siblings = append(siblings, optionItem{jq: o, st: st.Field(i)})
	// 	opts.Append(o)
	// }
}

// func handleOptionField(field reflect.StructField, val reflect.Value,
// 	siblings *[]optionItem) (opts jquery.JQuery) {
// 	d := field.Tag.Get("default")
// 	switch field.Type.Kind() {
// 	case reflect.Struct:
// 		opts = optionGroup(field.Name)
// 		g := jq("<div>").AddClass("option-group-options")
// 		t := field.Type
// 		var sibs []optionItem
// 		for i := 0; i < t.NumField(); i++ {
// 			f := t.Field(i)
// 			v := val.Field(i)
// 			o := handleOptionField(f, v, &sibs).SetAttr("title", f.Tag.Get("desc"))
// 			sibs = append(sibs, optionItem{jq: o, st: f})
// 			g.Append(o)
// 		}
// 		opts.Call(jquery.CLICK, func(event jquery.Event) {
// 			g.SlideToggle("fast")
// 			event.StopPropagation()
// 		})
// 		opts.Append(g)
// 	case reflect.Slice:
// 		console.Call("log", "handle slice", field.Name)
// 		opts = optionList(field.Name)
// 		items := jq("<div>").AddClass("option-list-items")
// 		console.Call("log", "val.Len", val.Len())
// 		for i := 0; i < val.Len(); i++ {
// 			console.Call("log", "i", i)
// 			// item := jq("<div>").AddClass("option-list-item")
// 			item := optionGroup(string(byte('0') + byte(i)))
// 			g := jq("<div>").AddClass("option-group-options")
// 			sv := val.Index(i)
// 			console.Call("log", "sv", sv)
// 			// console.Call("log", "sv", sv.Elem().Field(0))
// 			st := sv.Elem().Type()
// 			console.Call("log", "st", st.Name())
// 			console.Call("log", "st fields", st.NumField())
// 			var sibs []optionItem
// 			for i := 0; i < st.NumField(); i++ {
// 				f := st.Field(i)
// 				v := sv.Elem().Field(i)
// 				o := handleOptionField(f, v, &sibs).SetAttr("title", f.Tag.Get("desc"))
// 				sibs = append(sibs, optionItem{jq: o, st: f})
// 				g.Append(o)
// 			}
// 			item.Call(jquery.CLICK, func(event jquery.Event) {
// 				g.SlideToggle("fast")
// 				event.StopPropagation()
// 			})
// 			item.Append(g)
// 			console.Call("log", item)
// 			items.Append(item)
// 		}
// 		opts.Call(jquery.CLICK, func(event jquery.Event) {
// 			items.SlideToggle("fast")
// 			event.StopPropagation()
// 		})
// 		opts.Append(items)
// 	case reflect.Int:
// 		opts = numberOption(field.Name, field.Tag.Get("min"), field.Tag.Get("max"), d)
// 		opts.Call(jquery.CHANGE, func(event jquery.Event) {
// 			i, e := strconv.Atoi(event.Target.Get("value").String())
// 			if e != nil {
// 				panic(e)
// 			}
// 			val.SetInt(int64(i))
// 		})
// 		if d == "" {
// 			d = "0"
// 		}
// 		i, e := strconv.Atoi(d)
// 		if e != nil {
// 			panic(e)
// 		}
// 		val.SetInt(int64(i))
// 	case reflect.Float64:
// 		opts = numberOption(field.Name, field.Tag.Get("min"), field.Tag.Get("max"), d)
// 		opts.Call(jquery.CHANGE, func(event jquery.Event) {
// 			i, e := strconv.ParseFloat(event.Target.Get("value").String(), 64)
// 			if e != nil {
// 				panic(e)
// 			}
// 			val.SetFloat(i)
// 		})
// 		if d == "" {
// 			d = "0.0"
// 		}
// 		i, e := strconv.ParseFloat(d, 64)
// 		if e != nil {
// 			panic(e)
// 		}
// 		val.SetFloat(i)
// 	case reflect.String:
// 		if field.Name == "Type" {
// 			opts = choiceOption(field.Name, strings.Split(field.Tag.Get("Types"), ","))
// 			opts.Call(jquery.CHANGE, func(event jquery.Event) {
// 				val := event.Target.Get("value").String()
// 			outer:
// 				for _, s := range *siblings {
// 					if s.st.Name != field.Name {
// 						for _, t := range strings.Split(s.st.Tag.Get("Type"), ",") {
// 							if t == val {
// 								s.jq.SlideDown("fast")
// 								continue outer
// 							}
// 						}
// 						s.jq.SlideUp("fast")
// 					}
// 				}
// 			})
// 		} else {
// 			opts = textOption(field.Name)
// 		}
// 		val.SetString(d)
// 	default:
// 		console.Call("log", "is other")
// 		opts = numberOption(field.Name, "1", "1", "1")
// 	}
// 	opts.SetAttr("title", field.Tag.Get("desc"))
// outer:
// 	for _, s := range *siblings {
// 		if s.st.Name == "Type" {
// 			for _, t := range strings.Split(field.Tag.Get("Type"), ",") {
// 				if t == s.jq.Children("select").First().Val() {
// 					opts.Show()
// 					break outer
// 				}
// 			}
// 			opts.Hide()
// 		}
// 	}
// 	return
// }

// func optionGroup(title string) jquery.JQuery {
// 	g := jq("<div>").AddClass("option-group")
// 	h := jq("<div>").AddClass("option-group-header").SetText(title)
// 	return g.Append(h)
// }

// func optionList(title string) jquery.JQuery {
// 	g := jq("<div>").AddClass("option-list")
// 	h := jq("<div>").AddClass("option-group-header").SetText(title)
// 	return g.Append(h)
// }

// func numberOption(label, min, max, value string) jquery.JQuery {
// 	o := option(label)
// 	in := jq("<input>").SetAttr("type", "number").SetAttr("min", min).SetAttr("max", max).SetAttr("value", value)
// 	return o.Append(in)
// }

// func textOption(label string) jquery.JQuery {
// 	o := option(label)
// 	in := jq("<input>").SetAttr("type", "text")
// 	return o.Append(in)
// }

// func choiceOption(label string, choices []string) jquery.JQuery {
// 	o := option(label)
// 	s := jq("<select>")
// 	for _, c := range choices {
// 		s.Append(jq("<option>").SetAttr("value", c).SetText(c))
// 	}
// 	return o.Append(s)
// }

// func option(label string) jquery.JQuery {
// 	o := jq("<div>").AddClass("option")
// 	o.Append(label + " ")
// 	o.Call(jquery.CLICK, func(event jquery.Event) {
// 		event.StopPropagation()
// 	})
// 	return o
// }

// func vectorOption() jquery.JQuery {
// 	return jq("<span>") //.Append(numberOption("x", -999, 999, 0)).Append(numberOption("y", -999, 999, 0)).
// 	//	Append(numberOption("z", -999, 999, 0))
// }

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
