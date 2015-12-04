package main

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/Bredgren/gohtmlctrl/htmlctrl"
	"github.com/Bredgren/gotracer/trace"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
)

var jq = jquery.NewJQuery
var console = js.Global.Get("console")
var options *trace.Options

func main() {
	js.Global.Set("onBodyLoad", onBodyLoad)
	htmlctrl.RegisterValidator("validXPixel", htmlctrl.ValidateInt(func(x int) bool {
		return x < int(jq(".image-container").Data("initWidth").(float64))
	}))
	htmlctrl.RegisterValidator("validYPixel", htmlctrl.ValidateInt(func(y int) bool {
		return y < int(jq(".image-container").Data("initHeight").(float64))
	}))
}

func onBodyLoad() {
	initGlobalCallbacks()
	options = trace.NewOptions()
	initOptions()

	triggerRender()

	zoom := jq("#zoom")
	zoom.SetAttr("value", 1.0)
	zoom.SetAttr("min", 0.1)
	zoom.SetAttr("max", 10.0)
	zoom.SetAttr("step", 0.1)

	refreshHistory()
}

func initGlobalCallbacks() {
	jq(".option-tab").Call(jquery.CLICK, onToggleOptions)
	jq(".history-tab").Call(jquery.CLICK, onToggleHistory)
	jq("#save").Call(jquery.CLICK, onSave)
	jq("#load").Call(jquery.CLICK, onLoad)
	jq("#load-file").Call(jquery.CHANGE, onFileChange)
	jq("#render").Call(jquery.CLICK, onOptionChange)
	jq("#zoom").On("input change", onZoom)
	jq("#reset").Call(jquery.CLICK, onReset)
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

func onToggleOptions() {
	jq("#options").SlideToggle("fast")
	jq("#tab-up").Toggle()
	jq("#tab-down").Toggle()
}

func onToggleHistory() {
	jq("#history-list").SlideToggle("fast")
	jq("#history-tab-up").Toggle()
	jq("#history-tab-down").Toggle()
}

func onSave() {
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
	// Not allowing options to be changed while rendering
	jq("#options input").SetProp("disabled", true)

	// Re-add callbacks because when adding/removing items from slices they destroy the previous objects
	// and new ones will be missing the callbacks.
	addOptionSlides(jq("#all-options"))
	checkFastRender()

	triggerRender()
	refreshHistory()
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
	jquery.Post("/render", string(j), func(data, status, xhr string) {
		if status != "success" {
			console.Call("error", "Render wasn't success:", status, data, xhr)
		}
		// Not allowed to block in callback
		go func() { done <- true }()
		jq("#options input").SetProp("disabled", false)
		setImage(data)
	})
}

func refreshHistory() {
	jquery.Post("/history", "", func(data, status, xhr string) {
		if status != "success" {
			console.Call("error", "Failed to retrieve history:", status, data, xhr)
		}
		var history [][2]string
		e := json.Unmarshal([]byte(data), &history)
		if e != nil {
			console.Call("error", "Unmarshaling history:", e.Error())
			return
		}
		var done chan bool
		go populateHistory(history, done)
		go func() {
			<-done
			console.Call("log", "done populating history")
		}()
	})
}

func populateHistory(history [][2]string, done chan<- bool) {
	imgList := make([]*js.Object, len(history))
	var wg sync.WaitGroup
	for i, p := range history {
		i, p := i, p
		wg.Add(1)
		img := js.Global.Get("Image").New()
		img.Set("onload", func() {
			wg.Done()
			imgList[i] = img
		})
		img.Set("onerror", func() {
			wg.Done()
			console.Call("error", "unable to load image "+p[0])
		})
		img.Set("src", p[0])
	}
	wg.Wait()
	list := jq("#history-list")
	list.Empty()
	for i := len(imgList) - 1; i >= 0; i-- {
		img := imgList[i]
		item := jq("<div>").AddClass("history-item")
		text := strconv.Itoa(i)
		label := jq("<label>").SetText(text)
		item.Append(label)
		item.Append(img)
		item.On(jquery.CLICK, getOnHistoryItemSelect(item, history[i][0], history[i][1]))
		list.Append(item)
	}
	done <- true
}

func getOnHistoryItemSelect(item jquery.JQuery, img, scene string) func() {
	return func() {
		jq(".history-item").RemoveClass("history-selected")
		item.AddClass("history-selected")
		setImage(img)
		setOptions(scene)
	}
}

func setOptions(scene string) {
	jquery.Get(scene, "", func(data *js.Object, status, xhr string) {
		if status != "success" {
			console.Call("error", "Failed to retrieve scene:", status, data, xhr)
		}

		e := json.Unmarshal([]byte(data.String()), &options)
		if e != nil {
			console.Call("error", "Unmarshaling options:", e.Error())
			return
		}
		initOptions()
	})
}
