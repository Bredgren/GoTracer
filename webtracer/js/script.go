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

	addOptionSlides(o)
}

func addOptionSlides(opts jquery.JQuery) {
	structs := jq(".go-struct-field > .go-struct")
	structs.Each(func(i int, intf interface{}) {
		obj := intf.(*js.Object)
		label := jq(obj.Get("parentNode").Get("children").Index(0))
		st := jq(obj)
		jq(label).Call(jquery.CLICK, func(event jquery.Event) {
			st.SlideToggle("fast")
			event.StopPropagation()
		})
	})

	slices := jq(".go-struct-field > .go-slice")
	slices.Each(func(i int, intf interface{}) {
		obj := intf.(*js.Object)
		label := jq(obj.Get("parentNode").Get("children").Index(0))
		st := jq(obj)
		jq(label).Call(jquery.CLICK, func(event jquery.Event) {
			st.SlideToggle("fast")
			event.StopPropagation()
		})
	})
}

func onToggleControls() {
	jq("#options").SlideToggle("fast")
	jq("#tab-up").Toggle()
	jq("#tab-down").Toggle()
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
