package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/Bredgren/gotracer/lib"
)

const path = "/src/github.com/Bredgren/gotracer/webtracer"

var goPath = os.Getenv("GOPATH")

var templ *template.Template

var (
	// debug is set by the -D command line flag.
	debug bool
	port  int
)

func init() {
	flag.BoolVar(&debug, "D", false, "Debug mode. More logs, use unminified assets, etc.")
	flag.IntVar(&port, "p", 8080, "Set http port")
}

func main() {
	setup()

	http.HandleFunc("/", httpHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func setup() {
	e := os.Chdir(filepath.Join(goPath, path))
	if e != nil {
		log.Fatal(e)
	}
	flag.Parse()

	if debug {
		log.Println("Debug mode enabled!")
	}

	templ = template.Must(template.New("templ").Funcs(template.FuncMap{
		"debug": func() bool { return debug },
	}).ParseFiles("./tmpl/page.tmpl"))
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI != "/" {
		http.ServeFile(w, r, "./"+r.RequestURI)
		return
	}

	switch r.Method {
	case "POST":
		log.Println("POST", r.RequestURI)
		body, e := ioutil.ReadAll(r.Body)
		if e != nil {
			log.Println("Error reading post:", e)
			http.Error(w, e.Error(), http.StatusBadRequest)
			return
		}
		r.Body.Close()
		var options lib.Options
		e = json.Unmarshal(body, &options)
		if e != nil {
			log.Println("Error decoding json:", e)
			http.Error(w, e.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println(options)
		// Simulate render time with random wait
		<-time.After(time.Duration(rand.Intn(10)) * time.Second)
		fmt.Fprintln(w, "/img/render542.png")
	case "GET":
		log.Println("GET", r.RequestURI)
		renderTmpl(w, &page{})
	}
}

func renderTmpl(w http.ResponseWriter, p *page) {
	e := templ.ExecuteTemplate(w, "Main", p)
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
}

type page struct {
}
