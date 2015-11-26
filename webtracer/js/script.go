package main

import (
	"reflect"
	"strconv"

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
		w, h := img.Get("width").Float(), img.Get("height").Float()
		imgCon.SetData("initWidth", w)
		imgCon.SetData("initHeight", h)
		console.Call("log", w, h)
		setImageSize(w, h)
	})
	img.Set("src", "/img/render542.png")

	jqImg := jq(img)
	jqImg.Call("draggable")

	zoom := jq("#zoom")
	zoom.SetAttr("value", 1.0)
	zoom.SetAttr("min", 0.1)
	zoom.SetAttr("max", 10.0)
	zoom.SetAttr("step", 0.1)
}

func initCallbacks() {
	jq(".tab").Call(jquery.CLICK, onToggleControls)
	jq("#save").Call(jquery.CLICK, onSave)
	jq("#load").Call(jquery.CLICK, onLoad)
	jq("#zoom").On("input change", onZoom)
}

func setImageSize(w, h float64) {
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

func onZoom() {
	newZoom, e := strconv.ParseFloat(jq("#zoom").Val(), 64)
	if e != nil {
		console.Call("error", e.Error())
	}
	imgCon := jq(".image-container")
	newW := imgCon.Data("initWidth").(float64) * newZoom
	newH := imgCon.Data("initHeight").(float64) * newZoom
	setImageSize(newW, newH)
}
