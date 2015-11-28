package main

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"

	"github.com/Bredgren/gohtmlctrl/htmlctrl"
	"github.com/Bredgren/gotracer/lib"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
)

var jq = jquery.NewJQuery
var console = js.Global.Get("console")
var options *lib.Options

func main() {
	js.Global.Set("onBodyLoad", onBodyLoad)
}

func onBodyLoad() {
	initGlobalCallbacks()
	options = lib.NewOptions()
	initOptions()

	setImage("/img/render542.png")

	zoom := jq("#zoom")
	zoom.SetAttr("value", 1.0)
	zoom.SetAttr("min", 0.1)
	zoom.SetAttr("max", 10.0)
	zoom.SetAttr("step", 0.1)
}

func initGlobalCallbacks() {
	jq(".tab").Call(jquery.CLICK, onToggleControls)
	jq("#save").Call(jquery.CLICK, onSave)
	jq("#load").Call(jquery.CLICK, onLoad)
	jq("#load-file").Call(jquery.CHANGE, onFileChange)
	jq("#zoom").On("input change", onZoom)
	jq("#reset").Call(jquery.CLICK, onReset)
}

func initOptionCallbacks() {
	fn := func(i int, intf interface{}) {
		jq(intf.(*js.Object)).Off(".option")
		jq(intf.(*js.Object)).On(jquery.CLICK+".option", onOptionChange)
	}
	jq(".go-bool").Each(fn)
	jq(".go-int").Each(fn)
	jq(".go-float64").Each(fn)
	jq(".go-string").Each(fn)
	jq(".go-choice").Each(fn)
	jq(".go-slice button").Each(fn)
}

func setImage(path string) {
	img := js.Global.Get("Image").New()
	img.Set("onload", func() {
		imgCon := jq(".image-container")
		imgCon.Empty()
		imgCon.Append(img)
		w, h := img.Get("naturalWidth").Float(), img.Get("naturalHeight").Float()
		imgCon.SetData("initWidth", w)
		imgCon.SetData("initHeight", h)
		setImageSize(w, h)
		onReset()
	})
	img.Set("onerror", func() {
		console.Call("error", "unable to load image "+path)
	})
	img.Set("src", path)
	jq(img).Call("draggable")
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

func initOptions() {
	opts := jq("#options")
	opts.Empty()
	o, e := htmlctrl.Struct(options, "Options", "all-options", "")
	if e != nil {
		console.Call("error", e.Error())
		return
	}
	opts.Append(o)

	addOptionSlides(o)
	initOptionCallbacks()
	checkFastRender()
}

func addOptionSlides(opts jquery.JQuery) {
	structs := jq(".go-struct-field > .go-struct")
	structs.Each(func(i int, intf interface{}) {
		obj := intf.(*js.Object)
		label := jq(obj.Get("parentNode").Get("children").Index(0))
		st := jq(obj)
		jq(label).Off(".option")
		jq(label).On(jquery.CLICK+".option", func(event jquery.Event) {
			st.SlideToggle("fast")
			event.StopPropagation()
		})
	})

	slices := jq(".go-struct-field > .go-slice")
	slices.Each(func(i int, intf interface{}) {
		obj := intf.(*js.Object)
		label := jq(obj.Get("parentNode").Get("children").Index(0))
		st := jq(obj)
		jq(label).Off(".option")
		jq(label).On(jquery.CLICK+".option", func(event jquery.Event) {
			st.SlideToggle("fast")
			event.StopPropagation()
		})
	})
}

func checkFastRender() {
	fr := jq("#fast-render").Is(":checked")
	jq(".fast-render").Each(func(i int, intf interface{}) {
		obj := jq(intf.(*js.Object))
		if fr {
			obj.AddClass("fast-render-on")
		} else {
			obj.RemoveClass("fast-render-on")
		}
	})
}

func onToggleControls() {
	jq("#options").SlideToggle("fast")
	jq("#tab-up").Toggle()
	jq("#tab-down").Toggle()
}

func onSave() {
	console.Call("log", "save")

	j, e := json.Marshal(options)
	if e != nil {
		console.Call("error", e.Error())
		return
	}
	text := string(j)

	csvData := "data:application/csv;charset=utf-8," + js.Global.Call("encodeURIComponent", text).String()
	l := jq("#save")
	l.SetProp("href", csvData)
	l.SetProp("target", "_blank")
	l.SetProp("download", "scene.json")
}

func onLoad() {
	jq("#load-file").Call("trigger", "click")
}

func onFileChange(event *js.Object) {
	reader := js.Global.Get("FileReader").New()
	reader.Set("onload", func(evt *js.Object) {
		content := evt.Get("target").Get("result").String()
		e := json.Unmarshal([]byte(content), options)
		if e != nil {
			console.Call("error", e)
			return
		}
		initOptions()
		triggerRender()
	})
	reader.Call("readAsText", event.Get("target").Get("files").Index(0))
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

func onReset() {
	jq("#zoom").SetVal(1.0)
	onZoom()
	jq("img").SetCss("left", 0)
	jq("img").SetCss("top", 0)
}

func onOptionChange() {
	// Re-add callbacks because when adding/removing items from slices they destroy the previous objects
	// and new ones will be missing the callbacks.
	addOptionSlides(jq("#all-options"))
	initOptionCallbacks()
	checkFastRender()

	triggerRender()
}

func triggerRender() {
	j, e := json.Marshal(options)
	if e != nil {
		console.Call("error", e.Error())
		return
	}

	jq(".image-container").Empty()
	done := make(chan bool)
	go func(done <-chan bool) {
		anim := jq("#animation")
		anim.SetCss("color", "rgb(255,255,255)")
		anim.Show()
		for i := 0.0; ; i += 0.15 {
			select {
			case <-done:
				anim.Hide()
				return
			case <-time.After(time.Duration(33) * time.Millisecond):
				anim.SetCss("color", fmt.Sprintf("rgb(%d,%[1]d,%[1]d)", int((math.Cos(i)*0.4+(1-0.4))*255.0)))
			}
		}
	}(done)
	jquery.Post("/", string(j), func(data, status, xhr string) {
		if status != "success" {
			console.Call("error", "Render wasn't success:", status, data, xhr)
		}
		// Not allowed to block in callback
		go func() { done <- true }()
		setImage(data)
	})
}
